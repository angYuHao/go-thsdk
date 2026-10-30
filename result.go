package thsdk

import (
	"encoding/json"
	"fmt"
	"time"
)

// RawData is the untouched `data` payload of one runtime call.
// Call Into for the whole envelope, or Rows for the record set inside it.
type RawData struct {
	payload json.RawMessage
}

// Bytes exposes the raw JSON, for logging or custom decoding.
func (d RawData) Bytes() []byte { return d.payload }

// String renders the raw JSON.
func (d RawData) String() string { return string(d.payload) }

// IsEmpty reports a null or absent payload.
func (d RawData) IsEmpty() bool {
	return len(d.payload) == 0 || string(d.payload) == "null"
}

// UnmarshalJSON lets RawData nest inside another decoded structure.
func (d *RawData) UnmarshalJSON(data []byte) error {
	d.payload = append(d.payload[:0], data...)
	return nil
}

// MarshalJSON re-emits the stored payload.
func (d RawData) MarshalJSON() ([]byte, error) {
	if d.IsEmpty() {
		return []byte("null"), nil
	}
	return d.payload, nil
}

// Into decodes the entire payload into out.
func (d RawData) Into(out any) error {
	if d.IsEmpty() || out == nil {
		return nil
	}
	if err := json.Unmarshal(d.payload, out); err != nil {
		return fmt.Errorf("thsdk: 解析结果失败: %w", err)
	}
	return nil
}

// Rows decodes the record set: the payload itself when it is an array,
// otherwise the first array-valued field of the envelope object.
//
// This mirrors the Python SDK, which lifts the manifest's rows field into a
// DataFrame and keeps the remaining keys as metadata.
func (d RawData) Rows(out any) error {
	if d.IsEmpty() {
		return nil
	}
	if d.payload[0] == '[' {
		return d.Into(out)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(d.payload, &envelope); err != nil {
		return fmt.Errorf("thsdk: 结果不是数组或对象: %s", first(d.payload, 120))
	}
	for _, key := range rowsFieldOrder {
		if value, ok := envelope[key]; ok {
			return json.Unmarshal(value, out)
		}
	}
	return fmt.Errorf("thsdk: 结果里没有数组字段，可用字段: %v", mapKeys(envelope))
}

// Metadata returns the envelope keys other than the lifted row set, so a caller
// can read paging cursors alongside the rows.
func (d RawData) Metadata() (map[string]json.RawMessage, error) {
	if d.IsEmpty() {
		return nil, nil
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(d.payload, &envelope); err != nil {
		return nil, nil // a bare array carries no metadata
	}
	for _, key := range rowsFieldOrder {
		if _, ok := envelope[key]; ok {
			delete(envelope, key)
			break
		}
	}
	return envelope, nil
}

// rowsFieldOrder lists every rows field the manifest can name, most common first.
var rowsFieldOrder = []string{
	"rows", "items", "securities", "events", "tags", "sections", "groups",
}

func mapKeys(envelope map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(envelope))
	for key := range envelope {
		keys = append(keys, key)
	}
	return keys
}

func first(data []byte, limit int) string {
	if len(data) <= limit {
		return string(data)
	}
	return string(data[:limit])
}

// KLine is one historical candle. Time is YYYYMMDD for daily-or-longer periods
// and a component time index for intraday periods.
type KLine struct {
	Time          int64   `json:"time"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	Volume        int64   `json:"volume"`
	Turnover      int64   `json:"turnover"`
	PreviousClose float64 `json:"previous_close,omitempty"`
}

// TradeDate interprets a daily-or-longer candle's Time as a calendar date.
func (k KLine) TradeDate() (time.Time, error) {
	return time.ParseInLocation("20060102", fmt.Sprintf("%08d", k.Time), Shanghai)
}

// IntradayBar is one minute of a trading session.
type IntradayBar struct {
	Time          int64   `json:"time"`
	Price         float64 `json:"price"`
	Volume        int64   `json:"volume"`
	Turnover      int64   `json:"turnover"`
	LeadIndicator int64   `json:"lead_indicator"`
}

// Tick is one trade print.
type Tick struct {
	Code             string  `json:"code"`
	Time             int64   `json:"time"`
	Price            float64 `json:"price"`
	Volume           int64   `json:"volume"`
	DealType         int     `json:"deal_type"`
	TransactionCount int64   `json:"transaction_count"`
}

// Quote is one price level of the order book.
type Quote struct {
	Price  float64 `json:"price"`
	Volume int64   `json:"volume"`
}

// OrderBook is one instrument's bid/ask ladder.
type OrderBook struct {
	Security Security `json:"security"`
	Depth    string   `json:"depth"`
	Bids     []Quote  `json:"bids"`
	Asks     []Quote  `json:"asks"`
}

// SecurityRef is the identity plus display fields returned by catalog methods.
type SecurityRef struct {
	FullCode      string   `json:"full_code"`
	Code          string   `json:"code"`
	Market        string   `json:"market"`
	Name          string   `json:"name"`
	CodeDisplay   string   `json:"code_display"`
	MarketDisplay string   `json:"market_display"`
	Security      Security `json:"security"`
}

func requireText(argument, value string) error {
	if len(value) == 0 {
		return invalidArgument(argument, "必须是非空字符串")
	}
	return nil
}
