package thsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"
)

// PriceQuery mirrors the Python get_price/klines parameters.
//
// Set either Count or the Start/End pair; the two are mutually exclusive. With
// Count the window is relative ("the last N bars"), which is the only mode
// minute periods support.
type PriceQuery struct {
	Security  Security
	Start     time.Time
	End       time.Time
	Count     int
	Frequency string // "daily", "1m", "week", ... see ParsePeriod
	Fields    []string
	FQ        string // "", "none", "pre", "post"
}

// KLine is a projected candle row: the raw fields kept, minus any dropped by
// PriceQuery.Fields.
type PriceRow map[string]json.RawMessage

// GetPrice queries one instrument's KLine series.
func (c *Client) GetPrice(ctx context.Context, query PriceQuery) ([]PriceRow, error) {
	period, err := ParsePeriod(query.Frequency)
	if err != nil {
		return nil, err
	}
	adjust, err := ParseAdjust(query.FQ)
	if err != nil {
		return nil, err
	}

	start, end := -1, 0
	switch {
	case query.Count != 0 && (!query.Start.IsZero() || !query.End.IsZero()):
		return nil, invalidArgument("count", "不能与 start/end 同时使用")
	case query.Count < 0:
		return nil, invalidArgument("count", "必须大于 0")
	case query.Count > 0:
		start, end = -query.Count, 0
	case query.Start.IsZero() != query.End.IsZero():
		return nil, invalidArgument("start", "和 end 必须同时提供")
	case !query.Start.IsZero():
		if isMinutePeriod(period) {
			return nil, invalidArgument("frequency", "分钟 K 线的日期范围暂不支持；请改用 Count")
		}
		start, end = dateNumber(query.Start), dateNumber(query.End)
	}

	data, err := c.ListSecurityKLines(ctx, ListSecurityKLinesRequest{
		Security: query.Security,
		Start:    start,
		End:      end,
		Adjust:   adjust,
		Period:   period,
	})
	if err != nil {
		return nil, err
	}
	rows, err := asRecords(data)
	if err != nil {
		return nil, err
	}
	if query.Count > 0 && len(rows) > query.Count {
		rows = rows[len(rows)-query.Count:]
	}
	return project(rows, query.Fields)
}

// IntradayBars returns one trading day's minute timeline; a zero date means today.
func (c *Client) IntradayBars(ctx context.Context, security Security, date time.Time) ([]PriceRow, error) {
	data, err := c.ListSecurityIntradayBars(ctx, ListSecurityIntradayBarsRequest{
		Security: security,
		Date:     dateStringOrZero(date),
	})
	if err != nil {
		return nil, err
	}
	return asRecords(data)
}

// TimedIntradayBar attaches the wall-clock minute a bar belongs to.
type TimedIntradayBar struct {
	At  time.Time
	Bar IntradayBar
}

// IntradayBarsClock labels each minute bar with a Shanghai clock time.
//
// The component's `time` field is its internal minute index, not a timestamp:
// it advances by one per bar but also jumps mid-session, so it cannot be parsed
// directly. Bars are one per trading minute from 09:30, skipping the
// 11:30-13:00 lunch break and stopping at 15:00. Half-day sessions are not
// modelled — for exact event times use Ticks or GetSecurityTradingTimeline,
// whose `time` fields are real Unix seconds.
func (c *Client) IntradayBarsClock(ctx context.Context, security Security, date time.Time) ([]TimedIntradayBar, error) {
	data, err := c.ListSecurityIntradayBars(ctx, ListSecurityIntradayBarsRequest{
		Security: security,
		Date:     dateStringOrZero(date),
	})
	if err != nil {
		return nil, err
	}
	var bars []IntradayBar
	if err := data.Rows(&bars); err != nil {
		return nil, err
	}

	day := date
	if day.IsZero() {
		day = time.Now().In(Shanghai)
	}
	morning := time.Date(day.Year(), day.Month(), day.Day(), 9, 30, 0, 0, Shanghai)

	out := make([]TimedIntradayBar, 0, len(bars))
	for index, bar := range bars {
		out = append(out, TimedIntradayBar{At: tradingMinute(morning, index), Bar: bar})
	}
	return out, nil
}

const morningMinutes = 120 // 09:30–11:30

// tradingMinute maps a zero-based bar index to its clock minute.
func tradingMinute(sessionStart time.Time, index int) time.Time {
	if index < morningMinutes {
		return sessionStart.Add(time.Duration(index) * time.Minute)
	}
	afternoon := time.Date(sessionStart.Year(), sessionStart.Month(), sessionStart.Day(),
		13, 0, 0, 0, sessionStart.Location())
	return afternoon.Add(time.Duration(index-morningMinutes) * time.Minute)
}

// TickQuery selects either a recent-count window or an explicit time range.
type TickQuery struct {
	Security Security
	Count    int
	Start    time.Time
	End      time.Time
}

// Ticks returns Level-1 trade prints.
func (c *Client) Ticks(ctx context.Context, query TickQuery) ([]PriceRow, error) {
	if query.Count != 0 && (!query.Start.IsZero() || !query.End.IsZero()) {
		return nil, invalidArgument("count", "不能与 Start/End 同时使用")
	}
	if query.Count < 0 {
		return nil, invalidArgument("count", "必须大于 0")
	}

	var window any
	switch {
	case query.Count > 0:
		window = map[string]any{"count": query.Count}
	case !query.Start.IsZero() && !query.End.IsZero():
		window = map[string]any{"start": formatTimestamp(query.Start), "end": formatTimestamp(query.End)}
	case query.Start.IsZero() && query.End.IsZero():
		return nil, invalidArgument("count", "请提供 Count，或同时提供 Start 和 End")
	default:
		return nil, invalidArgument("start", "和 end 必须同时提供")
	}

	data, err := c.ListSecurityTicks(ctx, ListSecurityTicksRequest{
		Security: query.Security,
		Window:   window,
	})
	if err != nil {
		return nil, err
	}
	rows, err := asRecords(data)
	if err != nil {
		return nil, err
	}
	if query.Count > 0 && len(rows) > query.Count {
		rows = rows[len(rows)-query.Count:]
	}
	return rows, nil
}

// DepthRow is one price level of the expanded five-ladder table.
type DepthRow struct {
	FullCode string  `json:"full_code"`
	Side     string  `json:"side"`
	Level    int     `json:"level"`
	Price    float64 `json:"price"`
	Volume   int64   `json:"volume"`
}

// Depth reads the five-level book of one or more instruments and flattens it by
// side and level, like the Python depth() helper.
func (c *Client) Depth(ctx context.Context, securities ...Security) ([]DepthRow, error) {
	list, err := securityList(securities)
	if err != nil {
		return nil, err
	}
	data, err := c.ListSecurityOrderBooks(ctx, ListSecurityOrderBooksRequest{
		Securities: list,
		Depth:      DepthFive,
	})
	if err != nil {
		return nil, err
	}
	var books []OrderBook
	if err := data.Rows(&books); err != nil {
		return nil, err
	}
	rows := make([]DepthRow, 0, len(books)*10)
	for _, book := range books {
		for _, side := range []struct {
			name   string
			quotes []Quote
		}{{"bid", book.Bids}, {"ask", book.Asks}} {
			for level, quote := range side.quotes {
				rows = append(rows, DepthRow{
					FullCode: book.Security.FullCode(),
					Side:     side.name,
					Level:    level + 1,
					Price:    quote.Price,
					Volume:   quote.Volume,
				})
			}
		}
	}
	return rows, nil
}

// CorporateActions returns the most recent dividend and split records.
func (c *Client) CorporateActions(ctx context.Context, security Security, count int) ([]PriceRow, error) {
	if count <= 0 {
		count = 100
	}
	data, err := c.ListSecurityCorporateActions(ctx, ListSecurityCorporateActionsRequest{
		Security: security,
		Start:    -count,
		End:      0,
		Adjust:   AdjustNone,
		Period:   PeriodDay,
	})
	if err != nil {
		return nil, err
	}
	rows, err := asRecords(data)
	if err != nil {
		return nil, err
	}
	if len(rows) > count {
		rows = rows[len(rows)-count:]
	}
	return rows, nil
}

// SearchSymbols matches code, pinyin, or Chinese name.
func (c *Client) SearchSymbols(ctx context.Context, pattern string, market string) ([]PriceRow, error) {
	pattern = strings.TrimSpace(pattern)
	if pattern == "" {
		return nil, invalidArgument("pattern", "不能为空")
	}
	req := SearchSecuritiesRequest{Pattern: pattern}
	if market != "" {
		parsed, err := ParseMarket(market)
		if err != nil {
			return nil, err
		}
		req.Market = Text(parsed)
	}
	data, err := c.SearchSecurities(ctx, req)
	if err != nil {
		return nil, err
	}
	return identityRows(data)
}

// CompleteTHSCodes resolves short codes into full securities.
func (c *Client) CompleteTHSCodes(ctx context.Context, codes ...string) ([]PriceRow, error) {
	var cleaned []string
	for _, code := range codes {
		if code = strings.TrimSpace(strings.ToUpper(code)); code != "" {
			cleaned = append(cleaned, code)
		}
	}
	if len(cleaned) == 0 {
		return nil, invalidArgument("ths_code", "不能为空")
	}
	data, err := c.ResolveSecurities(ctx, ResolveSecuritiesRequest{Codes: cleaned})
	if err != nil {
		return nil, err
	}
	return identityRows(data)
}

// WencaiSelection is one row of a natural-language stock screen, with the
// dynamic fields flattened into columns keyed by display name.
type WencaiSelection struct {
	Rows   []map[string]any
	Query  string
	Fields []map[string]any
}

// Wencai runs a natural-language screen, reproducing the Python wencai_nlp
// identity merge and dynamic-column expansion.
func (c *Client) Wencai(ctx context.Context, condition string, limit int, markets ...string) (*WencaiSelection, error) {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return nil, invalidArgument("condition", "不能为空")
	}
	if limit <= 0 {
		limit = 100
	}
	req := QueryWencaiRequest{Query: condition, Limit: Num(limit)}
	if len(markets) > 0 {
		parsed, err := ParseMarkets(markets...)
		if err != nil {
			return nil, err
		}
		req.Markets = parsed
	}

	data, err := c.QueryWencai(ctx, req)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Query      string           `json:"query"`
		Fields     []map[string]any `json:"fields"`
		Securities []map[string]any `json:"securities"`
		Rows       []struct {
			Security map[string]any   `json:"security"`
			FullCode string           `json:"full_code"`
			Code     string           `json:"code"`
			Market   string           `json:"market"`
			Cells    []map[string]any `json:"cells"`
		} `json:"rows"`
	}
	if err := data.Into(&envelope); err != nil {
		return nil, err
	}

	fieldNames := map[string]string{}
	for _, field := range envelope.Fields {
		name := firstNonEmptyString(field, "display_name", "show_name", "index_name", "query_key")
		if name != "" {
			fieldNames[cellKey(field)] = name
		}
	}

	selection := &WencaiSelection{Query: envelope.Query, Fields: envelope.Fields}
	if selection.Query == "" {
		selection.Query = condition
	}
	seen := map[string]bool{}
	for index, raw := range envelope.Rows {
		code := strings.TrimSpace(raw.Code)
		market := strings.ToUpper(strings.TrimSpace(raw.Market))
		fullCode := strings.ToUpper(strings.TrimSpace(raw.FullCode))
		for _, cell := range raw.Cells {
			if firstNonEmptyString(cell, "name") == "code" {
				code = fmt.Sprint(cell["value"])
			}
		}
		if fullCode == "" && market != "" && code != "" {
			fullCode = market + code
		}
		if fullCode == "" && index < len(envelope.Securities) {
			fullCode = strings.ToUpper(firstNonEmptyString(envelope.Securities[index], "full_code"))
		}
		if fullCode != "" && len(fullCode) > 4 && marketPattern.MatchString(fullCode[:4]) {
			market, code = fullCode[:4], fullCode[4:]
		}

		row := map[string]any{}
		if fullCode != "" {
			row["full_code"] = fullCode
		}
		if code != "" {
			row["code"] = code
		}
		if market != "" {
			row["market"] = market
		}
		for _, cell := range raw.Cells {
			name := firstNonEmptyString(cell, "name")
			if column, ok := fieldNames[cellKey(cell)]; ok {
				name = column
			}
			if name == "" || name == "code" {
				continue
			}
			row[name] = cell["value"]
		}
		if len(row) == 0 {
			continue
		}
		key := fullCode + "\x00" + code
		if key != "\x00" && seen[key] {
			continue
		}
		seen[key] = true
		selection.Rows = append(selection.Rows, row)
	}

	// Some screens return only a security pool with no queryable fields; keep
	// those identities so the result can feed the other APIs.
	for _, pooled := range envelope.Securities {
		fullCode := strings.ToUpper(firstNonEmptyString(pooled, "full_code"))
		code := firstNonEmptyString(pooled, "code")
		if key := fullCode + "\x00" + code; key != "\x00" && !seen[key] {
			seen[key] = true
			row := map[string]any{}
			if fullCode != "" {
				row["full_code"] = fullCode
			}
			if code != "" {
				row["code"] = code
			}
			selection.Rows = append(selection.Rows, row)
		}
	}
	return selection, nil
}

// NewsResult is one page of the 7x24 flash feed plus its paging metadata.
type NewsResult struct {
	Items    []PriceRow
	Metadata map[string]json.RawMessage
}

// News reads a page of flash news.
func (c *Client) News(ctx context.Context, page int, pageSize int, tag string) (*NewsResult, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, invalidArgument("page", "和 pageSize 必须大于 0")
	}
	data, err := c.ListNewsFlashItems(ctx, ListNewsFlashItemsRequest{
		Traversal: FlashPage, Page: Num(page), PageSize: Num(pageSize), Tag: Text(tag),
	})
	if err != nil {
		return nil, err
	}
	var items []PriceRow
	if err := data.Rows(&items); err != nil {
		return nil, err
	}
	metadata, err := data.Metadata()
	if err != nil {
		return nil, err
	}
	return &NewsResult{Items: items, Metadata: metadata}, nil
}

// MarketSecurities lists one market with paging cursors in Metadata.
func (c *Client) MarketSecurities(ctx context.Context, market string, start, count int) (*NewsResult, error) {
	if start < 0 || count <= 0 {
		return nil, invalidArgument("start", "不能小于 0，count 必须大于 0")
	}
	parsed, err := ParseMarket(market)
	if err != nil {
		return nil, err
	}
	data, err := c.ListMarketSecurities(ctx, ListMarketSecuritiesRequest{
		Market: parsed, SortBegin: Num(start), SortCount: Num(count),
	})
	if err != nil {
		return nil, err
	}
	items, err := identityRows(data)
	if err != nil {
		return nil, err
	}
	metadata, err := data.Metadata()
	if err != nil {
		return nil, err
	}
	return &NewsResult{Items: items, Metadata: metadata}, nil
}

// BlockConstituents lists a board's members by numeric id or exact name.
func (c *Client) BlockConstituents(ctx context.Context, block string, start, count int) ([]PriceRow, error) {
	if start < 0 || count <= 0 {
		return nil, invalidArgument("start", "不能小于 0，count 必须大于 0")
	}
	identity := Block{}
	if id, err := parseBlockID(block); err == nil {
		identity.ID = &id
	} else {
		resolved, err := c.ResolveBlock(ctx, ResolveBlockRequest{Name: block})
		if err != nil {
			return nil, err
		}
		if err := resolved.Into(&identity); err != nil {
			return nil, err
		}
	}

	data, err := c.ListBlockConstituents(ctx, ListBlockConstituentsRequest{
		Block: identity, SortBegin: Num(start), SortCount: Num(count), SortOrder: Text("D"),
	})
	if err != nil {
		return nil, err
	}
	return identityRows(data)
}

func parseBlockID(value string) (int, error) {
	var id int
	if _, err := fmt.Sscanf(strings.TrimSpace(value), "%d", &id); err != nil {
		return 0, err
	}
	if id <= 0 || strings.TrimSpace(value) != fmt.Sprint(id) {
		return 0, invalidArgument("block", "板块 ID 必须大于 0")
	}
	return id, nil
}

func isMinutePeriod(period Period) bool { return period >= Period1Minute && period <= Period120Minutes }

func dateStringOrZero(value time.Time) string {
	if value.IsZero() {
		return "0"
	}
	return dateString(value)
}

// asRecords decodes a row array, keeping every field uninterpreted so callers
// can read whatever the account's permissions returned.
func asRecords(data RawData) ([]PriceRow, error) {
	var rows []PriceRow
	if err := data.Rows(&rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// identityRows flattens the nested `security` object into code/market/full_code
// columns, matching the Python _identity_rows helper.
func identityRows(data RawData) ([]PriceRow, error) {
	rows, err := asRecords(data)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		raw, ok := row["security"]
		if !ok {
			continue
		}
		var security Security
		if err := json.Unmarshal(raw, &security); err != nil {
			continue
		}
		delete(row, "security")
		if _, exists := row["code"]; !exists {
			row["code"] = jsonString(security.Code)
		}
		if _, exists := row["market"]; !exists {
			row["market"] = jsonString(security.Market)
		}
		if _, exists := row["full_code"]; !exists {
			row["full_code"] = jsonString(security.FullCode())
		}
	}
	return rows, nil
}

// project keeps only the named columns and errors when one is missing.
func project(rows []PriceRow, fields []string) ([]PriceRow, error) {
	if len(fields) == 0 {
		return rows, nil
	}
	for _, row := range rows {
		for _, field := range fields {
			if _, ok := row[field]; !ok {
				available := make([]string, 0, len(row))
				for key := range row {
					available = append(available, key)
				}
				sort.Strings(available)
				return nil, invalidArgument("fields", "返回结果缺少字段 %s（可用: %s）",
					field, strings.Join(available, ", "))
			}
		}
	}
	out := make([]PriceRow, 0, len(rows))
	for _, row := range rows {
		projected := PriceRow{}
		for _, field := range fields {
			projected[field] = row[field]
		}
		out = append(out, projected)
	}
	return out, nil
}

func jsonString(value string) json.RawMessage {
	encoded, _ := json.Marshal(value)
	return encoded
}

func firstNonEmptyString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if text, ok := values[key].(string); ok && strings.TrimSpace(text) != "" {
			return strings.TrimSpace(text)
		}
	}
	return ""
}

// cellKey groups a wencai cell with the field definition that names it.
func cellKey(values map[string]any) string {
	index := fmt.Sprint(values["key_index"])
	if _, ok := values["key_index"]; !ok {
		index = fmt.Sprint(values["index"])
	}
	timestamp, _ := values["timestamp"].(string)
	return index + "\x00" + timestamp
}
