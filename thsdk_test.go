package thsdk_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
)

// TestMain stages the runtime component in the working directory so a checkout
// with native/hq.* present can run without extra configuration.
func TestMain(m *testing.M) {
	if os.Getenv("THSDK_RUNTIME_PATH") == "" {
		for _, name := range []string{"hq.dll", "hq.so", "hq.dylib"} {
			if source, err := os.Stat(filepath.Join("native", name)); err == nil && !source.IsDir() {
				target := filepath.Join(os.TempDir(), "go-thsdk-"+name)
				if err := copyFile(filepath.Join("native", name), target); err == nil {
					os.Setenv("THSDK_RUNTIME_PATH", target)
				}
				break
			}
		}
	}
	os.Exit(m.Run())
}

func copyFile(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(target, data, 0o600)
}

// The live suite runs only when a runtime component is reachable; CI machines
// without the native library get a skip rather than a false failure.
func newClient(t *testing.T) (*thsdk.Client, context.Context) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping live service calls under -short")
	}
	client, err := thsdk.New(thsdk.Options{})
	if err != nil {
		t.Skipf("runtime component unavailable: %v", err)
	}
	t.Cleanup(func() { _ = client.Close() })
	if err := client.Auth(context.Background(), thsdk.Credentials{}); err != nil {
		t.Skipf("authentication unavailable: %v", err)
	}
	return client, context.Background()
}

func TestGetPriceRespectsCountAndFields(t *testing.T) {
	client, ctx := newClient(t)
	rows, err := client.GetPrice(ctx, thsdk.PriceQuery{
		Security:  thsdk.MustSecurity("USHA600519"),
		Count:     5,
		Frequency: "daily",
		FQ:        "pre",
		Fields:    []string{"time", "close"},
	})
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("GetPrice: %v", err)
	}
	if len(rows) != 5 {
		t.Fatalf("count=5 produced %d rows", len(rows))
	}
	for _, row := range rows {
		if len(row) != 2 {
			t.Fatalf("field projection kept %d columns: %v", len(row), row)
		}
		if _, ok := row["close"]; !ok {
			t.Fatalf("close column missing: %v", row)
		}
	}
}

func TestGetPriceRejectsUnknownField(t *testing.T) {
	client, ctx := newClient(t)
	_, err := client.GetPrice(ctx, thsdk.PriceQuery{
		Security: thsdk.MustSecurity("USHA600519"), Count: 2, Frequency: "daily",
		Fields: []string{"does_not_exist"},
	})
	if err == nil {
		t.Fatalf("expected an unknown-field error")
	}
	// Projection happens after the fetch, so a service-side failure here is not
	// what the test is about.
	var paramErr *thsdk.ParameterError
	if !asParameter(err, &paramErr) {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("expected *ParameterError, got %T: %v", err, err)
	}
}

func TestDepthExpandsLevels(t *testing.T) {
	client, ctx := newClient(t)
	rows, err := client.Depth(ctx, thsdk.MustSecurity("USHA600519"))
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("Depth: %v", err)
	}
	if len(rows) == 0 {
		t.Fatalf("Depth returned no levels")
	}
	bids, asks := 0, 0
	for _, row := range rows {
		switch row.Side {
		case "bid":
			bids++
		case "ask":
			asks++
		default:
			t.Fatalf("unexpected side %q", row.Side)
		}
		if row.Level < 1 {
			t.Fatalf("levels start at 1, got %d", row.Level)
		}
	}
	if bids == 0 || asks == 0 {
		t.Fatalf("expected both sides, got bids=%d asks=%d", bids, asks)
	}
}

func TestTicksAndIntradayReturnRows(t *testing.T) {
	client, ctx := newClient(t)
	ticks, err := client.Ticks(ctx, thsdk.TickQuery{
		Security: thsdk.MustSecurity("USHA600519"), Count: 10,
	})
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("Ticks: %v", err)
	}
	if len(ticks) == 0 || len(ticks) > 10 {
		t.Fatalf("count=10 produced %d ticks", len(ticks))
	}

	bars, err := client.IntradayBars(ctx, thsdk.MustSecurity("USHA600519"), time.Time{})
	if err != nil {
		skipWhenServiceUnavailable(t, err)
		t.Fatalf("IntradayBars: %v", err)
	}
	if len(bars) == 0 {
		t.Fatalf("IntradayBars returned nothing")
	}
}

func TestWencaiFlattensDynamicColumns(t *testing.T) {
	client, ctx := newClient(t)
	result, err := client.Wencai(ctx, "市盈率小于20", 5)
	if err != nil {
		t.Fatalf("Wencai: %v", err)
	}
	if result.Query == "" {
		t.Fatalf("query echo missing")
	}
	if len(result.Rows) == 0 {
		t.Fatalf("wencai produced no rows")
	}
	for _, row := range result.Rows {
		if len(row) == 0 {
			t.Fatalf("row has no dynamic columns: %v", row)
		}
	}
}

func TestNewsPageCarriesMetadata(t *testing.T) {
	client, ctx := newClient(t)
	page, err := client.News(ctx, 1, 5, "")
	if err != nil {
		t.Fatalf("News: %v", err)
	}
	if len(page.Items) == 0 {
		t.Fatalf("no flash items")
	}
	if len(page.Metadata) == 0 {
		t.Fatalf("paging metadata dropped")
	}
}

func TestStandardMethodRequiredFieldGuard(t *testing.T) {
	client, ctx := newClient(t)
	_, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
		Start: -3, End: 0, Adjust: thsdk.AdjustNone, Period: thsdk.PeriodDay,
	})
	if err == nil {
		t.Fatalf("expected a missing-security error")
	}
	if !asParameter(err, new(*thsdk.ParameterError)) {
		t.Fatalf("expected *ParameterError, got %T: %v", err, err)
	}
}

func TestNotAuthenticatedGuard(t *testing.T) {
	client, err := thsdk.New(thsdk.Options{})
	if err != nil {
		t.Skipf("runtime component unavailable: %v", err)
	}
	defer client.Close()
	if _, err := client.AccountPermissions(context.Background()); err != thsdk.ErrNotAuthenticated {
		t.Fatalf("expected ErrNotAuthenticated, got %v", err)
	}
}

func TestContextDeadlineReachesComponent(t *testing.T) {
	client, ctx := newClient(t)
	short, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()
	_, err := client.ListMarketSecurities(short, thsdk.ListMarketSecuritiesRequest{
		Market: "USHA", SortBegin: thsdk.Num(0), SortCount: thsdk.Num(20),
	})
	if err == nil {
		t.Fatalf("expected the short deadline to fail")
	}
}

func TestRawDataLiftsRowsFromEitherShape(t *testing.T) {
	var into thsdk.RawData
	if err := json.Unmarshal([]byte(`[{"time":20261029,"close":1.5}]`), &into); err != nil {
		t.Fatalf("unmarshal fixture: %v", err)
	}
	var array []thsdk.KLine
	if err := into.Rows(&array); err != nil {
		t.Fatalf("bare array: %v", err)
	}
	if len(array) != 1 || array[0].Close != 1.5 {
		t.Fatalf("bare array decoded to %+v", array)
	}

	var envelope thsdk.RawData
	if err := json.Unmarshal([]byte(`{"items":[{"name":"a"}],"total":7}`), &envelope); err != nil {
		t.Fatalf("unmarshal envelope fixture: %v", err)
	}
	var items []struct {
		Name string `json:"name"`
	}
	if err := envelope.Rows(&items); err != nil {
		t.Fatalf("envelope rows: %v", err)
	}
	if len(items) != 1 || items[0].Name != "a" {
		t.Fatalf("envelope decoded to %+v", items)
	}
	metadata, err := envelope.Metadata()
	if err != nil {
		t.Fatalf("metadata: %v", err)
	}
	if _, ok := metadata["items"]; ok {
		t.Fatalf("metadata should exclude the lifted row field")
	}
	if _, ok := metadata["total"]; !ok {
		t.Fatalf("metadata lost the paging cursor: %v", metadata)
	}
}

func asParameter(err error, target **thsdk.ParameterError) bool {
	typed, ok := err.(*thsdk.ParameterError)
	if ok {
		*target = typed
	}
	return ok
}
