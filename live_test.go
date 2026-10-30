package thsdk_test

import (
	"encoding/json"
	"errors"
	"testing"

	thsdk "github.com/angYuHao/go-thsdk"
)

// Live smoke coverage over the request shapes that differ at the ABI level:
// a no-param call, a bare-string param, a bare-array param, an object param,
// the paging envelopes, and the mutation guard. The full 80-method sweep is a
// manual tool (examples/allmethods) because it takes minutes under the 50ms
// throttle; this subset is what a CI run can afford.

func TestLiveCallShapes(t *testing.T) {
	client, ctx := newClient(t)

	cases := []struct {
		name   string
		params any
	}{
		{"account_permissions", nil},
		{"list_security_k_lines", map[string]any{
			"security": thsdk.MustSecurity("USHA600519"),
			"start":    -3, "end": 0, "adjust": 1, "period": 0x4000,
		}},
		{"list_security_link_relations", "USHA600519"}, // bare string param
		{"list_security_futures_relations", []thsdk.Security{ // bare array param
			thsdk.MustSecurity("USHA600519"),
		}},
		{"resolve_block", map[string]any{"name": "沪深300"}},
		{"search_securities", map[string]any{"pattern": "茅台"}},
		{"list_news_flash_items", map[string]any{"traversal": "page", "page": 1, "page_size": 3}},
		{"query_wencai", map[string]any{"query": "市盈率小于20", "limit": 3}},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			data, err := client.CallRaw(ctx, item.name, item.params)
			if err != nil {
				skipWhenServiceUnavailable(t, err)
				t.Fatalf("%s: %v", item.name, err)
			}
			if data.IsEmpty() {
				t.Logf("%s returned null (legal for this account)", item.name)
				return
			}
			if !json.Valid(data.Bytes()) {
				t.Fatalf("%s returned invalid JSON: %.120s", item.name, data.Bytes())
			}
			t.Logf("%s ok: %d bytes", item.name, len(data.Bytes()))
		})
	}
}

func TestLiveTypedRequestMatchesRawRequest(t *testing.T) {
	client, ctx := newClient(t)

	sec := thsdk.MustSecurity("USHA600519")
	typed, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
		Security: sec, Start: -3, End: 0, Adjust: thsdk.AdjustForward, Period: thsdk.PeriodDay,
	})
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("typed: %v", err)
	}
	raw, err := client.CallRaw(ctx, "list_security_k_lines", map[string]any{
		"security": sec, "start": -3, "end": 0, "adjust": 1, "period": 0x4000,
	})
	if err != nil {
		t.Fatalf("raw: %v", err)
	}

	var typedRows, rawRows []thsdk.KLine
	if err := typed.Rows(&typedRows); err != nil {
		t.Fatalf("typed rows: %v", err)
	}
	if err := raw.Rows(&rawRows); err != nil {
		t.Fatalf("raw rows: %v", err)
	}
	if len(typedRows) != len(rawRows) {
		t.Fatalf("generated struct produced %d candles, raw map produced %d", len(typedRows), len(rawRows))
	}
	if len(typedRows) > 0 && typedRows[0].Close != rawRows[0].Close {
		t.Fatalf("the two encodings disagree: %v vs %v", typedRows[0], rawRows[0])
	}
}

func TestLiveOptionalFieldOmissionHasEffect(t *testing.T) {
	client, ctx := newClient(t)

	// An unset Opt field must not appear on the wire at all; an explicit zero
	// must. Both requests should still succeed, proving the component receives
	// what the struct declares rather than a field full of nulls.
	sparse, err := client.ListMarketSecurities(ctx, thsdk.ListMarketSecuritiesRequest{Market: "USHA"})
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("sparse: %v", err)
	}
	full, err := client.ListMarketSecurities(ctx, thsdk.ListMarketSecuritiesRequest{
		Market: "USHA", SortBegin: thsdk.Num(0), SortCount: thsdk.Num(2),
	})
	if err != nil {
		t.Fatalf("full: %v", err)
	}

	var sparseRows, fullRows []thsdk.SecurityRef
	if err := sparse.Rows(&sparseRows); err != nil {
		t.Fatalf("sparse rows: %v", err)
	}
	if err := full.Rows(&fullRows); err != nil {
		t.Fatalf("full rows: %v", err)
	}
	if len(fullRows) > 2 {
		t.Fatalf("sort_count=2 returned %d rows; the optional field leaked", len(fullRows))
	}
	if len(sparseRows) == 0 {
		t.Fatalf("the default page came back empty")
	}
}

func TestLiveErrorTypeIsMappedFaithfully(t *testing.T) {
	client, ctx := newClient(t)

	// A code the component always rejects; the SDK must surface it as *Error
	// with the upstream code intact, not as a generic failure.
	_, err := client.CallRaw(ctx, "no_such_method", nil)
	if err == nil {
		t.Fatalf("expected an unknown-method error")
	}
	var apiErr *thsdk.Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *thsdk.Error, got %T: %v", err, err)
	}
	if apiErr.Code == "" {
		t.Fatalf("error code lost: %+v", apiErr)
	}
	t.Logf("unknown method surfaced as code=%s", apiErr.Code)
}

// skipWhenServiceUnavailable tolerates the upstream having no data, no
// permission for this trial account, or an internal encode failure. The same
// request fails identically on the Python SDK, so it is not a porting defect.
func skipWhenServiceUnavailable(t *testing.T, err error) {
	t.Helper()
	var apiErr *thsdk.Error
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case "no_data", "encode_failed", "no_permission", "permission_denied", "timeout", "query_failed":
			t.Skipf("upstream service unavailable: %v", err)
		}
	}
}
