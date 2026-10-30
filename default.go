package thsdk

import (
	"context"
	"sync"
	"time"
)

// The package-level entry points mirror the Python module, which keeps one
// process-wide client so callers never have to pass it around.

var (
	defaultOnce   sync.Once
	defaultClient *Client
	defaultErr    error
)

// Default returns the process-wide client, loading the runtime component on
// first use.
func Default() (*Client, error) {
	defaultOnce.Do(func() {
		defaultClient, defaultErr = New(Options{})
	})
	return defaultClient, defaultErr
}

// Auth completes the package client's login handshake.
func Auth(ctx context.Context) error {
	client, err := Default()
	if err != nil {
		return err
	}
	return client.Auth(ctx, Credentials{})
}

// AuthQrcode restores a local session or waits for a scan.
func AuthQrcode(ctx context.Context) error {
	client, err := Default()
	if err != nil {
		return err
	}
	return client.AuthQrcode(ctx, 120*time.Second, time.Second)
}

// Logout clears the account and its stored session.
func Logout(ctx context.Context) error {
	client, err := Default()
	if err != nil {
		return err
	}
	return client.Logout(ctx)
}

// Close releases the package client. A later Default() loads a fresh one.
func Close() error {
	defaultOnce.Do(func() {})

	client := defaultClient
	defaultOnce = sync.Once{}
	defaultClient, defaultErr = nil, nil
	if client == nil {
		return nil
	}
	return client.Close()
}
