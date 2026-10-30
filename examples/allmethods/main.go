package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
)

// This tool sends one live request per native method and classifies the answer:
//
//	OK       the component returned data
//	SERVICE  no_data / encode_failed / no_permission / permission_denied — the
//	         upstream service, not the SDK. The Python SDK fails identically.
//	FAIL     a real problem worth investigating
//
// Use `go test ./... -run Live` for the CI-friendly form of this sweep.
var includeWrites = flag.Bool("writes", false, "also exercise the 10 account-write methods (mutates the shared trial account)")

func main() {
	flag.Parse()

	client, err := thsdk.New(thsdk.Options{})
	if err != nil {
		fmt.Fprintln(os.Stderr, "open:", err)
		os.Exit(1)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Auth(ctx, thsdk.Credentials{}); err != nil {
		fmt.Fprintln(os.Stderr, "auth:", err)
		os.Exit(1)
	}

	mout, _ := json.Marshal(thsdk.MustSecurity("USHA600519"))
	fmt.Println("component handshake + auth OK; security fixture:", string(mout))
	fmt.Println()

	sec := thsdk.MustSecurity("USHA600519")
	secSZ := thsdk.MustSecurity("USZA300033")
	secs := []thsdk.Security{sec, secSZ}
	blockID := 100
	wencaiField := map[string]any{
		"query_group": "换手率", "query_key": "换手率", "index_name": "换手率",
		"show_name": "换手率", "timestamp": "20261029", "index": 2, "unit": "%", "type": "DOUBLE",
	}
	indexRef := map[string]any{"security": sec}

	// One live request per read-only native method.
	cases := []struct {
		method string
		params any
	}{
		{"account_permissions", nil},
		{"analyze_security_limit_up", map[string]any{"security": sec}},
		{"calculate_security_realtime_statistics", map[string]any{"securities": secs, "metrics": []int{330323, 330329}}},
		{"check_market_timeline_readiness", map[string]any{"market": "USHA", "time_index_day": "0"}},
		{"get_account_watchlist", nil},
		{"get_account_watchlist_groups", nil},
		{"get_market_metadata", map[string]any{"market": "USHA"}},
		{"get_market_security_names", map[string]any{"market": "USHA"}},
		{"get_security_concept_tags", map[string]any{"security": sec}},
		{"get_security_industry", map[string]any{"security": sec}},
		{"get_security_popularity_rank", map[string]any{"security": sec}},
		{"get_security_price_limit_events", map[string]any{"security": sec, "start_year": 2025, "end_year": 2026}},
		{"get_security_short_term_highlights", map[string]any{"security": sec}},
		{"get_security_trading_timeline", map[string]any{"security": sec}},
		{"list_block_constituents", map[string]any{"block": map[string]any{"id": blockID}, "sort_begin": 0, "sort_count": 5, "sort_order": "D"}},
		{"list_block_descriptions", map[string]any{"blocks": []map[string]any{{"id": blockID}}}},
		{"list_commodity_stock_linkage_news", map[string]any{"commodity_index": map[string]any{"id": "S004161736"}}},
		{"list_futures_related_securities", nil},
		{"list_hot_block_calendar_events", nil},
		{"list_hot_event_news", map[string]any{"block": map[string]any{"id": blockID}, "range": "recent"}},
		{"list_index_call_auction_quotes", map[string]any{"index": indexRef, "phase": "opening"}},
		{"list_industry_children", map[string]any{"industry": map[string]any{"security": map[string]string{"code": "881273", "market": "URFI"}}}},
		{"list_market_call_auction_unusual_events", map[string]any{"market": "USHA", "trade_date": "20261029"}},
		{"list_market_change_events", nil},
		{"list_market_securities", map[string]any{"market": "USHA", "sort_begin": 0, "sort_count": 5}},
		{"list_market_short_term_events", map[string]any{"market": "USHA", "mode": "security_scoped", "securities": []thsdk.Security{sec}, "max_count": 10}},
		{"list_news_flash_items", map[string]any{"traversal": "page", "page": 1, "page_size": 5}},
		{"list_news_items", map[string]any{"category": 14339, "security": sec, "traversal": map[string]any{"mode": "page"}, "summary_mode": 2, "advertorial": false}},
		{"list_related_security_performances", map[string]any{"security": sec}},
		{"list_securities_by_realtime_signal", map[string]any{"markets": []string{"USHA", "USZA"}, "signal": "limit_up"}},
		{"list_security_ah_relations", map[string]any{"security": sec}},
		{"list_security_block_memberships", map[string]any{"security": sec}},
		{"list_security_call_auction_quotes", map[string]any{"security": secSZ, "phase": "opening", "date": "2026-09-03T09:15:00+08:00"}},
		{"list_security_corporate_actions", map[string]any{"security": sec, "start": -3, "end": 0, "adjust": 0, "period": 0x4000}},
		{"list_security_daily_capital_flows", map[string]any{"securities": secs}},
		{"list_security_daily_k_lines_with_previous_close", map[string]any{"security": sec, "start": -3, "end": 0, "adjust": 1}},
		{"list_security_extended_hours_intraday_bars", map[string]any{"security": sec, "session": "pre_market", "date": "0"}},
		{"list_security_extended_hours_ticks", map[string]any{"security": sec, "session": "pre_market", "window": map[string]any{"count": 5}}},
		{"list_security_financial_snapshots", map[string]any{"securities": secs, "fields": []string{"total_shares"}}},
		{"list_security_futures_relations", secs},
		{"list_security_industry_mappings", map[string]any{"securities": secs}},
		{"list_security_intraday_bars", map[string]any{"security": sec, "date": "0"}},
		{"list_security_k_lines", map[string]any{"security": sec, "start": -3, "end": 0, "adjust": 1, "period": 0x4000}},
		{"list_security_link_relations", "USHA600519"},
		{"list_security_news_events", map[string]any{"security": sec, "timeline": "historical", "window": map[string]any{"start": "2026-08-01T00:00:00+08:00", "end": "2026-09-02T23:59:59+08:00"}}},
		{"list_security_news_markers", map[string]any{"security": sec, "timeline": "historical"}},
		{"list_security_option_greeks", map[string]any{"securities": secs, "metrics": []int{723854}}},
		{"list_security_order_books", map[string]any{"securities": secs, "depth": "5"}},
		{"list_security_price_volume_levels", map[string]any{"security": sec}},
		{"list_security_research_reports", map[string]any{"security": sec, "traversal": map[string]any{"mode": "page"}}},
		{"list_security_ticks", map[string]any{"security": sec, "window": map[string]any{"count": 5}}},
		{"list_security_time_and_sales_ticks", map[string]any{"security": sec, "count": 5}},
		{"list_wencai_expression_fields", "市盈率"},
		{"list_wencai_hot_blocks", map[string]any{"kind": "concept"}},
		{"query_wencai", map[string]any{"query": "市盈率小于20", "limit": 3}},
		{"query_wencai_realtime_fields", map[string]any{"securities": secs, "fields": []map[string]any{wencaiField}}},
		{"query_wencai_securities", map[string]any{"query": "市值大于1000亿", "limit": 3}},
		{"rank_block_securities", map[string]any{"block": map[string]any{"id": blockID}, "sort_begin": 0, "sort_count": 5}},
		{"rank_block_securities_by_industry", map[string]any{"block": map[string]any{"id": blockID}, "sort_begin": 0, "sort_count": 5}},
		{"rank_block_securities_by_wencai_field", map[string]any{"block": map[string]any{"id": 61322}, "field": wencaiField, "offset": 0, "limit": 10, "sort_order": "D"}},
		{"rank_constituents_by_performance_contribution", map[string]any{"code": "000300", "market": "USHA", "sort_begin": 0, "sort_count": 5}},
		{"rank_related_securities", map[string]any{"security": sec, "sort_begin": 0, "sort_count": 5}},
		{"rank_securities_by_popularity", map[string]any{"type": 1}},
		{"rank_securities_by_realtime_statistic", map[string]any{"markets": []string{"USHA"}, "sort_by": 330329, "sort_begin": 0, "sort_count": 5}},
		{"resolve_block", map[string]any{"name": "沪深300"}},
		{"resolve_securities", map[string]any{"codes": []string{"600519", "300033"}}},
		{"resolve_wencai_securities", map[string]any{"query": "高送转", "limit": 3}},
		{"search_securities", map[string]any{"pattern": "茅台"}},
		{"sort_securities", map[string]any{"securities": secs, "sort_begin": 0, "sort_count": 5}},
		{"sort_securities_by_industry", map[string]any{"securities": secs, "sort_begin": 0, "sort_count": 5}},
	}

	writes := []struct {
		method string
		params any
	}{
		{"create_account_watchlist_group", map[string]any{"name": "gothsdk-tmp", "securities": secs}},
		{"add_account_watchlist_securities", map[string]any{"securities": secs}},
		{"add_account_watchlist_group_securities", map[string]any{"group_id": blockID, "securities": secs}},
		{"remove_account_watchlist_securities", map[string]any{"securities": secs}},
		{"remove_account_watchlist_group_securities", map[string]any{"group_id": blockID, "securities": secs}},
		{"replace_account_watchlist_securities", map[string]any{"version": 0, "securities": secs}},
		{"replace_account_watchlist_group_securities", map[string]any{"group_id": blockID, "securities": secs}},
		{"rename_account_watchlist_group", map[string]any{"group_id": blockID, "name": "gothsdk-renamed"}},
		{"delete_account_watchlist_group", map[string]any{"group_id": blockID}},
		{"clear_account_watchlist", nil},
	}

	if *includeWrites {
		cases = append(cases, writes...)
		fmt.Println("!! account-write methods included")
	} else {
		fmt.Printf("skipped %d account-write methods (pass -writes to include)\n\n", len(writes))
	}

	var ok, service int
	serviceCodes := map[string]int{}
	failures := []string{}

	for _, item := range cases {
		start := time.Now()
		data, err := client.CallRaw(ctx, item.method, item.params)
		elapsed := time.Since(start)

		if err != nil {
			var apiErr *thsdk.Error
			if asError(err, &apiErr) {
				switch apiErr.Code {
				case "no_data", "encode_failed", "no_permission", "permission_denied":
					service++
					serviceCodes[apiErr.Code]++
					fmt.Printf("SERVICE  %-52s %-16s %s\n", item.method, apiErr.Code, truncate(apiErr.Message, 40))
					continue
				}
			}
			failures = append(failures, fmt.Sprintf("%s -> %v", item.method, err))
			fmt.Printf("FAIL     %-52s %v\n", item.method, err)
			continue
		}

		summary := describe(data.Bytes())
		ok++
		fmt.Printf("OK       %-52s %6dms  %s\n", item.method, elapsed.Milliseconds(), summary)
	}

	fmt.Printf("\n=== %d requested: %d ok, %d service-side (%v), %d real failures ===\n",
		len(cases), ok, service, sortedCounts(serviceCodes), len(failures))
	for _, f := range failures {
		fmt.Println("  -", f)
	}
	if len(failures) > 0 {
		os.Exit(1)
	}
}

func describe(payload []byte) string {
	if len(payload) == 0 || string(payload) == "null" {
		return "null"
	}
	var array []map[string]json.RawMessage
	if json.Unmarshal(payload, &array) == nil {
		return fmt.Sprintf("array[%d] %s", len(array), keysOf(array, 0))
	}
	var object map[string]json.RawMessage
	if json.Unmarshal(payload, &object) != nil {
		return fmt.Sprintf("bytes[%d]", len(payload))
	}
	names := make([]string, 0, len(object))
	for key := range object {
		names = append(names, key)
	}
	sort.Strings(names)
	for _, rowField := range []string{"rows", "items", "securities", "events", "tags", "sections", "groups"} {
		if value, exists := object[rowField]; exists {
			var rows []map[string]json.RawMessage
			if json.Unmarshal(value, &rows) == nil {
				return fmt.Sprintf("object{%s: array[%d]} %s", rowField, len(rows), keysOf(rows, 0))
			}
		}
	}
	return fmt.Sprintf("object{%s}", strings.Join(names, ","))
}

func keysOf(rows []map[string]json.RawMessage, index int) string {
	if index >= len(rows) {
		return ""
	}
	names := make([]string, 0, len(rows[index]))
	for key := range rows[index] {
		names = append(names, key)
	}
	if len(names) == 0 {
		return ""
	}
	sort.Strings(names)
	if len(names) > 6 {
		names = append(names[:6:6], "...")
	}
	return "keys=" + strings.Join(names, ",")
}

func sortedCounts(counts map[string]int) string {
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s:%d", key, counts[key]))
	}
	return strings.Join(parts, " ")
}

func truncate(text string, limit int) string {
	if len(text) <= limit {
		return text
	}
	return text[:limit]
}

func asError(err error, target **thsdk.Error) bool {
	typed, ok := err.(*thsdk.Error)
	if ok {
		*target = typed
	}
	return ok
}
