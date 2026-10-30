package thsdk

import (
	"encoding/json"
	"strings"
	"testing"
)

// The generated structs are the only place where the wire contract is spelled
// out in Go, so a wrong json tag or an over-eager omitempty would silently
// change what the component receives. These tests pin the emitted document.

func encode(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

func TestRequiredFieldsAlwaysPresent(t *testing.T) {
	got := encode(t, ListSecurityKLinesRequest{
		Security: Security{Code: "600519", Market: "USHA"},
		Start:    -3, End: 0, Adjust: AdjustNone, Period: PeriodDay,
	})
	want := `{"security":{"code":"600519","market":"USHA"},"start":-3,"end":0,"adjust":0,"period":16384}`
	if got != want {
		t.Fatalf("kline request\n got %s\nwant %s", got, want)
	}
}

func TestZeroValuedEnumsAreSent(t *testing.T) {
	// adjust=0 (no adjustment) and end=0 (up to now) are meaningful; dropping
	// them would change which series the component returns.
	body := encode(t, ListSecurityCorporateActionsRequest{
		Security: Security{Code: "600519", Market: "USHA"},
		Start:    -3, End: 0, Adjust: AdjustNone, Period: PeriodDay,
	})
	for _, needle := range []string{`"end":0`, `"adjust":0`, `"period":16384`} {
		if !strings.Contains(body, needle) {
			t.Fatalf("missing %s in %s", needle, body)
		}
	}
}

func TestUnsetOptionalsAreOmitted(t *testing.T) {
	got := encode(t, ListNewsFlashItemsRequest{Traversal: "page"})
	want := `{"traversal":"page"}`
	if got != want {
		t.Fatalf("optional fields leaked: %s", got)
	}
}

func TestExplicitZeroOptionalIsSent(t *testing.T) {
	// SortBegin 0 is a real offset; the caller opting in must survive marshaling.
	got := encode(t, ListMarketSecuritiesRequest{
		Market: "USHA", SortBegin: Num(0), SortCount: Num(20),
	})
	want := `{"market":"USHA","sort_begin":0,"sort_count":20}`
	if got != want {
		t.Fatalf("zero-valued optionals\n got %s\nwant %s", got, want)
	}
}

func TestBoolAndStringOptionals(t *testing.T) {
	got := encode(t, ListNewsItemsRequest{
		Security: Security{Code: "600519", Market: "USHA"},
		Category: 14339, Traversal: map[string]any{"mode": "page"},
		Advertorial: Flag(false), SummaryMode: Num(0),
	})
	want := `{"category":14339,"security":{"code":"600519","market":"USHA"},"traversal":{"mode":"page"},"summary_mode":0,"advertorial":false}`
	if got != want {
		t.Fatalf("optional scalars\n got %s\nwant %s", got, want)
	}
}

func TestSecurityNormalizesCase(t *testing.T) {
	got := encode(t, Security{Code: "600519", Market: "usha"})
	if want := `{"code":"600519","market":"USHA"}`; got != want {
		t.Fatalf("market not normalized: %s", got)
	}
}

func TestSecurityRejectsBadMarket(t *testing.T) {
	if _, err := json.Marshal(Security{Code: "600519", Market: "SH"}); err == nil {
		t.Fatalf("expected a 2-letter market to be rejected")
	}
}

func TestBlockAcceptsIDOrName(t *testing.T) {
	id := 61322
	if got, want := encode(t, Block{ID: &id}), `{"id":61322}`; got != want {
		t.Fatalf("block by id: %s", got)
	}
	if got, want := encode(t, Block{Name: "沪深300"}), `{"name":"沪深300"}`; got != want {
		t.Fatalf("block by name: %s", got)
	}
	if _, err := json.Marshal(Block{}); err == nil {
		t.Fatalf("empty block should be rejected")
	}
}

func TestSecurityListMarshalsAsArray(t *testing.T) {
	// list_security_futures_relations sends a bare array, not an object.
	securities, err := securityList([]Security{
		{Code: "600519", Market: "USHA"}, {Code: "300033", Market: "USZA"},
	})
	if err != nil {
		t.Fatalf("securityList: %v", err)
	}
	want := `[{"code":"600519","market":"USHA"},{"code":"300033","market":"USZA"}]`
	if got := encode(t, securities); got != want {
		t.Fatalf("securities array\n got %s\nwant %s", got, want)
	}
}

func TestEmptySecurityListRejected(t *testing.T) {
	if _, err := securityList(nil); err == nil {
		t.Fatalf("empty securities should be rejected before the ABI is touched")
	}
}

func TestGeneratedMethodValidationGuardsRequired(t *testing.T) {
	client := &Client{authenticated: true}
	// Security is required; the guard must fire without reaching the component.
	_, err := client.ListSecurityKLines(nil, ListSecurityKLinesRequest{
		Start: -3, End: 0, Adjust: AdjustNone, Period: PeriodDay,
	})
	if err == nil {
		t.Fatalf("expected the missing-security guard to fire")
	}
	var paramErr *ParameterError
	if !asParameterError(err, &paramErr) {
		t.Fatalf("expected *ParameterError, got %T: %v", err, err)
	}
	if paramErr.Argument != "security" {
		t.Fatalf("guard named %q, want security", paramErr.Argument)
	}
}

func asParameterError(err error, target **ParameterError) bool {
	typed, ok := err.(*ParameterError)
	if ok {
		*target = typed
	}
	return ok
}
