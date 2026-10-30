package thsdk

import "strings"

// Period is the KLine cycle passed to the runtime component.
type Period int

// KLine periods. Minute periods only support relative windows (negative start).
const (
	Period1Minute    Period = 0x3001
	Period5Minutes   Period = 0x3005
	Period15Minutes  Period = 0x300F
	Period30Minutes  Period = 0x301E
	Period60Minutes  Period = 0x303C
	Period120Minutes Period = 0x3078
	PeriodDay        Period = 0x4000
	PeriodWeek       Period = 0x5001
	PeriodMonth      Period = 0x6001
	PeriodQuarter    Period = 0x6003
	PeriodYear       Period = 0x7001
)

// Adjust is the price-adjustment mode for KLine queries.
type Adjust int

const (
	AdjustNone     Adjust = 0
	AdjustForward  Adjust = 1 // 前复权
	AdjustBackward Adjust = 2 // 后复权
)

// CallAuctionPhase selects the auction window to read.
type CallAuctionPhase string

const (
	PhaseOpening CallAuctionPhase = "opening"
	PhaseClosing CallAuctionPhase = "closing"
	PhaseCustom  CallAuctionPhase = "custom"
)

// ExtendedHoursSession selects a pre- or post-market session.
type ExtendedHoursSession string

const (
	SessionPreMarket  ExtendedHoursSession = "pre_market"
	SessionPostMarket ExtendedHoursSession = "post_market"
)

// OrderBookDepth is passed as a string, not a number.
type OrderBookDepth string

const DepthFive OrderBookDepth = "5"

// NewsFlashTraversalMode paginates the 7x24 flash feed.
type NewsFlashTraversalMode string

const (
	FlashPage  NewsFlashTraversalMode = "page"
	FlashOlder NewsFlashTraversalMode = "older"
	FlashNewer NewsFlashTraversalMode = "newer"
)

// NewsTraversalMode paginates article listings.
type NewsTraversalMode string

const (
	NewsPage  NewsTraversalMode = "page"
	NewsSince NewsTraversalMode = "since"
)

// NewsTimeline selects historical or live security news.
type NewsTimeline string

const (
	TimelineHistorical NewsTimeline = "historical"
	TimelineRealtime   NewsTimeline = "realtime"
)

// NewsDirection steps through a news cursor.
type NewsDirection string

const (
	DirectionFirst NewsDirection = "first"
	DirectionOlder NewsDirection = "older"
	DirectionNewer NewsDirection = "newer"
)

// RealtimeSecuritySignal filters securities by a live signal.
type RealtimeSecuritySignal string

const (
	SignalLimitUp        RealtimeSecuritySignal = "limit_up"
	SignalLimitDown      RealtimeSecuritySignal = "limit_down"
	SignalLimitUpBreak   RealtimeSecuritySignal = "limit_up_break"
	SignalLimitDownBreak RealtimeSecuritySignal = "limit_down_break"
)

// RealtimeStatsMetric identifies a computed realtime statistic.
type RealtimeStatsMetric int

const (
	MetricFirstLimitUpTime RealtimeStatsMetric = 330323
	MetricFinalLimitUpTime RealtimeStatsMetric = 330324
	MetricLimitUpType      RealtimeStatsMetric = 330325
	MetricLimitUpState     RealtimeStatsMetric = 330329
)

// OptionGreeksMetric identifies an option risk metric.
type OptionGreeksMetric int

const (
	GreeksImpliedVolatility OptionGreeksMetric = 920341
	GreeksHistoryVolatility OptionGreeksMetric = 723853
	GreeksRealLeverRatio    OptionGreeksMetric = 2980
	GreeksTheoryPrice       OptionGreeksMetric = 920342
	GreeksDelta             OptionGreeksMetric = 723854
	GreeksGamma             OptionGreeksMetric = 723855
	GreeksRho               OptionGreeksMetric = 723856
	GreeksTheta             OptionGreeksMetric = 723857
	GreeksVega              OptionGreeksMetric = 723858
)

// FinancialSnapshotField selects a share-count or per-share metric.
type FinancialSnapshotField string

const (
	FieldTotalShares      FinancialSnapshotField = "total_shares"
	FieldFloatShares      FinancialSnapshotField = "float_shares"
	FieldNetAssetPerShare FinancialSnapshotField = "net_asset_per_share"
)

// WencaiHotBlockKind chooses concept or industry hot boards.
type WencaiHotBlockKind string

const (
	HotBlockConcept  WencaiHotBlockKind = "concept"
	HotBlockIndustry WencaiHotBlockKind = "industry"
)

// HotEventNewsRange selects recent or archived event news.
type HotEventNewsRange string

const (
	RangeRecent     HotEventNewsRange = "recent"
	RangeHistorical HotEventNewsRange = "historical"
)

// ShortTermEventMode scopes market moving-event queries.
type ShortTermEventMode string

const (
	EventMarketRecent   ShortTermEventMode = "market_recent"
	EventSecurityScoped ShortTermEventMode = "security_scoped"
)

// PerformanceContributionTarget ranks index or board constituents.
type PerformanceContributionTarget string

const (
	ContributionIndex    PerformanceContributionTarget = "index"
	ContributionSecurity PerformanceContributionTarget = "security"
)

// PopularityRankTopType selects the popularity leaderboard.
type PopularityRankTopType int

const PopularityAllStocks PopularityRankTopType = 1

// InfoType is the news category channel id.
type InfoType int

const (
	InfoBad                    InfoType = -1
	InfoHistoryTime            InfoType = 2053
	InfoRealtime               InfoType = 2054
	InfoNewsReport             InfoType = 8193
	InfoShSzFund               InfoType = 10240
	InfoForeignExchangeNews    InfoType = 10241
	InfoForeignExchangeFX      InfoType = 10242
	InfoForeignExchangeCJ      InfoType = 10243
	InfoGlobalMarket           InfoType = 10244
	InfoGlobalFinance          InfoType = 10245
	InfoGlobalEconomicData     InfoType = 10246
	InfoStockInfo              InfoType = 14339
	InfoIndustryInfo           InfoType = 14340
	InfoAnnouncement           InfoType = 14341
	InfoStockResearchReport    InfoType = 14342
	InfoIndustryResearchReport InfoType = 14343
	InfoSecurityTapeReading    InfoType = 14348
	InfoIndexFutures           InfoType = 14349
	InfoMarketTapeReading      InfoType = 14355
	InfoMarketNews             InfoType = 14356
	InfoMarketTrendReading     InfoType = 14359
	InfoGLTMessage             InfoType = 14363
	InfoNewStockInfo           InfoType = 14364
	InfoMultiplyFutures        InfoType = 14365
	InfoNEEQNews               InfoType = 14366
	InfoNEEQSecondaryNews      InfoType = 14367
	InfoCommodityIndex         InfoType = 14368
	InfoSecurityInsight        InfoType = 14371
	InfoShToHKNews             InfoType = 32769
	InfoHKNews                 InfoType = 32770
	InfoHKPanInfo              InfoType = 32771
)

// InfoUSANews shares the channel id of InfoHKNews in the upstream SDK.
const InfoUSANews = InfoHKNews

// periodAliases maps the convenience frequency names onto Period ids.
var periodAliases = map[string]Period{
	"1m": Period1Minute, "5m": Period5Minutes, "15m": Period15Minutes,
	"30m": Period30Minutes, "60m": Period60Minutes, "120m": Period120Minutes,
	"day": PeriodDay, "daily": PeriodDay, "1d": PeriodDay,
	"week": PeriodWeek, "weekly": PeriodWeek,
	"month": PeriodMonth, "monthly": PeriodMonth,
	"quarter": PeriodQuarter, "year": PeriodYear,
}

// ParsePeriod accepts "daily", "1m", "week" and similar convenience names.
func ParsePeriod(frequency string) (Period, error) {
	period, ok := periodAliases[strings.ToLower(strings.TrimSpace(frequency))]
	if !ok {
		return 0, &ParameterError{Argument: "frequency", Message: "不支持的周期"}
	}
	return period, nil
}

var adjustAliases = map[string]Adjust{
	"":         AdjustNone,
	"none":     AdjustNone,
	"pre":      AdjustForward,
	"forward":  AdjustForward,
	"post":     AdjustBackward,
	"backward": AdjustBackward,
}

// ParseAdjust accepts "pre" / "post" / "none" and similar names.
func ParseAdjust(value string) (Adjust, error) {
	adjust, ok := adjustAliases[strings.ToLower(strings.TrimSpace(value))]
	if !ok {
		return 0, &ParameterError{Argument: "fq", Message: "不支持的复权方式"}
	}
	return adjust, nil
}
