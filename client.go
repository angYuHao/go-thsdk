package thsdk

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/angYuHao/go-thsdk/internal/native"
)

// Environment variables read by Auth when no credentials are passed.
const (
	EnvUsername = "THS_USERNAME"
	EnvPassword = "THS_PASSWORD"
)

// SessionFileName is the reusable session stored inside SessionDir.
const SessionFileName = "account.session"

const (
	accountMACSeed        = "thsdk-account-mac-v1\x00"
	maxTimeoutSeconds     = 600
	minQRPollMilliseconds = 500
	maxQRPollMilliseconds = 30000
	maxRateLimitRetryMS   = 1000
)

// Options configures a Client.
type Options struct {
	// SessionDir holds account.session. Defaults to the working directory.
	SessionDir string
	// RuntimePath points at hq.dll / hq.so / hq.dylib. When empty the SDK looks
	// at THSDK_RUNTIME_PATH and the usual package layouts.
	RuntimePath string
	// Mac overrides the device address derived from the account name.
	Mac string
	// DisableTemporaryAccount stops Auth from falling back to the bundled trial
	// accounts. Production deployments should set it and use real credentials.
	DisableTemporaryAccount bool
	// AuthTimeout bounds the login handshake when ctx carries no deadline.
	// Zero means 60 seconds; the maximum is 600.
	AuthTimeout time.Duration
	// OnQRChallenge receives the login QR payload during AuthQrcode. When nil
	// and Verbose is set, the scan URL is printed to stderr.
	OnQRChallenge func(QRChallenge)
	// Verbose prints the scan URL and image URL of a QR login challenge.
	Verbose bool
}

// QRChallenge is one polling state of a QR login job.
type QRChallenge struct {
	ScanURL       string `json:"scan_url"`
	ImageURL      string `json:"image_url"`
	State         string `json:"state"`
	Authenticated bool   `json:"authenticated"`
}

// Client owns one runtime handle and the authentication state on top of it.
type Client struct {
	options Options
	runtime *native.Runtime

	callMu        sync.Mutex
	authenticated bool
}

// New loads the runtime component and returns an unauthenticated client.
func New(options Options) (*Client, error) {
	path, err := resolveComponent(options.RuntimePath)
	if err != nil {
		return nil, err
	}
	if options.SessionDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		options.SessionDir = dir
	}
	sessionDir, err := filepath.Abs(options.SessionDir)
	if err != nil {
		return nil, err
	}
	options.SessionDir = sessionDir

	runtime, err := native.Open(path)
	if err != nil {
		return nil, err
	}
	return &Client{options: options, runtime: runtime}, nil
}

// Close releases the runtime component.
func (c *Client) Close() error {
	c.callMu.Lock()
	defer c.callMu.Unlock()
	c.authenticated = false
	return c.runtime.Close()
}

// Authenticated reports whether this client has completed an auth handshake.
func (c *Client) Authenticated() bool {
	c.callMu.Lock()
	defer c.callMu.Unlock()
	return c.authenticated
}

// CallRaw performs one runtime call, returning the untouched payload.
// Use it for methods without a typed wrapper.
func (c *Client) CallRaw(ctx context.Context, method string, params any) (RawData, error) {
	return c.request(ctx, method, params)
}

// request enforces authentication and the component's single rate-limit retry.
func (c *Client) request(ctx context.Context, method string, params any) (RawData, error) {
	c.callMu.Lock()
	authenticated := c.authenticated
	c.callMu.Unlock()
	if !authenticated {
		return RawData{}, ErrNotAuthenticated
	}
	return c.invoke(ctx, method, params)
}

// invoke is used by both authenticated calls and the auth handshake itself.
// A context deadline also becomes the component-side timeout_ms budget.
func (c *Client) invoke(ctx context.Context, method string, params any) (RawData, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	deadline, hasDeadline := time.Time{}, false
	if value, ok := ctx.Deadline(); ok {
		deadline, hasDeadline = value, true
	}

	timeoutMS := 0
	if hasDeadline {
		if time.Until(deadline) <= 0 {
			return RawData{}, ctx.Err()
		}
		timeoutMS = millisecondsOf(time.Until(deadline))
	}

	for attempt := 0; attempt < 2; attempt++ {
		if err := ctx.Err(); err != nil {
			return RawData{}, err
		}
		data, err := c.runtime.Invoke(method, params, timeoutMS)
		if err == nil {
			return RawData{payload: data}, nil
		}

		var apiErr *native.APIError
		if !errors.As(err, &apiErr) {
			return RawData{}, translateError(err)
		}
		retryAfter, ok := apiErr.RetryAfterMS()
		canRetry := attempt == 0 && ok && apiErr.Code == "rate_limited" && retryAfter <= maxRateLimitRetryMS
		if canRetry && hasDeadline {
			canRetry = time.Now().Add(time.Duration(retryAfter) * time.Millisecond).Before(deadline)
		}
		if !canRetry {
			return RawData{}, translateError(err)
		}
		select {
		case <-time.After(time.Duration(retryAfter) * time.Millisecond):
		case <-ctx.Done():
			return RawData{}, translateError(err)
		}
		if hasDeadline {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return RawData{}, translateError(err)
			}
			timeoutMS = millisecondsOf(remaining)
		}
	}
	return RawData{}, translateError(errors.New("thsdk: 调用未达预期结束状态"))
}

func millisecondsOf(d time.Duration) int {
	ms := int(d / time.Millisecond)
	if ms < 1 {
		return 1
	}
	if ms > native.MaxCallTimeoutMS {
		return native.MaxCallTimeoutMS
	}
	return ms
}

func translateError(err error) error {
	var apiErr *native.APIError
	if errors.As(err, &apiErr) {
		return &Error{
			Code:          apiErr.Code,
			Message:       apiErr.Message,
			RetryAfter:    func() int { value, _ := apiErr.RetryAfterMS(); return value }(),
			HasRetryAfter: func() bool { _, ok := apiErr.RetryAfterMS(); return ok }(),
		}
	}
	return err
}

// Auth logs in with explicit credentials, then THS_USERNAME/THS_PASSWORD, then
// a stored session, then (unless disabled) a bundled temporary account.
//
// Pass an empty Credentials to skip straight to the environment/session chain.
// The budget is ctx's deadline when set, otherwise Options.AuthTimeout,
// otherwise 60 seconds.
func (c *Client) Auth(ctx context.Context, cred Credentials) error {
	ctx, finish, err := c.authContext(ctx)
	defer finish()
	if err != nil {
		return err
	}
	if cred.set() {
		if strings.TrimSpace(cred.Username) == "" || cred.Password == "" {
			return newAuthError("账号和密码不能为空", nil)
		}
		return c.authenticate(ctx, c.accountParams(cred))
	}

	env, found, err := c.environmentCredentials()
	if err != nil {
		return err
	}
	if found {
		return c.authenticate(ctx, c.accountParams(env))
	}

	c.callMu.Lock()
	authenticated := c.authenticated
	c.callMu.Unlock()
	if authenticated {
		return nil
	}

	if _, statErr := os.Stat(filepath.Join(c.options.SessionDir, SessionFileName)); statErr == nil {
		err := c.authenticate(ctx, map[string]any{"session_dir": c.options.SessionDir})
		if err == nil {
			return nil
		}
		var apiErr *Error
		if !errors.As(err, &apiErr) || apiErr.Code != "session_unavailable" {
			return err
		}
	}

	if !c.temporaryAllowed() {
		return newAuthError("没有可用凭据；请提供账号密码、设置 THS_USERNAME/THS_PASSWORD，或先执行扫码登录", nil)
	}
	username, password, mac := pickTemporaryAccount()
	return c.authenticate(ctx, c.accountParams(Credentials{Username: username, Password: password, Mac: mac}))
}

// Credentials is one account's login material. Mac is optional.
type Credentials struct {
	Username string
	Password string
	Mac      string
}

func (cred Credentials) set() bool {
	return cred.Username != "" || cred.Password != ""
}

// authContext bounds the handshake so a hung login cannot outlive the budget.
func (c *Client) authContext(ctx context.Context) (context.Context, func(), error) {
	if ctx == nil {
		ctx = context.Background()
	}
	budget := c.options.AuthTimeout
	if budget == 0 {
		budget = defaultAuthTimeout
	}
	if budget < 0 || budget > maxTimeoutSeconds*time.Second {
		return ctx, func() {}, invalidArgument("AuthTimeout", "必须大于 0 且不超过 %d 秒", maxTimeoutSeconds)
	}
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}, nil
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, budget)
	return timeoutCtx, cancel, nil
}

const defaultAuthTimeout = 60 * time.Second

func (c *Client) temporaryAllowed() bool {
	return !c.options.DisableTemporaryAccount
}

func (c *Client) environmentCredentials() (Credentials, bool, error) {
	username, hasUser := os.LookupEnv(EnvUsername)
	password, hasPass := os.LookupEnv(EnvPassword)
	if !hasUser && !hasPass {
		return Credentials{}, false, nil
	}
	if strings.TrimSpace(username) == "" || password == "" {
		return Credentials{}, false, newAuthError(
			"环境变量 THS_USERNAME 和 THS_PASSWORD 必须同时设置且不能为空", nil)
	}
	return Credentials{Username: username, Password: password}, true, nil
}

func (c *Client) accountParams(cred Credentials) map[string]any {
	mac := cred.Mac
	if mac == "" {
		mac = c.options.Mac
	}
	if mac == "" {
		mac = deriveAccountMAC(cred.Username)
	}
	return map[string]any{
		"session_dir": c.options.SessionDir,
		"username":    cred.Username,
		"password":    cred.Password,
		"mac":         mac,
	}
}

func (c *Client) authenticate(ctx context.Context, params map[string]any) error {
	data, err := c.invoke(ctx, "auth", params)
	if err != nil {
		return newAuthError("认证失败", err)
	}
	var state struct {
		Authenticated bool `json:"authenticated"`
	}
	if len(data.payload) > 0 {
		if err := json.Unmarshal(data.payload, &state); err != nil {
			return newAuthError("认证响应无法解析", err)
		}
	}
	if !state.Authenticated {
		return newAuthError("认证未完成", nil)
	}
	c.callMu.Lock()
	c.authenticated = true
	c.callMu.Unlock()
	return nil
}

// Logout clears the account and deletes the reusable session file.
func (c *Client) Logout(ctx context.Context) error {
	_, err := c.invoke(ctx, "logout", map[string]any{"session_dir": c.options.SessionDir})
	c.callMu.Lock()
	c.authenticated = false
	c.callMu.Unlock()
	return err
}

// AuthQrcode restores a local session when possible, otherwise it runs one QR
// login job and waits for the scan. Challenges go to Options.OnQRChallenge.
func (c *Client) AuthQrcode(ctx context.Context, timeout time.Duration, pollInterval time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 || timeout > maxTimeoutSeconds*time.Second {
		return invalidArgument("timeout", "必须大于 0 且不超过 %d 秒", maxTimeoutSeconds)
	}
	if pollInterval < minQRPollMilliseconds*time.Millisecond || pollInterval > maxQRPollMilliseconds*time.Millisecond {
		return invalidArgument("poll_interval", "必须在 %d 到 %d 毫秒之间", minQRPollMilliseconds, maxQRPollMilliseconds)
	}

	deadline := time.Now().Add(timeout)
	start := map[string]any{
		"session_dir":      c.options.SessionDir,
		"timeout_seconds":  secondsCeil(timeout),
		"poll_interval_ms": millisecondsOf(pollInterval),
	}

	c.callMu.Lock()
	defer c.callMu.Unlock()

	state, err := c.pollQRCode(context.Background(), start, deadline)
	if err != nil {
		return newAuthError("扫码登录失败", err)
	}
	shown := QRChallenge{}
	if completed := c.takeQRState(state, &shown); completed {
		c.authenticated = true
		return nil
	}

	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			break
		}
		state, err = c.pollQRCode(ctx, nil, deadline)
		if err != nil {
			return newAuthError("扫码登录失败", err)
		}
		if c.takeQRState(state, &shown) {
			c.authenticated = true
			return nil
		}
		remaining = time.Until(deadline)
		if remaining <= 0 {
			break
		}
		wait := pollInterval
		if wait > remaining {
			wait = remaining
		}
		time.Sleep(wait)
	}

	// Release the pending job on the component side before reporting the timeout.
	_, _ = c.pollQRCode(context.Background(), map[string]any{"cancel": true}, time.Now().Add(time.Second))
	return newAuthError("二维码登录超时", nil)
}

func (c *Client) pollQRCode(ctx context.Context, params any, deadline time.Time) (RawData, error) {
	remaining := time.Until(deadline)
	if remaining <= 0 {
		return RawData{}, newAuthError("认证超时", nil)
	}
	return c.invoke(ctx, "auth_qrcode", params)
}

// takeQRState reports whether the challenge means "authenticated", and shows it
// when it carries a new scan target.
func (c *Client) takeQRState(data RawData, shown *QRChallenge) bool {
	if len(data.payload) == 0 {
		return false
	}
	var challenge QRChallenge
	if err := json.Unmarshal(data.payload, &challenge); err != nil {
		return false
	}
	if challenge.Authenticated || challenge.State == "authenticated" {
		return true
	}
	if (challenge.ScanURL == "" && challenge.ImageURL == "") || *shown == challenge {
		return false
	}
	*shown = challenge
	c.displayChallenge(challenge)
	return false
}

func (c *Client) displayChallenge(challenge QRChallenge) {
	if c.options.OnQRChallenge != nil {
		c.options.OnQRChallenge(challenge)
		return
	}
	if c.options.Verbose {
		fmt.Fprintf(os.Stderr, "请使用同花顺扫码登录：\n  %s\n", challenge.ScanURL)
		if challenge.ImageURL != "" {
			fmt.Fprintf(os.Stderr, "  二维码图片: %s\n", challenge.ImageURL)
		}
	}
}

// deriveAccountMAC maps one account name to a stable, locally administered
// unicast MAC, matching the Python SDK's _derive_account_mac.
func deriveAccountMAC(username string) string {
	digest := sha256.Sum256(append([]byte(accountMACSeed), username...))
	octets := digest[:6]
	octets[0] = (octets[0] & 0xFC) | 0x02
	parts := make([]string, 6)
	for i, value := range octets {
		parts[i] = fmt.Sprintf("%02x", value)
	}
	return strings.Join(parts, ":")
}

func secondsCeil(d time.Duration) int {
	seconds := int(d / time.Second)
	if d%time.Second != 0 {
		seconds++
	}
	if seconds < 1 {
		return 1
	}
	return seconds
}
