# API 参考

80 个标准方法与 12 个便捷方法。字段名就是原生组件请求体里的 snake_case 键；Go 侧导出名用大驼峰，
`ID`、`AH`、`NEEQ` 等缩写保持大写。本文件与 `methods_gen.go` 同源，均来自上游方法清单。

- [怎么读这张表](#conventions)
- [行情、交易序列与基本面](#cat-market_data)（17 个）
- [市场目录、证券发现与关系](#cat-securities)（18 个）
- [资讯、事件与异动](#cat-news)（14 个）
- [实时统计、排行与排序](#cat-rankings)（11 个）
- [问财](#cat-wencai)（7 个）
- [账户与自选股](#cat-account)（13 个）
- [便捷方法](#convenience)

<a id="conventions"></a>

## 怎么读这张表

```go
data, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Start: -20, End: 0, Adjust: thsdk.AdjustForward, Period: thsdk.PeriodDay,
})

var rows []thsdk.KLine
err = data.Rows(&rows)       // 顶层数组直接解；对象则抬出 rows/items/securities/…
meta, err := data.Metadata() // 外壳里除行集以外的键（分页游标、总数等）
err = data.Into(&single)     // 单对象结果
```

| 记号 | 含义 |
| --- | --- |
| `Opt` | 可选字段。不赋值即不发送；要发送零值用 `thsdk.Num(0)`、`Text("")`、`Flag(false)` |
| `Security` | `code` + 四位字母 `market`；用 `MustSecurity("USHA600519")` 或 `NewSecurity` 构造 |
| `Block` | `{ID: &n}` 与 `{Name: "沪深300"}` 二选一 |
| `any` | 自由结构（多为对象），按上游请求体形状填写 |
| 写操作 | 会改动远端账号数据，默认不返回行集 |

必填字段有静态检查：缺失时返回 `*ParameterError`，**不会**发出请求，也不消耗限流额度。
返回 `数组` 时用 `data.Rows(&x)`；返回 `对象` 时同样调用 `Rows`，SDK 会自动抬出行集键。


<a id="cat-market_data"></a>
## 行情、交易序列与基本面

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`CheckMarketTimelineReadiness`](#checkmarkettimelinereadiness) | `check_market_timeline_readiness` | `CheckMarketTimelineReadinessRequest` | `market`、`time_index_day` | — | 对象 `MarketTimelineReadinessResult` |
| [`GetSecurityTradingTimeline`](#getsecuritytradingtimeline) | `get_security_trading_timeline` | `GetSecurityTradingTimelineRequest` | `security` | — | 对象 `SecurityTradingTimeline`，行集在 `items` |
| [`ListIndexCallAuctionQuotes`](#listindexcallauctionquotes) | `list_index_call_auction_quotes` | `ListIndexCallAuctionQuotesRequest` | `index`、`phase` | — | 数组 `IndexCallAuctionQuote` |
| [`ListSecurityCallAuctionQuotes`](#listsecuritycallauctionquotes) | `list_security_call_auction_quotes` | `ListSecurityCallAuctionQuotesRequest` | `security`、`phase` | `date`、`window` | 数组 `CallAuctionQuote` |
| [`ListSecurityCorporateActions`](#listsecuritycorporateactions) | `list_security_corporate_actions` | `ListSecurityCorporateActionsRequest` | `security`、`start`、`end`、`adjust`、`period` | — | 数组 `CorporateAction` |
| [`ListSecurityDailyCapitalFlows`](#listsecuritydailycapitalflows) | `list_security_daily_capital_flows` | `ListSecurityDailyCapitalFlowsRequest` | `securities` | — | 数组 `DailyCapitalFlow` |
| [`ListSecurityDailyKLinesWithPreviousClose`](#listsecuritydailyklineswithpreviousclose) | `list_security_daily_k_lines_with_previous_close` | `ListSecurityDailyKLinesWithPreviousCloseRequest` | `security`、`start`、`end`、`adjust` | — | 数组 `DailyKLineWithPreviousClose` |
| [`ListSecurityExtendedHoursIntradayBars`](#listsecurityextendedhoursintradaybars) | `list_security_extended_hours_intraday_bars` | `ListSecurityExtendedHoursIntradayBarsRequest` | `security`、`session`、`date` | — | 数组 `IntradayBar` |
| [`ListSecurityExtendedHoursTicks`](#listsecurityextendedhoursticks) | `list_security_extended_hours_ticks` | `ListSecurityExtendedHoursTicksRequest` | `security`、`session`、`window` | — | 数组 `Tick` |
| [`ListSecurityFinancialSnapshots`](#listsecurityfinancialsnapshots) | `list_security_financial_snapshots` | `ListSecurityFinancialSnapshotsRequest` | `securities`、`fields` | — | 数组 `SecurityFinancialSnapshot` |
| [`ListSecurityIntradayBars`](#listsecurityintradaybars) | `list_security_intraday_bars` | `ListSecurityIntradayBarsRequest` | `security`、`date` | — | 数组 `IntradayBar` |
| [`ListSecurityKLines`](#listsecurityklines) | `list_security_k_lines` | `ListSecurityKLinesRequest` | `security`、`start`、`end`、`adjust`、`period` | — | 数组 `KLine` |
| [`ListSecurityOptionGreeks`](#listsecurityoptiongreeks) | `list_security_option_greeks` | `ListSecurityOptionGreeksRequest` | `securities` | `metrics` | 数组 `OptionGreeksRow` |
| [`ListSecurityOrderBooks`](#listsecurityorderbooks) | `list_security_order_books` | `ListSecurityOrderBooksRequest` | `securities`、`depth` | — | 数组 `SecurityOrderBook` |
| [`ListSecurityPriceVolumeLevels`](#listsecuritypricevolumelevels) | `list_security_price_volume_levels` | `ListSecurityPriceVolumeLevelsRequest` | `security` | — | 数组 `PriceVolumeLevel` |
| [`ListSecurityTicks`](#listsecurityticks) | `list_security_ticks` | `ListSecurityTicksRequest` | `security`、`window` | — | 数组 `Tick` |
| [`ListSecurityTimeAndSalesTicks`](#listsecuritytimeandsalesticks) | `list_security_time_and_sales_ticks` | `ListSecurityTimeAndSalesTicksRequest` | `security`、`count` | — | 数组 `Tick` |

### 字段明细


<a id="checkmarkettimelinereadiness"></a>
#### `CheckMarketTimelineReadiness`

原生方法 `check_market_timeline_readiness`；返回 对象 `MarketTimelineReadinessResult`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `time_index_day` | `TimeIndexDay` | `string` | **是** |

```go
client.CheckMarketTimelineReadiness(ctx, thsdk.CheckMarketTimelineReadinessRequest{
    Market: "USHA",
    TimeIndexDay: "0",
})
```

<a id="getsecuritytradingtimeline"></a>
#### `GetSecurityTradingTimeline`

原生方法 `get_security_trading_timeline`；返回 对象 `SecurityTradingTimeline`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.GetSecurityTradingTimeline(ctx, thsdk.GetSecurityTradingTimelineRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listindexcallauctionquotes"></a>
#### `ListIndexCallAuctionQuotes`

原生方法 `list_index_call_auction_quotes`；返回 数组 `IndexCallAuctionQuote`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `index` | `Index` | `any` | **是** |
| `phase` | `Phase` | `string` | **是** |

```go
client.ListIndexCallAuctionQuotes(ctx, thsdk.ListIndexCallAuctionQuotesRequest{
    Index: /* 见上游请求体 */,
    Phase: "opening",
})
```

<a id="listsecuritycallauctionquotes"></a>
#### `ListSecurityCallAuctionQuotes`

原生方法 `list_security_call_auction_quotes`；返回 数组 `CallAuctionQuote`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `phase` | `Phase` | `string` | **是** |
| `date` | `Date` | `Opt` | 否 |
| `window` | `Window` | `any` | 否 |

```go
client.ListSecurityCallAuctionQuotes(ctx, thsdk.ListSecurityCallAuctionQuotesRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Phase: "opening",
})
```

<a id="listsecuritycorporateactions"></a>
#### `ListSecurityCorporateActions`

原生方法 `list_security_corporate_actions`；返回 数组 `CorporateAction`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `start` | `Start` | `int` | **是** |
| `end` | `End` | `int` | **是** |
| `adjust` | `Adjust` | `Adjust` | **是** |
| `period` | `Period` | `Period` | **是** |

```go
client.ListSecurityCorporateActions(ctx, thsdk.ListSecurityCorporateActionsRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Start: -20,
    End: 0,
    Adjust: thsdk.AdjustForward,
    Period: thsdk.PeriodDay,
})
```

<a id="listsecuritydailycapitalflows"></a>
#### `ListSecurityDailyCapitalFlows`

原生方法 `list_security_daily_capital_flows`；返回 数组 `DailyCapitalFlow`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.ListSecurityDailyCapitalFlows(ctx, thsdk.ListSecurityDailyCapitalFlowsRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="listsecuritydailyklineswithpreviousclose"></a>
#### `ListSecurityDailyKLinesWithPreviousClose`

原生方法 `list_security_daily_k_lines_with_previous_close`；返回 数组 `DailyKLineWithPreviousClose`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `start` | `Start` | `int` | **是** |
| `end` | `End` | `int` | **是** |
| `adjust` | `Adjust` | `Adjust` | **是** |

```go
client.ListSecurityDailyKLinesWithPreviousClose(ctx, thsdk.ListSecurityDailyKLinesWithPreviousCloseRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Start: -20,
    End: 0,
    Adjust: thsdk.AdjustForward,
})
```

<a id="listsecurityextendedhoursintradaybars"></a>
#### `ListSecurityExtendedHoursIntradayBars`

原生方法 `list_security_extended_hours_intraday_bars`；返回 数组 `IntradayBar`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `session` | `Session` | `string` | **是** |
| `date` | `Date` | `string` | **是** |

```go
client.ListSecurityExtendedHoursIntradayBars(ctx, thsdk.ListSecurityExtendedHoursIntradayBarsRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Session: "pre_market",
    Date: "0",
})
```

<a id="listsecurityextendedhoursticks"></a>
#### `ListSecurityExtendedHoursTicks`

原生方法 `list_security_extended_hours_ticks`；返回 数组 `Tick`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `session` | `Session` | `string` | **是** |
| `window` | `Window` | `any` | **是** |

```go
client.ListSecurityExtendedHoursTicks(ctx, thsdk.ListSecurityExtendedHoursTicksRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Session: "pre_market",
    Window: /* 见上游请求体 */,
})
```

<a id="listsecurityfinancialsnapshots"></a>
#### `ListSecurityFinancialSnapshots`

原生方法 `list_security_financial_snapshots`；返回 数组 `SecurityFinancialSnapshot`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `fields` | `Fields` | `any` | **是** |

```go
client.ListSecurityFinancialSnapshots(ctx, thsdk.ListSecurityFinancialSnapshotsRequest{
    Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
    Fields: /* 见上游请求体 */,
})
```

<a id="listsecurityintradaybars"></a>
#### `ListSecurityIntradayBars`

原生方法 `list_security_intraday_bars`；返回 数组 `IntradayBar`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `date` | `Date` | `string` | **是** |

```go
client.ListSecurityIntradayBars(ctx, thsdk.ListSecurityIntradayBarsRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Date: "0",
})
```

<a id="listsecurityklines"></a>
#### `ListSecurityKLines`

原生方法 `list_security_k_lines`；返回 数组 `KLine`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `start` | `Start` | `int` | **是** |
| `end` | `End` | `int` | **是** |
| `adjust` | `Adjust` | `Adjust` | **是** |
| `period` | `Period` | `Period` | **是** |

```go
client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Start: -20,
    End: 0,
    Adjust: thsdk.AdjustForward,
    Period: thsdk.PeriodDay,
})
```

<a id="listsecurityoptiongreeks"></a>
#### `ListSecurityOptionGreeks`

原生方法 `list_security_option_greeks`；返回 数组 `OptionGreeksRow`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `metrics` | `Metrics` | `any` | 否 |

```go
client.ListSecurityOptionGreeks(ctx, thsdk.ListSecurityOptionGreeksRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="listsecurityorderbooks"></a>
#### `ListSecurityOrderBooks`

原生方法 `list_security_order_books`；返回 数组 `SecurityOrderBook`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `depth` | `Depth` | `OrderBookDepth` | **是** |

```go
client.ListSecurityOrderBooks(ctx, thsdk.ListSecurityOrderBooksRequest{
    Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
    Depth: thsdk.DepthFive,
})
```

<a id="listsecuritypricevolumelevels"></a>
#### `ListSecurityPriceVolumeLevels`

原生方法 `list_security_price_volume_levels`；返回 数组 `PriceVolumeLevel`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.ListSecurityPriceVolumeLevels(ctx, thsdk.ListSecurityPriceVolumeLevelsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listsecurityticks"></a>
#### `ListSecurityTicks`

原生方法 `list_security_ticks`；返回 数组 `Tick`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `window` | `Window` | `any` | **是** |

```go
client.ListSecurityTicks(ctx, thsdk.ListSecurityTicksRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Window: /* 见上游请求体 */,
})
```

<a id="listsecuritytimeandsalesticks"></a>
#### `ListSecurityTimeAndSalesTicks`

原生方法 `list_security_time_and_sales_ticks`；返回 数组 `Tick`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `count` | `Count` | `int` | **是** |

```go
client.ListSecurityTimeAndSalesTicks(ctx, thsdk.ListSecurityTimeAndSalesTicksRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Count: 50,
})
```

<a id="cat-securities"></a>
## 市场目录、证券发现与关系

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`GetMarketMetadata`](#getmarketmetadata) | `get_market_metadata` | `GetMarketMetadataRequest` | `market` | `version_flag` | 对象 `MarketMetadata`，行集在 `sections` |
| [`GetMarketSecurityNames`](#getmarketsecuritynames) | `get_market_security_names` | `GetMarketSecurityNamesRequest` | `market` | `realtime_version_ini`、`history_version_ini`、`base_version_ini` | 对象 `MarketSecurityNamesResult`，行集在 `items` |
| [`GetSecurityConceptTags`](#getsecurityconcepttags) | `get_security_concept_tags` | `GetSecurityConceptTagsRequest` | `security` | — | 对象 `SecurityConceptTagsResult`，行集在 `tags` |
| [`GetSecurityIndustry`](#getsecurityindustry) | `get_security_industry` | `GetSecurityIndustryRequest` | `security` | — | 对象 `SecurityIndustryMapping` |
| [`ListBlockConstituents`](#listblockconstituents) | `list_block_constituents` | `ListBlockConstituentsRequest` | `block` | `sort_begin`、`sort_count`、`sort_order`、`sort_id` | 数组 `BlockConstituent` |
| [`ListBlockDescriptions`](#listblockdescriptions) | `list_block_descriptions` | `ListBlockDescriptionsRequest` | `blocks` | — | 数组 `BlockDescription` |
| [`ListFuturesRelatedSecurities`](#listfuturesrelatedsecurities) | `list_futures_related_securities` | 无参数 | — | — | 数组 `FuturesRelatedSecurity` |
| [`ListIndustryChildren`](#listindustrychildren) | `list_industry_children` | `ListIndustryChildrenRequest` | `industry` | — | 数组 `IndustryChild` |
| [`ListMarketSecurities`](#listmarketsecurities) | `list_market_securities` | `ListMarketSecuritiesRequest` | `market` | `sort_begin`、`sort_count`、`sort_id`、`sort_order`、`func_period` | 对象 `MarketSecuritiesResult`，行集在 `securities` |
| [`ListRelatedSecurityPerformances`](#listrelatedsecurityperformances) | `list_related_security_performances` | `ListRelatedSecurityPerformancesRequest` | `security` | — | 数组 `RelatedSecurityPerformance` |
| [`ListSecurityAHRelations`](#listsecurityahrelations) | `list_security_ah_relations` | `ListSecurityAHRelationsRequest` | `security` | — | 数组 `SecurityAHRelation` |
| [`ListSecurityBlockMemberships`](#listsecurityblockmemberships) | `list_security_block_memberships` | `ListSecurityBlockMembershipsRequest` | `security` | — | 数组 `SecurityBlockMembership` |
| [`ListSecurityFuturesRelations`](#listsecurityfuturesrelations) | `list_security_futures_relations` | `[]Security` | `securities` | — | 数组 `SecurityFuturesRelation` |
| [`ListSecurityIndustryMappings`](#listsecurityindustrymappings) | `list_security_industry_mappings` | `ListSecurityIndustryMappingsRequest` | `securities` | — | 数组 `SecurityIndustryMapping` |
| [`ListSecurityLinkRelations`](#listsecuritylinkrelations) | `list_security_link_relations` | `string` | `link_key` | — | 数组 `SecurityLinkRelation` |
| [`ResolveBlock`](#resolveblock) | `resolve_block` | `ResolveBlockRequest` | `name` | — | 对象 `Block` |
| [`ResolveSecurities`](#resolvesecurities) | `resolve_securities` | `ResolveSecuritiesRequest` | `codes` | — | 数组 `ResolvedSecurity` |
| [`SearchSecurities`](#searchsecurities) | `search_securities` | `SearchSecuritiesRequest` | `pattern` | `market` | 数组 `SecurityCandidate` |

### 字段明细


<a id="getmarketmetadata"></a>
#### `GetMarketMetadata`

原生方法 `get_market_metadata`；返回 对象 `MarketMetadata`，行集在 `sections`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `version_flag` | `VersionFlag` | `Opt` | 否 |

```go
client.GetMarketMetadata(ctx, thsdk.GetMarketMetadataRequest{Market: "USHA"})
```

<a id="getmarketsecuritynames"></a>
#### `GetMarketSecurityNames`

原生方法 `get_market_security_names`；返回 对象 `MarketSecurityNamesResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `realtime_version_ini` | `RealtimeVersionIni` | `Opt` | 否 |
| `history_version_ini` | `HistoryVersionIni` | `Opt` | 否 |
| `base_version_ini` | `BaseVersionIni` | `Opt` | 否 |

```go
client.GetMarketSecurityNames(ctx, thsdk.GetMarketSecurityNamesRequest{Market: "USHA"})
```

<a id="getsecurityconcepttags"></a>
#### `GetSecurityConceptTags`

原生方法 `get_security_concept_tags`；返回 对象 `SecurityConceptTagsResult`，行集在 `tags`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.GetSecurityConceptTags(ctx, thsdk.GetSecurityConceptTagsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="getsecurityindustry"></a>
#### `GetSecurityIndustry`

原生方法 `get_security_industry`；返回 对象 `SecurityIndustryMapping`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.GetSecurityIndustry(ctx, thsdk.GetSecurityIndustryRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listblockconstituents"></a>
#### `ListBlockConstituents`

原生方法 `list_block_constituents`；返回 数组 `BlockConstituent`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `block` | `Block` | `Block` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |

```go
client.ListBlockConstituents(ctx, thsdk.ListBlockConstituentsRequest{Block: thsdk.Block{Name: "沪深300"}})
```

<a id="listblockdescriptions"></a>
#### `ListBlockDescriptions`

原生方法 `list_block_descriptions`；返回 数组 `BlockDescription`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `blocks` | `Blocks` | `[]Block` | **是** |

```go
client.ListBlockDescriptions(ctx, thsdk.ListBlockDescriptionsRequest{Blocks: []thsdk.Block{{Name: "沪深300"}}})
```

<a id="listfuturesrelatedsecurities"></a>
#### `ListFuturesRelatedSecurities`

原生方法 `list_futures_related_securities`；返回 数组 `FuturesRelatedSecurity`。

```go
client.ListFuturesRelatedSecurities(ctx)
```

<a id="listindustrychildren"></a>
#### `ListIndustryChildren`

原生方法 `list_industry_children`；返回 数组 `IndustryChild`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `industry` | `Industry` | `any` | **是** |

```go
client.ListIndustryChildren(ctx, thsdk.ListIndustryChildrenRequest{Industry: /* 见上游请求体 */})
```

<a id="listmarketsecurities"></a>
#### `ListMarketSecurities`

原生方法 `list_market_securities`；返回 对象 `MarketSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `func_period` | `FuncPeriod` | `Opt` | 否 |

```go
client.ListMarketSecurities(ctx, thsdk.ListMarketSecuritiesRequest{Market: "USHA"})
```

<a id="listrelatedsecurityperformances"></a>
#### `ListRelatedSecurityPerformances`

原生方法 `list_related_security_performances`；返回 数组 `RelatedSecurityPerformance`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.ListRelatedSecurityPerformances(ctx, thsdk.ListRelatedSecurityPerformancesRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listsecurityahrelations"></a>
#### `ListSecurityAHRelations`

原生方法 `list_security_ah_relations`；返回 数组 `SecurityAHRelation`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.ListSecurityAHRelations(ctx, thsdk.ListSecurityAHRelationsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listsecurityblockmemberships"></a>
#### `ListSecurityBlockMemberships`

原生方法 `list_security_block_memberships`；返回 数组 `SecurityBlockMembership`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.ListSecurityBlockMemberships(ctx, thsdk.ListSecurityBlockMembershipsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listsecurityfuturesrelations"></a>
#### `ListSecurityFuturesRelations`

原生方法 `list_security_futures_relations`；返回 数组 `SecurityFuturesRelation`。

```go
client.ListSecurityFuturesRelations(ctx, []thsdk.Security{thsdk.MustSecurity("USHA600519")})
```

<a id="listsecurityindustrymappings"></a>
#### `ListSecurityIndustryMappings`

原生方法 `list_security_industry_mappings`；返回 数组 `SecurityIndustryMapping`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.ListSecurityIndustryMappings(ctx, thsdk.ListSecurityIndustryMappingsRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="listsecuritylinkrelations"></a>
#### `ListSecurityLinkRelations`

原生方法 `list_security_link_relations`；返回 数组 `SecurityLinkRelation`。

```go
client.ListSecurityLinkRelations(ctx, "USHA600519")
```

<a id="resolveblock"></a>
#### `ResolveBlock`

原生方法 `resolve_block`；返回 对象 `Block`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `name` | `Name` | `string` | **是** |

```go
client.ResolveBlock(ctx, thsdk.ResolveBlockRequest{Name: "沪深300"})
```

<a id="resolvesecurities"></a>
#### `ResolveSecurities`

原生方法 `resolve_securities`；返回 数组 `ResolvedSecurity`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `codes` | `Codes` | `[]string` | **是** |

```go
client.ResolveSecurities(ctx, thsdk.ResolveSecuritiesRequest{Codes: []string{"600519"}})
```

<a id="searchsecurities"></a>
#### `SearchSecurities`

原生方法 `search_securities`；返回 数组 `SecurityCandidate`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `pattern` | `Pattern` | `string` | **是** |
| `market` | `Market` | `Opt` | 否 |

```go
client.SearchSecurities(ctx, thsdk.SearchSecuritiesRequest{Pattern: "茅台"})
```

<a id="cat-news"></a>
## 资讯、事件与异动

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`AnalyzeSecurityLimitUp`](#analyzesecuritylimitup) | `analyze_security_limit_up` | `AnalyzeSecurityLimitUpRequest` | `security` | `date` | 对象 `SecurityLimitUpAnalysisResult` |
| [`GetSecurityPriceLimitEvents`](#getsecuritypricelimitevents) | `get_security_price_limit_events` | `GetSecurityPriceLimitEventsRequest` | `security` | `start_year`、`end_year` | 对象 `PriceLimitEventsResult`，行集在 `events` |
| [`GetSecurityShortTermHighlights`](#getsecurityshorttermhighlights) | `get_security_short_term_highlights` | `GetSecurityShortTermHighlightsRequest` | `security` | — | 对象 `SecurityShortTermHighlightsResult` |
| [`ListCommodityStockLinkageNews`](#listcommoditystocklinkagenews) | `list_commodity_stock_linkage_news` | `ListCommodityStockLinkageNewsRequest` | `commodity_index` | — | 对象 `CommodityStockLinkageNewsResult`，行集在 `items` |
| [`ListHotBlockCalendarEvents`](#listhotblockcalendarevents) | `list_hot_block_calendar_events` | 无参数 | — | — | 对象 `HotBlockCalendarResult`，行集在 `items` |
| [`ListHotEventNews`](#listhoteventnews) | `list_hot_event_news` | `ListHotEventNewsRequest` | `block` | `range` | 对象 `HotEventNewsResult`，行集在 `items` |
| [`ListMarketCallAuctionUnusualEvents`](#listmarketcallauctionunusualevents) | `list_market_call_auction_unusual_events` | `ListMarketCallAuctionUnusualEventsRequest` | `market` | `trade_date` | 数组 `ShortTermEvent` |
| [`ListMarketChangeEvents`](#listmarketchangeevents) | `list_market_change_events` | 无参数 | — | — | 数组 `MarketChangeEvent` |
| [`ListMarketShortTermEvents`](#listmarketshorttermevents) | `list_market_short_term_events` | `ListMarketShortTermEventsRequest` | `market`、`mode` | `securities`、`max_count`、`before_timestamp_micros` | 数组 `ShortTermEvent` |
| [`ListNewsFlashItems`](#listnewsflashitems) | `list_news_flash_items` | `ListNewsFlashItemsRequest` | `traversal` | `tag`、`tag_id`、`page`、`page_size`、`page_time`、`sequence`、`created_at`、`environment` | 对象 `NewsFlashItemsResult`，行集在 `items` |
| [`ListNewsItems`](#listnewsitems) | `list_news_items` | `ListNewsItemsRequest` | `category`、`traversal` | `security`、`market`、`summary_mode`、`advertorial` | 对象 `NewsItemsResult`，行集在 `items` |
| [`ListSecurityNewsEvents`](#listsecuritynewsevents) | `list_security_news_events` | `ListSecurityNewsEventsRequest` | `security`、`timeline`、`window` | — | 对象 `SecurityNewsItemsResult`，行集在 `items` |
| [`ListSecurityNewsMarkers`](#listsecuritynewsmarkers) | `list_security_news_markers` | `ListSecurityNewsMarkersRequest` | `security`、`timeline` | — | 对象 `SecurityNewsItemsResult`，行集在 `items` |
| [`ListSecurityResearchReports`](#listsecurityresearchreports) | `list_security_research_reports` | `ListSecurityResearchReportsRequest` | `security`、`traversal` | — | 对象 `SecurityResearchReportsResult`，行集在 `items` |

### 字段明细


<a id="analyzesecuritylimitup"></a>
#### `AnalyzeSecurityLimitUp`

原生方法 `analyze_security_limit_up`；返回 对象 `SecurityLimitUpAnalysisResult`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `date` | `Date` | `Opt` | 否 |

```go
client.AnalyzeSecurityLimitUp(ctx, thsdk.AnalyzeSecurityLimitUpRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="getsecuritypricelimitevents"></a>
#### `GetSecurityPriceLimitEvents`

原生方法 `get_security_price_limit_events`；返回 对象 `PriceLimitEventsResult`，行集在 `events`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `start_year` | `StartYear` | `Opt` | 否 |
| `end_year` | `EndYear` | `Opt` | 否 |

```go
client.GetSecurityPriceLimitEvents(ctx, thsdk.GetSecurityPriceLimitEventsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="getsecurityshorttermhighlights"></a>
#### `GetSecurityShortTermHighlights`

原生方法 `get_security_short_term_highlights`；返回 对象 `SecurityShortTermHighlightsResult`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |

```go
client.GetSecurityShortTermHighlights(ctx, thsdk.GetSecurityShortTermHighlightsRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listcommoditystocklinkagenews"></a>
#### `ListCommodityStockLinkageNews`

原生方法 `list_commodity_stock_linkage_news`；返回 对象 `CommodityStockLinkageNewsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `commodity_index` | `CommodityIndex` | `any` | **是** |

```go
client.ListCommodityStockLinkageNews(ctx, thsdk.ListCommodityStockLinkageNewsRequest{CommodityIndex: /* 见上游请求体 */})
```

<a id="listhotblockcalendarevents"></a>
#### `ListHotBlockCalendarEvents`

原生方法 `list_hot_block_calendar_events`；返回 对象 `HotBlockCalendarResult`，行集在 `items`。

```go
client.ListHotBlockCalendarEvents(ctx)
```

<a id="listhoteventnews"></a>
#### `ListHotEventNews`

原生方法 `list_hot_event_news`；返回 对象 `HotEventNewsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `block` | `Block` | `Block` | **是** |
| `range` | `Range` | `Opt` | 否 |

```go
client.ListHotEventNews(ctx, thsdk.ListHotEventNewsRequest{Block: thsdk.Block{Name: "沪深300"}})
```

<a id="listmarketcallauctionunusualevents"></a>
#### `ListMarketCallAuctionUnusualEvents`

原生方法 `list_market_call_auction_unusual_events`；返回 数组 `ShortTermEvent`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `trade_date` | `TradeDate` | `Opt` | 否 |

```go
client.ListMarketCallAuctionUnusualEvents(ctx, thsdk.ListMarketCallAuctionUnusualEventsRequest{Market: "USHA"})
```

<a id="listmarketchangeevents"></a>
#### `ListMarketChangeEvents`

原生方法 `list_market_change_events`；返回 数组 `MarketChangeEvent`。

```go
client.ListMarketChangeEvents(ctx)
```

<a id="listmarketshorttermevents"></a>
#### `ListMarketShortTermEvents`

原生方法 `list_market_short_term_events`；返回 数组 `ShortTermEvent`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `market` | `Market` | `string` | **是** |
| `mode` | `Mode` | `string` | **是** |
| `securities` | `Securities` | `[]Security` | 否 |
| `max_count` | `MaxCount` | `Opt` | 否 |
| `before_timestamp_micros` | `BeforeTimestampMicros` | `Opt` | 否 |

```go
client.ListMarketShortTermEvents(ctx, thsdk.ListMarketShortTermEventsRequest{
    Market: "USHA",
    Mode: "market_recent",
})
```

<a id="listnewsflashitems"></a>
#### `ListNewsFlashItems`

原生方法 `list_news_flash_items`；返回 对象 `NewsFlashItemsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `traversal` | `Traversal` | `any` | **是** |
| `tag` | `Tag` | `Opt` | 否 |
| `tag_id` | `TagID` | `Opt` | 否 |
| `page` | `Page` | `Opt` | 否 |
| `page_size` | `PageSize` | `Opt` | 否 |
| `page_time` | `PageTime` | `Opt` | 否 |
| `sequence` | `Sequence` | `Opt` | 否 |
| `created_at` | `CreatedAt` | `Opt` | 否 |
| `environment` | `Environment` | `Opt` | 否 |

```go
client.ListNewsFlashItems(ctx, thsdk.ListNewsFlashItemsRequest{Traversal: "page"})
```

<a id="listnewsitems"></a>
#### `ListNewsItems`

原生方法 `list_news_items`；返回 对象 `NewsItemsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `category` | `Category` | `int` | **是** |
| `security` | `Security` | `Security` | 否 |
| `market` | `Market` | `Opt` | 否 |
| `traversal` | `Traversal` | `any` | **是** |
| `summary_mode` | `SummaryMode` | `Opt` | 否 |
| `advertorial` | `Advertorial` | `Opt` | 否 |

```go
client.ListNewsItems(ctx, thsdk.ListNewsItemsRequest{
    Category: thsdk.InfoStockInfo,
    Traversal: "page",
})
```

<a id="listsecuritynewsevents"></a>
#### `ListSecurityNewsEvents`

原生方法 `list_security_news_events`；返回 对象 `SecurityNewsItemsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `timeline` | `Timeline` | `string` | **是** |
| `window` | `Window` | `any` | **是** |

```go
client.ListSecurityNewsEvents(ctx, thsdk.ListSecurityNewsEventsRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Timeline: "historical",
    Window: /* 见上游请求体 */,
})
```

<a id="listsecuritynewsmarkers"></a>
#### `ListSecurityNewsMarkers`

原生方法 `list_security_news_markers`；返回 对象 `SecurityNewsItemsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `timeline` | `Timeline` | `string` | **是** |

```go
client.ListSecurityNewsMarkers(ctx, thsdk.ListSecurityNewsMarkersRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Timeline: "historical",
})
```

<a id="listsecurityresearchreports"></a>
#### `ListSecurityResearchReports`

原生方法 `list_security_research_reports`；返回 对象 `SecurityResearchReportsResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `traversal` | `Traversal` | `any` | **是** |

```go
client.ListSecurityResearchReports(ctx, thsdk.ListSecurityResearchReportsRequest{
    Security: thsdk.MustSecurity("USHA600519"),
    Traversal: "page",
})
```

<a id="cat-rankings"></a>
## 实时统计、排行与排序

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`CalculateSecurityRealtimeStatistics`](#calculatesecurityrealtimestatistics) | `calculate_security_realtime_statistics` | `CalculateSecurityRealtimeStatisticsRequest` | `securities`、`metrics` | — | 对象 `RealtimeStatsResult`，行集在 `rows` |
| [`GetSecurityPopularityRank`](#getsecuritypopularityrank) | `get_security_popularity_rank` | `GetSecurityPopularityRankRequest` | `security` | `benchmark` | 对象 `PopularityRankResult` |
| [`ListSecuritiesByRealtimeSignal`](#listsecuritiesbyrealtimesignal) | `list_securities_by_realtime_signal` | `ListSecuritiesByRealtimeSignalRequest` | `markets`、`signal` | — | 数组 `RealtimeSignalSecurity` |
| [`RankBlockSecurities`](#rankblocksecurities) | `rank_block_securities` | `RankBlockSecuritiesRequest` | `block` | `sort_begin`、`sort_count`、`sort_id`、`sort_order`、`exclude_securities`、`exclude_blocks`、`func_period` | 对象 `SortedSecuritiesResult`，行集在 `securities` |
| [`RankBlockSecuritiesByIndustry`](#rankblocksecuritiesbyindustry) | `rank_block_securities_by_industry` | `RankBlockSecuritiesByIndustryRequest` | `block` | `sort_begin`、`sort_count`、`sort_order` | 对象 `SortedSecuritiesResult`，行集在 `securities` |
| [`RankConstituentsByPerformanceContribution`](#rankconstituentsbyperformancecontribution) | `rank_constituents_by_performance_contribution` | `RankConstituentsByPerformanceContributionRequest` | `code`、`market` | `target`、`sort_id`、`sort_order`、`sort_begin`、`sort_count`、`valid_begin`、`valid_end` | 数组 `PerformanceContributionItem` |
| [`RankRelatedSecurities`](#rankrelatedsecurities) | `rank_related_securities` | `RankRelatedSecuritiesRequest` | `security` | `sort_begin`、`sort_count`、`sort_id`、`sort_order`、`func_period` | 对象 `SortedSecuritiesResult`，行集在 `securities` |
| [`RankSecuritiesByPopularity`](#ranksecuritiesbypopularity) | `rank_securities_by_popularity` | `RankSecuritiesByPopularityRequest` | — | `type`、`benchmark` | 对象 `PopularityRankTopResult`，行集在 `items` |
| [`RankSecuritiesByRealtimeStatistic`](#ranksecuritiesbyrealtimestatistic) | `rank_securities_by_realtime_statistic` | `RankSecuritiesByRealtimeStatisticRequest` | `sort_by` | `markets`、`securities`、`sort_dir`、`sort_begin`、`sort_count`、`exclude_securities` | 对象 `RealtimeStatsResult`，行集在 `rows` |
| [`SortSecurities`](#sortsecurities) | `sort_securities` | `SortSecuritiesRequest` | `securities` | `sort_begin`、`sort_count`、`sort_id`、`sort_order`、`func_period` | 对象 `SortedSecuritiesResult`，行集在 `securities` |
| [`SortSecuritiesByIndustry`](#sortsecuritiesbyindustry) | `sort_securities_by_industry` | `SortSecuritiesByIndustryRequest` | `securities` | `sort_begin`、`sort_count`、`sort_order` | 对象 `SortedSecuritiesResult`，行集在 `securities` |

### 字段明细


<a id="calculatesecurityrealtimestatistics"></a>
#### `CalculateSecurityRealtimeStatistics`

原生方法 `calculate_security_realtime_statistics`；返回 对象 `RealtimeStatsResult`，行集在 `rows`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `metrics` | `Metrics` | `any` | **是** |

```go
client.CalculateSecurityRealtimeStatistics(ctx, thsdk.CalculateSecurityRealtimeStatisticsRequest{
    Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
    Metrics: /* 见上游请求体 */,
})
```

<a id="getsecuritypopularityrank"></a>
#### `GetSecurityPopularityRank`

原生方法 `get_security_popularity_rank`；返回 对象 `PopularityRankResult`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `benchmark` | `Benchmark` | `any` | 否 |

```go
client.GetSecurityPopularityRank(ctx, thsdk.GetSecurityPopularityRankRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="listsecuritiesbyrealtimesignal"></a>
#### `ListSecuritiesByRealtimeSignal`

原生方法 `list_securities_by_realtime_signal`；返回 数组 `RealtimeSignalSecurity`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `markets` | `Markets` | `[]string` | **是** |
| `signal` | `Signal` | `string` | **是** |

```go
client.ListSecuritiesByRealtimeSignal(ctx, thsdk.ListSecuritiesByRealtimeSignalRequest{
    Markets: []string{"USHA", "USZA"},
    Signal: "limit_up",
})
```

<a id="rankblocksecurities"></a>
#### `RankBlockSecurities`

原生方法 `rank_block_securities`；返回 对象 `SortedSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `block` | `Block` | `Block` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `exclude_securities` | `ExcludeSecurities` | `[]Security` | 否 |
| `exclude_blocks` | `ExcludeBlocks` | `[]Block` | 否 |
| `func_period` | `FuncPeriod` | `Opt` | 否 |

```go
client.RankBlockSecurities(ctx, thsdk.RankBlockSecuritiesRequest{Block: thsdk.Block{Name: "沪深300"}})
```

<a id="rankblocksecuritiesbyindustry"></a>
#### `RankBlockSecuritiesByIndustry`

原生方法 `rank_block_securities_by_industry`；返回 对象 `SortedSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `block` | `Block` | `Block` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |

```go
client.RankBlockSecuritiesByIndustry(ctx, thsdk.RankBlockSecuritiesByIndustryRequest{Block: thsdk.Block{Name: "沪深300"}})
```

<a id="rankconstituentsbyperformancecontribution"></a>
#### `RankConstituentsByPerformanceContribution`

原生方法 `rank_constituents_by_performance_contribution`；返回 数组 `PerformanceContributionItem`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `code` | `Code` | `string` | **是** |
| `market` | `Market` | `string` | **是** |
| `target` | `Target` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `valid_begin` | `ValidBegin` | `Opt` | 否 |
| `valid_end` | `ValidEnd` | `Opt` | 否 |

```go
client.RankConstituentsByPerformanceContribution(ctx, thsdk.RankConstituentsByPerformanceContributionRequest{
    Code: "000300",
    Market: "USHA",
})
```

<a id="rankrelatedsecurities"></a>
#### `RankRelatedSecurities`

原生方法 `rank_related_securities`；返回 对象 `SortedSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `security` | `Security` | `Security` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `func_period` | `FuncPeriod` | `Opt` | 否 |

```go
client.RankRelatedSecurities(ctx, thsdk.RankRelatedSecuritiesRequest{Security: thsdk.MustSecurity("USHA600519")})
```

<a id="ranksecuritiesbypopularity"></a>
#### `RankSecuritiesByPopularity`

原生方法 `rank_securities_by_popularity`；返回 对象 `PopularityRankTopResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `type` | `Type` | `Opt` | 否 |
| `benchmark` | `Benchmark` | `any` | 否 |

<a id="ranksecuritiesbyrealtimestatistic"></a>
#### `RankSecuritiesByRealtimeStatistic`

原生方法 `rank_securities_by_realtime_statistic`；返回 对象 `RealtimeStatsResult`，行集在 `rows`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `markets` | `Markets` | `[]string` | 否 |
| `securities` | `Securities` | `[]Security` | 否 |
| `sort_by` | `SortBy` | `int` | **是** |
| `sort_dir` | `SortDir` | `Opt` | 否 |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `exclude_securities` | `ExcludeSecurities` | `[]Security` | 否 |

```go
client.RankSecuritiesByRealtimeStatistic(ctx, thsdk.RankSecuritiesByRealtimeStatisticRequest{SortBy: 1})
```

<a id="sortsecurities"></a>
#### `SortSecurities`

原生方法 `sort_securities`；返回 对象 `SortedSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_id` | `SortID` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |
| `func_period` | `FuncPeriod` | `Opt` | 否 |

```go
client.SortSecurities(ctx, thsdk.SortSecuritiesRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="sortsecuritiesbyindustry"></a>
#### `SortSecuritiesByIndustry`

原生方法 `sort_securities_by_industry`；返回 对象 `SortedSecuritiesResult`，行集在 `securities`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `sort_begin` | `SortBegin` | `Opt` | 否 |
| `sort_count` | `SortCount` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |

```go
client.SortSecuritiesByIndustry(ctx, thsdk.SortSecuritiesByIndustryRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="cat-wencai"></a>
## 问财

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`ListWencaiExpressionFields`](#listwencaiexpressionfields) | `list_wencai_expression_fields` | `string` | `query` | — | 数组 `WencaiField` |
| [`ListWencaiHotBlocks`](#listwencaihotblocks) | `list_wencai_hot_blocks` | `ListWencaiHotBlocksRequest` | `kind` | — | 对象 `WencaiHotBlocksResult`，行集在 `items` |
| [`QueryWencai`](#querywencai) | `query_wencai` | `QueryWencaiRequest` | `query` | `markets`、`limit` | 对象 `WencaiQueryResult`，行集在 `rows` |
| [`QueryWencaiRealtimeFields`](#querywencairealtimefields) | `query_wencai_realtime_fields` | `QueryWencaiRealtimeFieldsRequest` | `securities`、`fields` | — | 对象 `WencaiRealtimeFieldsResult`，行集在 `rows` |
| [`QueryWencaiSecurities`](#querywencaisecurities) | `query_wencai_securities` | `QueryWencaiSecuritiesRequest` | `query` | `markets`、`limit` | 对象 `WencaiSecuritiesResult`，行集在 `items` |
| [`RankBlockSecuritiesByWencaiField`](#rankblocksecuritiesbywencaifield) | `rank_block_securities_by_wencai_field` | `RankBlockSecuritiesByWencaiFieldRequest` | `block`、`field` | `offset`、`limit`、`sort_order` | 对象 `WencaiBlockRankingResult`，行集在 `items` |
| [`ResolveWencaiSecurities`](#resolvewencaisecurities) | `resolve_wencai_securities` | `ResolveWencaiSecuritiesRequest` | `query` | `limit`、`stock_suffix` | 对象 `ResolvedWencaiSecuritiesResult`，行集在 `items` |

### 字段明细


<a id="listwencaiexpressionfields"></a>
#### `ListWencaiExpressionFields`

原生方法 `list_wencai_expression_fields`；返回 数组 `WencaiField`。

```go
client.ListWencaiExpressionFields(ctx, "USHA600519")
```

<a id="listwencaihotblocks"></a>
#### `ListWencaiHotBlocks`

原生方法 `list_wencai_hot_blocks`；返回 对象 `WencaiHotBlocksResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `kind` | `Kind` | `string` | **是** |

```go
client.ListWencaiHotBlocks(ctx, thsdk.ListWencaiHotBlocksRequest{Kind: "concept"})
```

<a id="querywencai"></a>
#### `QueryWencai`

原生方法 `query_wencai`；返回 对象 `WencaiQueryResult`，行集在 `rows`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `query` | `Query` | `string` | **是** |
| `markets` | `Markets` | `[]string` | 否 |
| `limit` | `Limit` | `Opt` | 否 |

```go
client.QueryWencai(ctx, thsdk.QueryWencaiRequest{Query: "市盈率小于20"})
```

<a id="querywencairealtimefields"></a>
#### `QueryWencaiRealtimeFields`

原生方法 `query_wencai_realtime_fields`；返回 对象 `WencaiRealtimeFieldsResult`，行集在 `rows`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `fields` | `Fields` | `any` | **是** |

```go
client.QueryWencaiRealtimeFields(ctx, thsdk.QueryWencaiRealtimeFieldsRequest{
    Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
    Fields: /* 见上游请求体 */,
})
```

<a id="querywencaisecurities"></a>
#### `QueryWencaiSecurities`

原生方法 `query_wencai_securities`；返回 对象 `WencaiSecuritiesResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `query` | `Query` | `string` | **是** |
| `markets` | `Markets` | `[]string` | 否 |
| `limit` | `Limit` | `Opt` | 否 |

```go
client.QueryWencaiSecurities(ctx, thsdk.QueryWencaiSecuritiesRequest{Query: "市盈率小于20"})
```

<a id="rankblocksecuritiesbywencaifield"></a>
#### `RankBlockSecuritiesByWencaiField`

原生方法 `rank_block_securities_by_wencai_field`；返回 对象 `WencaiBlockRankingResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `block` | `Block` | `Block` | **是** |
| `field` | `Field` | `any` | **是** |
| `offset` | `Offset` | `Opt` | 否 |
| `limit` | `Limit` | `Opt` | 否 |
| `sort_order` | `SortOrder` | `Opt` | 否 |

```go
client.RankBlockSecuritiesByWencaiField(ctx, thsdk.RankBlockSecuritiesByWencaiFieldRequest{
    Block: thsdk.Block{Name: "沪深300"},
    Field: /* 见上游请求体 */,
})
```

<a id="resolvewencaisecurities"></a>
#### `ResolveWencaiSecurities`

原生方法 `resolve_wencai_securities`；返回 对象 `ResolvedWencaiSecuritiesResult`，行集在 `items`。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `query` | `Query` | `string` | **是** |
| `limit` | `Limit` | `Opt` | 否 |
| `stock_suffix` | `StockSuffix` | `Opt` | 否 |

```go
client.ResolveWencaiSecurities(ctx, thsdk.ResolveWencaiSecuritiesRequest{Query: "市盈率小于20"})
```

<a id="cat-account"></a>
## 账户与自选股

| Go 方法 | 原生方法 | 请求 | 必填字段 | 可选字段 | 返回 |
| --- | --- | --- | --- | --- | --- |
| [`AccountPermissions`](#accountpermissions) | `account_permissions` | 无参数 | — | — | 对象 `AccountPermissions` |
| [`AddAccountWatchlistGroupSecurities`](#addaccountwatchlistgroupsecurities) | `add_account_watchlist_group_securities` | `AddAccountWatchlistGroupSecuritiesRequest` | `group_id`、`securities` | — | `AccountWatchlistGroups`（写操作） |
| [`AddAccountWatchlistSecurities`](#addaccountwatchlistsecurities) | `add_account_watchlist_securities` | `AddAccountWatchlistSecuritiesRequest` | `securities` | `add_to_front` | `AccountWatchlistUpdate`（写操作） |
| [`ClearAccountWatchlist`](#clearaccountwatchlist) | `clear_account_watchlist` | 无参数 | — | — | `AccountWatchlistUpdate`（写操作） |
| [`CreateAccountWatchlistGroup`](#createaccountwatchlistgroup) | `create_account_watchlist_group` | `CreateAccountWatchlistGroupRequest` | `name` | `securities` | `AccountWatchlistGroups`（写操作） |
| [`DeleteAccountWatchlistGroup`](#deleteaccountwatchlistgroup) | `delete_account_watchlist_group` | `DeleteAccountWatchlistGroupRequest` | `group_id` | — | `AccountWatchlistGroups`（写操作） |
| [`GetAccountWatchlist`](#getaccountwatchlist) | `get_account_watchlist` | 无参数 | — | — | 对象 `AccountWatchlistSnapshot`，行集在 `securities` |
| [`GetAccountWatchlistGroups`](#getaccountwatchlistgroups) | `get_account_watchlist_groups` | 无参数 | — | — | 对象 `AccountWatchlistGroups`，行集在 `groups` |
| [`RemoveAccountWatchlistGroupSecurities`](#removeaccountwatchlistgroupsecurities) | `remove_account_watchlist_group_securities` | `RemoveAccountWatchlistGroupSecuritiesRequest` | `group_id`、`securities` | — | `AccountWatchlistGroups`（写操作） |
| [`RemoveAccountWatchlistSecurities`](#removeaccountwatchlistsecurities) | `remove_account_watchlist_securities` | `RemoveAccountWatchlistSecuritiesRequest` | `securities` | — | `AccountWatchlistUpdate`（写操作） |
| [`RenameAccountWatchlistGroup`](#renameaccountwatchlistgroup) | `rename_account_watchlist_group` | `RenameAccountWatchlistGroupRequest` | `group_id`、`name` | — | `AccountWatchlistGroups`（写操作） |
| [`ReplaceAccountWatchlistGroupSecurities`](#replaceaccountwatchlistgroupsecurities) | `replace_account_watchlist_group_securities` | `ReplaceAccountWatchlistGroupSecuritiesRequest` | `group_id`、`securities` | — | `AccountWatchlistGroups`（写操作） |
| [`ReplaceAccountWatchlistSecurities`](#replaceaccountwatchlistsecurities) | `replace_account_watchlist_securities` | `ReplaceAccountWatchlistSecuritiesRequest` | `version`、`securities` | — | `AccountWatchlistUpdate`（写操作） |

### 字段明细


<a id="accountpermissions"></a>
#### `AccountPermissions`

原生方法 `account_permissions`；返回 对象 `AccountPermissions`。

```go
client.AccountPermissions(ctx)
```

<a id="addaccountwatchlistgroupsecurities"></a>
#### `AddAccountWatchlistGroupSecurities`

原生方法 `add_account_watchlist_group_securities`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `group_id` | `GroupID` | `int` | **是** |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.AddAccountWatchlistGroupSecurities(ctx, thsdk.AddAccountWatchlistGroupSecuritiesRequest{GroupID: 1, Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="addaccountwatchlistsecurities"></a>
#### `AddAccountWatchlistSecurities`

原生方法 `add_account_watchlist_securities`；返回 `AccountWatchlistUpdate`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |
| `add_to_front` | `AddToFront` | `Opt` | 否 |

```go
client.AddAccountWatchlistSecurities(ctx, thsdk.AddAccountWatchlistSecuritiesRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="clearaccountwatchlist"></a>
#### `ClearAccountWatchlist`

原生方法 `clear_account_watchlist`；返回 `AccountWatchlistUpdate`（写操作）。

```go
client.ClearAccountWatchlist(ctx)
```

<a id="createaccountwatchlistgroup"></a>
#### `CreateAccountWatchlistGroup`

原生方法 `create_account_watchlist_group`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `name` | `Name` | `string` | **是** |
| `securities` | `Securities` | `[]Security` | 否 |

```go
client.CreateAccountWatchlistGroup(ctx, thsdk.CreateAccountWatchlistGroupRequest{Name: "研究池", Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="deleteaccountwatchlistgroup"></a>
#### `DeleteAccountWatchlistGroup`

原生方法 `delete_account_watchlist_group`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `group_id` | `GroupID` | `int` | **是** |

```go
client.DeleteAccountWatchlistGroup(ctx, thsdk.DeleteAccountWatchlistGroupRequest{GroupID: 1})
```

<a id="getaccountwatchlist"></a>
#### `GetAccountWatchlist`

原生方法 `get_account_watchlist`；返回 对象 `AccountWatchlistSnapshot`，行集在 `securities`。

```go
client.GetAccountWatchlist(ctx)
```

<a id="getaccountwatchlistgroups"></a>
#### `GetAccountWatchlistGroups`

原生方法 `get_account_watchlist_groups`；返回 对象 `AccountWatchlistGroups`，行集在 `groups`。

```go
client.GetAccountWatchlistGroups(ctx)
```

<a id="removeaccountwatchlistgroupsecurities"></a>
#### `RemoveAccountWatchlistGroupSecurities`

原生方法 `remove_account_watchlist_group_securities`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `group_id` | `GroupID` | `int` | **是** |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.RemoveAccountWatchlistGroupSecurities(ctx, thsdk.RemoveAccountWatchlistGroupSecuritiesRequest{GroupID: 1, Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="removeaccountwatchlistsecurities"></a>
#### `RemoveAccountWatchlistSecurities`

原生方法 `remove_account_watchlist_securities`；返回 `AccountWatchlistUpdate`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.RemoveAccountWatchlistSecurities(ctx, thsdk.RemoveAccountWatchlistSecuritiesRequest{Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="renameaccountwatchlistgroup"></a>
#### `RenameAccountWatchlistGroup`

原生方法 `rename_account_watchlist_group`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `group_id` | `GroupID` | `int` | **是** |
| `name` | `Name` | `string` | **是** |

```go
client.RenameAccountWatchlistGroup(ctx, thsdk.RenameAccountWatchlistGroupRequest{GroupID: 1, Name: "研究池"})
```

<a id="replaceaccountwatchlistgroupsecurities"></a>
#### `ReplaceAccountWatchlistGroupSecurities`

原生方法 `replace_account_watchlist_group_securities`；返回 `AccountWatchlistGroups`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `group_id` | `GroupID` | `int` | **是** |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.ReplaceAccountWatchlistGroupSecurities(ctx, thsdk.ReplaceAccountWatchlistGroupSecuritiesRequest{GroupID: 1, Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")}})
```

<a id="replaceaccountwatchlistsecurities"></a>
#### `ReplaceAccountWatchlistSecurities`

原生方法 `replace_account_watchlist_securities`；返回 `AccountWatchlistUpdate`（写操作）。

| 字段 | Go 成员 | 类型 | 必填 |
| --- | --- | --- | --- |
| `version` | `Version` | `int` | **是** |
| `securities` | `Securities` | `[]Security` | **是** |

```go
client.ReplaceAccountWatchlistSecurities(ctx, thsdk.ReplaceAccountWatchlistSecuritiesRequest{
    Version: 0,
    Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
})
```

<a id="convenience"></a>

## 便捷方法

对应上游 Python 的 12 个封装函数，已经做好参数换算和结果展开。

| Go 方法 | 上游 Python | 做什么 | 返回 |
| --- | --- | --- | --- |
| `GetPrice(ctx, PriceQuery)` | `get_price` / `klines` | 多周期 K 线，条数或日期范围、复权、列筛选 | `[]PriceRow` |
| `IntradayBars(ctx, security, date)` | `intraday_data` | 指定交易日分时；零值日期 = 当天 | `[]PriceRow` |
| `Ticks(ctx, TickQuery)` | `tick_level1` | 逐笔成交，按最近条数或时间范围 | `[]PriceRow` |
| `Depth(ctx, securities...)` | `depth` | 五档盘口，按买卖方向和档位展开 | `[]DepthRow` |
| `CorporateActions(ctx, security, count)` | `corporate_action` | 分红送转配股记录 | `[]PriceRow` |
| `SearchSymbols(ctx, pattern, market)` | `search_symbols` | 代码/拼音/名称搜索 | `[]PriceRow` |
| `CompleteTHSCodes(ctx, codes...)` | `complete_ths_code` | 短代码补全为完整证券 | `[]PriceRow` |
| `Wencai(ctx, condition, limit, markets...)` | `wencai_nlp` | 自然语言选股，动态字段摊成列 | `*WencaiSelection` |
| `News(ctx, page, pageSize, tag)` | `news` | 7x24 快讯分页 | `*NewsResult` |
| `MarketSecurities(ctx, market, start, count)` | `market_securities` | 市场证券分页 | `*NewsResult` |
| `BlockConstituents(ctx, block, start, count)` | `block_constituents` | 板块成分（ID 或名称） | `[]PriceRow` |

`PriceRow` 是 `map[string]json.RawMessage`，保留组件原样的键；要强类型就改用标准方法配
`KLine` / `Tick` / `OrderBook` 等结构体（定义见 [result.go](../result.go)）。

### 周期与复权

`PriceQuery.Frequency` 接受 `1m 5m 15m 30m 60m 120m daily 1d week weekly month monthly quarter year`；
`FQ` 接受 `""`、`none`、`pre`/`forward`、`post`/`backward`。分钟周期只支持相对条数（`Count`），
传日期范围会返回 `*ParameterError`——与上游限制一致。

