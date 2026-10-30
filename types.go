// Package thsdk is a Go port of the THSDK Python SDK for real-time and
// historical securities data, backed by the same native runtime component.
package thsdk

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// Shanghai is the timezone the runtime component expects for timestamps.
var Shanghai = time.FixedZone("Asia/Shanghai", 8*60*60)

// Security identifies one instrument: a 4-letter market plus the code itself.
// Full codes look like USHA600519 (market USHA, code 600519).
type Security struct {
	Code   string `json:"code"`
	Market string `json:"market"`
}

// Block identifies a sector or index board, by numeric id or by name.
type Block struct {
	ID   *int   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type blockAlias Block

func (b Block) MarshalJSON() ([]byte, error) {
	if b.ID == nil && b.Name == "" {
		return nil, fmt.Errorf("thsdk: 板块标识必须包含 id 或 name")
	}
	return json.Marshal(struct {
		ID   *int   `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}{ID: b.ID, Name: b.Name})
}

var marketPattern = regexp.MustCompile(`^[A-Za-z]{4}$`)

// NewSecurity parses a full code such as USHA600519.
func NewSecurity(fullCode string) (Security, error) {
	fullCode = strings.TrimSpace(fullCode)
	if len(fullCode) <= 4 {
		return Security{}, fmt.Errorf("thsdk: 请使用完整证券代码，例如 USHA600519 或 USZA300033，收到 %q", fullCode)
	}
	market, code := fullCode[:4], fullCode[4:]
	if !marketPattern.MatchString(market) || code == "" {
		return Security{}, fmt.Errorf("thsdk: 请使用完整证券代码，例如 USHA600519 或 USZA300033，收到 %q", fullCode)
	}
	return Security{Code: code, Market: strings.ToUpper(market)}, nil
}

// FullCode returns the market-prefixed form, e.g. USHA600519.
func (s Security) FullCode() string { return s.Market + s.Code }

// String lets a Security be passed directly where a full code is expected.
func (s Security) String() string { return s.FullCode() }

// ParseMarket validates a 4-letter market code.
func ParseMarket(market string) (string, error) {
	market = strings.TrimSpace(market)
	if !marketPattern.MatchString(market) {
		return "", fmt.Errorf("thsdk: 市场代码必须是四位字母，例如 USHA 或 USZA，收到 %q", market)
	}
	return strings.ToUpper(market), nil
}

// MarshalJSON normalizes the identity on the way out, so a Security built from
// "USHA600519" via NewSecurity, or hand-written in lower case, both wire up the
// same way.
func (s Security) MarshalJSON() ([]byte, error) {
	market, err := ParseMarket(s.Market)
	if err != nil {
		return nil, err
	}
	code := strings.TrimSpace(s.Code)
	if code == "" {
		return nil, fmt.Errorf("thsdk: 证券代码不能为空")
	}
	return json.Marshal(struct {
		Code   string `json:"code"`
		Market string `json:"market"`
	}{Code: code, Market: market})
}

// ParseSecurity accepts either a full code ("USHA600519") or a value that
// already is one, and is what the convenience layer uses.
func ParseSecurity(value any) (Security, error) {
	switch typed := value.(type) {
	case Security:
		return typed, nil
	case *Security:
		if typed == nil {
			return Security{}, fmt.Errorf("thsdk: 证券不能为 nil")
		}
		return *typed, nil
	case string:
		return NewSecurity(typed)
	case map[string]any:
		code, _ := typed["code"].(string)
		market, _ := typed["market"].(string)
		return Security{Code: strings.TrimSpace(code), Market: strings.ToUpper(strings.TrimSpace(market))}, nil
	default:
		return Security{}, fmt.Errorf("thsdk: 不支持的证券类型 %T", value)
	}
}

// ParseMarkets accepts one space- or comma-separated string, or a list.
func ParseMarkets(markets ...string) ([]string, error) {
	var out []string
	for _, item := range markets {
		for _, token := range strings.FieldsFunc(item, func(r rune) bool {
			return r == ',' || unicode.IsSpace(r)
		}) {
			market, err := ParseMarket(token)
			if err != nil {
				return nil, err
			}
			out = append(out, market)
		}
	}
	return out, nil
}

func securityList(securities []Security) ([]Security, error) {
	if len(securities) == 0 {
		return nil, fmt.Errorf("thsdk: 证券列表不能为空")
	}
	for _, item := range securities {
		if !marketPattern.MatchString(item.Market) || item.Code == "" {
			return nil, fmt.Errorf("thsdk: 请使用完整证券代码，例如 USHA600519，收到 %q", item.String())
		}
	}
	out := make([]Security, len(securities))
	for i, item := range securities {
		out[i] = Security{Code: item.Code, Market: strings.ToUpper(item.Market)}
	}
	return out, nil
}

// datePattern matches the YYYYMMDD form the component uses for date arguments.
func dateNumber(value time.Time) int {
	return value.Year()*10000 + int(value.Month())*100 + value.Day()
}

func dateString(value time.Time) string { return fmt.Sprintf("%08d", dateNumber(value)) }

// formatTimestamp renders an ISO-8601 timestamp in Shanghai time, matching the
// Python SDK's datetime.isoformat() output for naive and aware values alike.
func formatTimestamp(value time.Time) string {
	value = value.In(Shanghai)
	base := value.Format("2006-01-02T15:04:05")
	if nanos := value.Nanosecond(); nanos != 0 {
		base = fmt.Sprintf("%s.%06d", base, nanos/1000)
	}
	return base + value.Format("-07:00")
}

// MustSecurity panics when fullCode is not a valid complete code; it is for
// literals in examples and configuration.
func MustSecurity(fullCode string) Security {
	security, err := NewSecurity(fullCode)
	if err != nil {
		panic(err)
	}
	return security
}
