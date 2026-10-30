// Code generated from the upstream thsdk method manifest. DO NOT EDIT.

package thsdk

import (
	"context"
)

// AccountPermissions calls the native account_permissions method (category account).
// Returns AccountPermissions.
func (c *Client) AccountPermissions(ctx context.Context) (RawData, error) {
	return c.request(ctx, "account_permissions", nil)
}

// AddAccountWatchlistGroupSecuritiesRequest is the request envelope of AddAccountWatchlistGroupSecurities.
type AddAccountWatchlistGroupSecuritiesRequest struct {
	GroupID    int        `json:"group_id"`
	Securities []Security `json:"securities"`
}

// AddAccountWatchlistGroupSecurities calls the native add_account_watchlist_group_securities method (category account).
// Required fields: group_id, securities.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) AddAccountWatchlistGroupSecurities(ctx context.Context, req AddAccountWatchlistGroupSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by add_account_watchlist_group_securities")
	}
	return c.request(ctx, "add_account_watchlist_group_securities", &req)
}

// AddAccountWatchlistSecuritiesRequest is the request envelope of AddAccountWatchlistSecurities.
type AddAccountWatchlistSecuritiesRequest struct {
	Securities []Security `json:"securities"`
	AddToFront Opt        `json:"add_to_front,omitempty"`
}

// AddAccountWatchlistSecurities calls the native add_account_watchlist_securities method (category account).
// Required fields: securities.
// Mutates account state and returns AccountWatchlistUpdate.
func (c *Client) AddAccountWatchlistSecurities(ctx context.Context, req AddAccountWatchlistSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by add_account_watchlist_securities")
	}
	return c.request(ctx, "add_account_watchlist_securities", &req)
}

// AnalyzeSecurityLimitUpRequest is the request envelope of AnalyzeSecurityLimitUp.
type AnalyzeSecurityLimitUpRequest struct {
	Security Security `json:"security"`
	Date     Opt      `json:"date,omitempty"`
}

// AnalyzeSecurityLimitUp calls the native analyze_security_limit_up method (category news).
// Required fields: security.
// Returns SecurityLimitUpAnalysisResult.
func (c *Client) AnalyzeSecurityLimitUp(ctx context.Context, req AnalyzeSecurityLimitUpRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by analyze_security_limit_up")
	}
	return c.request(ctx, "analyze_security_limit_up", &req)
}

// CalculateSecurityRealtimeStatisticsRequest is the request envelope of CalculateSecurityRealtimeStatistics.
type CalculateSecurityRealtimeStatisticsRequest struct {
	Securities []Security `json:"securities"`
	Metrics    any        `json:"metrics"`
}

// CalculateSecurityRealtimeStatistics calls the native calculate_security_realtime_statistics method (category rankings).
// Required fields: securities, metrics.
// Returns RealtimeStatsResult; the row set is under `rows`.
func (c *Client) CalculateSecurityRealtimeStatistics(ctx context.Context, req CalculateSecurityRealtimeStatisticsRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by calculate_security_realtime_statistics")
	}
	if req.Metrics == nil {
		return RawData{}, invalidArgument("metrics", "is required by calculate_security_realtime_statistics")
	}
	return c.request(ctx, "calculate_security_realtime_statistics", &req)
}

// CheckMarketTimelineReadinessRequest is the request envelope of CheckMarketTimelineReadiness.
type CheckMarketTimelineReadinessRequest struct {
	Market       string `json:"market"`
	TimeIndexDay string `json:"time_index_day"`
}

// CheckMarketTimelineReadiness calls the native check_market_timeline_readiness method (category market_data).
// Required fields: market, time_index_day.
// Returns MarketTimelineReadinessResult.
func (c *Client) CheckMarketTimelineReadiness(ctx context.Context, req CheckMarketTimelineReadinessRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by check_market_timeline_readiness")
	}
	if req.TimeIndexDay == "" {
		return RawData{}, invalidArgument("time_index_day", "is required by check_market_timeline_readiness")
	}
	return c.request(ctx, "check_market_timeline_readiness", &req)
}

// ClearAccountWatchlist calls the native clear_account_watchlist method (category account).
// Mutates account state and returns AccountWatchlistUpdate.
func (c *Client) ClearAccountWatchlist(ctx context.Context) (RawData, error) {
	return c.request(ctx, "clear_account_watchlist", nil)
}

// CreateAccountWatchlistGroupRequest is the request envelope of CreateAccountWatchlistGroup.
type CreateAccountWatchlistGroupRequest struct {
	Name       string     `json:"name"`
	Securities []Security `json:"securities"`
}

// CreateAccountWatchlistGroup calls the native create_account_watchlist_group method (category account).
// Required fields: name.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) CreateAccountWatchlistGroup(ctx context.Context, req CreateAccountWatchlistGroupRequest) (RawData, error) {
	if req.Name == "" {
		return RawData{}, invalidArgument("name", "is required by create_account_watchlist_group")
	}
	return c.request(ctx, "create_account_watchlist_group", &req)
}

// DeleteAccountWatchlistGroupRequest is the request envelope of DeleteAccountWatchlistGroup.
type DeleteAccountWatchlistGroupRequest struct {
	GroupID int `json:"group_id"`
}

// DeleteAccountWatchlistGroup calls the native delete_account_watchlist_group method (category account).
// Required fields: group_id.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) DeleteAccountWatchlistGroup(ctx context.Context, req DeleteAccountWatchlistGroupRequest) (RawData, error) {
	return c.request(ctx, "delete_account_watchlist_group", &req)
}

// GetAccountWatchlist calls the native get_account_watchlist method (category account).
// Returns AccountWatchlistSnapshot; the row set is under `securities`.
func (c *Client) GetAccountWatchlist(ctx context.Context) (RawData, error) {
	return c.request(ctx, "get_account_watchlist", nil)
}

// GetAccountWatchlistGroups calls the native get_account_watchlist_groups method (category account).
// Returns AccountWatchlistGroups; the row set is under `groups`.
func (c *Client) GetAccountWatchlistGroups(ctx context.Context) (RawData, error) {
	return c.request(ctx, "get_account_watchlist_groups", nil)
}

// GetMarketMetadataRequest is the request envelope of GetMarketMetadata.
type GetMarketMetadataRequest struct {
	Market      string `json:"market"`
	VersionFlag Opt    `json:"version_flag,omitempty"`
}

// GetMarketMetadata calls the native get_market_metadata method (category securities).
// Required fields: market.
// Returns MarketMetadata; the row set is under `sections`.
func (c *Client) GetMarketMetadata(ctx context.Context, req GetMarketMetadataRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by get_market_metadata")
	}
	return c.request(ctx, "get_market_metadata", &req)
}

// GetMarketSecurityNamesRequest is the request envelope of GetMarketSecurityNames.
type GetMarketSecurityNamesRequest struct {
	Market             string `json:"market"`
	RealtimeVersionIni Opt    `json:"realtime_version_ini,omitempty"`
	HistoryVersionIni  Opt    `json:"history_version_ini,omitempty"`
	BaseVersionIni     Opt    `json:"base_version_ini,omitempty"`
}

// GetMarketSecurityNames calls the native get_market_security_names method (category securities).
// Required fields: market.
// Returns MarketSecurityNamesResult; the row set is under `items`.
func (c *Client) GetMarketSecurityNames(ctx context.Context, req GetMarketSecurityNamesRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by get_market_security_names")
	}
	return c.request(ctx, "get_market_security_names", &req)
}

// GetSecurityConceptTagsRequest is the request envelope of GetSecurityConceptTags.
type GetSecurityConceptTagsRequest struct {
	Security Security `json:"security"`
}

// GetSecurityConceptTags calls the native get_security_concept_tags method (category securities).
// Required fields: security.
// Returns SecurityConceptTagsResult; the row set is under `tags`.
func (c *Client) GetSecurityConceptTags(ctx context.Context, req GetSecurityConceptTagsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_concept_tags")
	}
	return c.request(ctx, "get_security_concept_tags", &req)
}

// GetSecurityIndustryRequest is the request envelope of GetSecurityIndustry.
type GetSecurityIndustryRequest struct {
	Security Security `json:"security"`
}

// GetSecurityIndustry calls the native get_security_industry method (category securities).
// Required fields: security.
// Returns SecurityIndustryMapping.
func (c *Client) GetSecurityIndustry(ctx context.Context, req GetSecurityIndustryRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_industry")
	}
	return c.request(ctx, "get_security_industry", &req)
}

// GetSecurityPopularityRankRequest is the request envelope of GetSecurityPopularityRank.
type GetSecurityPopularityRankRequest struct {
	Security  Security `json:"security"`
	Benchmark any      `json:"benchmark"`
}

// GetSecurityPopularityRank calls the native get_security_popularity_rank method (category rankings).
// Required fields: security.
// Returns PopularityRankResult.
func (c *Client) GetSecurityPopularityRank(ctx context.Context, req GetSecurityPopularityRankRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_popularity_rank")
	}
	return c.request(ctx, "get_security_popularity_rank", &req)
}

// GetSecurityPriceLimitEventsRequest is the request envelope of GetSecurityPriceLimitEvents.
type GetSecurityPriceLimitEventsRequest struct {
	Security  Security `json:"security"`
	StartYear Opt      `json:"start_year,omitempty"`
	EndYear   Opt      `json:"end_year,omitempty"`
}

// GetSecurityPriceLimitEvents calls the native get_security_price_limit_events method (category news).
// Required fields: security.
// Returns PriceLimitEventsResult; the row set is under `events`.
func (c *Client) GetSecurityPriceLimitEvents(ctx context.Context, req GetSecurityPriceLimitEventsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_price_limit_events")
	}
	return c.request(ctx, "get_security_price_limit_events", &req)
}

// GetSecurityShortTermHighlightsRequest is the request envelope of GetSecurityShortTermHighlights.
type GetSecurityShortTermHighlightsRequest struct {
	Security Security `json:"security"`
}

// GetSecurityShortTermHighlights calls the native get_security_short_term_highlights method (category news).
// Required fields: security.
// Returns SecurityShortTermHighlightsResult.
func (c *Client) GetSecurityShortTermHighlights(ctx context.Context, req GetSecurityShortTermHighlightsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_short_term_highlights")
	}
	return c.request(ctx, "get_security_short_term_highlights", &req)
}

// GetSecurityTradingTimelineRequest is the request envelope of GetSecurityTradingTimeline.
type GetSecurityTradingTimelineRequest struct {
	Security Security `json:"security"`
}

// GetSecurityTradingTimeline calls the native get_security_trading_timeline method (category market_data).
// Required fields: security.
// Returns SecurityTradingTimeline; the row set is under `items`.
func (c *Client) GetSecurityTradingTimeline(ctx context.Context, req GetSecurityTradingTimelineRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by get_security_trading_timeline")
	}
	return c.request(ctx, "get_security_trading_timeline", &req)
}

// ListBlockConstituentsRequest is the request envelope of ListBlockConstituents.
type ListBlockConstituentsRequest struct {
	Block     Block `json:"block"`
	SortBegin Opt   `json:"sort_begin,omitempty"`
	SortCount Opt   `json:"sort_count,omitempty"`
	SortOrder Opt   `json:"sort_order,omitempty"`
	SortID    Opt   `json:"sort_id,omitempty"`
}

// ListBlockConstituents calls the native list_block_constituents method (category securities).
// Required fields: block.
// Returns a JSON array of BlockConstituent objects.
func (c *Client) ListBlockConstituents(ctx context.Context, req ListBlockConstituentsRequest) (RawData, error) {
	if req.Block.ID == nil && req.Block.Name == "" {
		return RawData{}, invalidArgument("block", "is required by list_block_constituents")
	}
	return c.request(ctx, "list_block_constituents", &req)
}

// ListBlockDescriptionsRequest is the request envelope of ListBlockDescriptions.
type ListBlockDescriptionsRequest struct {
	Blocks []Block `json:"blocks"`
}

// ListBlockDescriptions calls the native list_block_descriptions method (category securities).
// Required fields: blocks.
// Returns a JSON array of BlockDescription objects.
func (c *Client) ListBlockDescriptions(ctx context.Context, req ListBlockDescriptionsRequest) (RawData, error) {
	if len(req.Blocks) == 0 {
		return RawData{}, invalidArgument("blocks", "is required by list_block_descriptions")
	}
	return c.request(ctx, "list_block_descriptions", &req)
}

// ListCommodityStockLinkageNewsRequest is the request envelope of ListCommodityStockLinkageNews.
type ListCommodityStockLinkageNewsRequest struct {
	CommodityIndex any `json:"commodity_index"`
}

// ListCommodityStockLinkageNews calls the native list_commodity_stock_linkage_news method (category news).
// Required fields: commodity_index.
// Returns CommodityStockLinkageNewsResult; the row set is under `items`.
func (c *Client) ListCommodityStockLinkageNews(ctx context.Context, req ListCommodityStockLinkageNewsRequest) (RawData, error) {
	if req.CommodityIndex == nil {
		return RawData{}, invalidArgument("commodity_index", "is required by list_commodity_stock_linkage_news")
	}
	return c.request(ctx, "list_commodity_stock_linkage_news", &req)
}

// ListFuturesRelatedSecurities calls the native list_futures_related_securities method (category securities).
// Returns a JSON array of FuturesRelatedSecurity objects.
func (c *Client) ListFuturesRelatedSecurities(ctx context.Context) (RawData, error) {
	return c.request(ctx, "list_futures_related_securities", nil)
}

// ListHotBlockCalendarEvents calls the native list_hot_block_calendar_events method (category news).
// Returns HotBlockCalendarResult; the row set is under `items`.
func (c *Client) ListHotBlockCalendarEvents(ctx context.Context) (RawData, error) {
	return c.request(ctx, "list_hot_block_calendar_events", nil)
}

// ListHotEventNewsRequest is the request envelope of ListHotEventNews.
type ListHotEventNewsRequest struct {
	Block Block `json:"block"`
	Range Opt   `json:"range,omitempty"`
}

// ListHotEventNews calls the native list_hot_event_news method (category news).
// Required fields: block.
// Returns HotEventNewsResult; the row set is under `items`.
func (c *Client) ListHotEventNews(ctx context.Context, req ListHotEventNewsRequest) (RawData, error) {
	if req.Block.ID == nil && req.Block.Name == "" {
		return RawData{}, invalidArgument("block", "is required by list_hot_event_news")
	}
	return c.request(ctx, "list_hot_event_news", &req)
}

// ListIndexCallAuctionQuotesRequest is the request envelope of ListIndexCallAuctionQuotes.
type ListIndexCallAuctionQuotesRequest struct {
	Index any    `json:"index"`
	Phase string `json:"phase"`
}

// ListIndexCallAuctionQuotes calls the native list_index_call_auction_quotes method (category market_data).
// Required fields: index, phase.
// Returns a JSON array of IndexCallAuctionQuote objects.
func (c *Client) ListIndexCallAuctionQuotes(ctx context.Context, req ListIndexCallAuctionQuotesRequest) (RawData, error) {
	if req.Index == nil {
		return RawData{}, invalidArgument("index", "is required by list_index_call_auction_quotes")
	}
	if req.Phase == "" {
		return RawData{}, invalidArgument("phase", "is required by list_index_call_auction_quotes")
	}
	return c.request(ctx, "list_index_call_auction_quotes", &req)
}

// ListIndustryChildrenRequest is the request envelope of ListIndustryChildren.
type ListIndustryChildrenRequest struct {
	Industry any `json:"industry"`
}

// ListIndustryChildren calls the native list_industry_children method (category securities).
// Required fields: industry.
// Returns a JSON array of IndustryChild objects.
func (c *Client) ListIndustryChildren(ctx context.Context, req ListIndustryChildrenRequest) (RawData, error) {
	if req.Industry == nil {
		return RawData{}, invalidArgument("industry", "is required by list_industry_children")
	}
	return c.request(ctx, "list_industry_children", &req)
}

// ListMarketCallAuctionUnusualEventsRequest is the request envelope of ListMarketCallAuctionUnusualEvents.
type ListMarketCallAuctionUnusualEventsRequest struct {
	Market    string `json:"market"`
	TradeDate Opt    `json:"trade_date,omitempty"`
}

// ListMarketCallAuctionUnusualEvents calls the native list_market_call_auction_unusual_events method (category news).
// Required fields: market.
// Returns a JSON array of ShortTermEvent objects.
func (c *Client) ListMarketCallAuctionUnusualEvents(ctx context.Context, req ListMarketCallAuctionUnusualEventsRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by list_market_call_auction_unusual_events")
	}
	return c.request(ctx, "list_market_call_auction_unusual_events", &req)
}

// ListMarketChangeEvents calls the native list_market_change_events method (category news).
// Returns a JSON array of MarketChangeEvent objects.
func (c *Client) ListMarketChangeEvents(ctx context.Context) (RawData, error) {
	return c.request(ctx, "list_market_change_events", nil)
}

// ListMarketSecuritiesRequest is the request envelope of ListMarketSecurities.
type ListMarketSecuritiesRequest struct {
	Market     string `json:"market"`
	SortBegin  Opt    `json:"sort_begin,omitempty"`
	SortCount  Opt    `json:"sort_count,omitempty"`
	SortID     Opt    `json:"sort_id,omitempty"`
	SortOrder  Opt    `json:"sort_order,omitempty"`
	FuncPeriod Opt    `json:"func_period,omitempty"`
}

// ListMarketSecurities calls the native list_market_securities method (category securities).
// Required fields: market.
// Returns MarketSecuritiesResult; the row set is under `securities`.
func (c *Client) ListMarketSecurities(ctx context.Context, req ListMarketSecuritiesRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by list_market_securities")
	}
	return c.request(ctx, "list_market_securities", &req)
}

// ListMarketShortTermEventsRequest is the request envelope of ListMarketShortTermEvents.
type ListMarketShortTermEventsRequest struct {
	Market                string     `json:"market"`
	Mode                  string     `json:"mode"`
	Securities            []Security `json:"securities"`
	MaxCount              Opt        `json:"max_count,omitempty"`
	BeforeTimestampMicros Opt        `json:"before_timestamp_micros,omitempty"`
}

// ListMarketShortTermEvents calls the native list_market_short_term_events method (category news).
// Required fields: market, mode.
// Returns a JSON array of ShortTermEvent objects.
func (c *Client) ListMarketShortTermEvents(ctx context.Context, req ListMarketShortTermEventsRequest) (RawData, error) {
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by list_market_short_term_events")
	}
	if req.Mode == "" {
		return RawData{}, invalidArgument("mode", "is required by list_market_short_term_events")
	}
	return c.request(ctx, "list_market_short_term_events", &req)
}

// ListNewsFlashItemsRequest is the request envelope of ListNewsFlashItems.
type ListNewsFlashItemsRequest struct {
	Traversal   any `json:"traversal"`
	Tag         Opt `json:"tag,omitempty"`
	TagID       Opt `json:"tag_id,omitempty"`
	Page        Opt `json:"page,omitempty"`
	PageSize    Opt `json:"page_size,omitempty"`
	PageTime    Opt `json:"page_time,omitempty"`
	Sequence    Opt `json:"sequence,omitempty"`
	CreatedAt   Opt `json:"created_at,omitempty"`
	Environment Opt `json:"environment,omitempty"`
}

// ListNewsFlashItems calls the native list_news_flash_items method (category news).
// Required fields: traversal.
// Returns NewsFlashItemsResult; the row set is under `items`.
func (c *Client) ListNewsFlashItems(ctx context.Context, req ListNewsFlashItemsRequest) (RawData, error) {
	if req.Traversal == nil {
		return RawData{}, invalidArgument("traversal", "is required by list_news_flash_items")
	}
	return c.request(ctx, "list_news_flash_items", &req)
}

// ListNewsItemsRequest is the request envelope of ListNewsItems.
type ListNewsItemsRequest struct {
	Category    int      `json:"category"`
	Security    Security `json:"security"`
	Market      Opt      `json:"market,omitempty"`
	Traversal   any      `json:"traversal"`
	SummaryMode Opt      `json:"summary_mode,omitempty"`
	Advertorial Opt      `json:"advertorial,omitempty"`
}

// ListNewsItems calls the native list_news_items method (category news).
// Required fields: category, traversal.
// Returns NewsItemsResult; the row set is under `items`.
func (c *Client) ListNewsItems(ctx context.Context, req ListNewsItemsRequest) (RawData, error) {
	if req.Traversal == nil {
		return RawData{}, invalidArgument("traversal", "is required by list_news_items")
	}
	return c.request(ctx, "list_news_items", &req)
}

// ListRelatedSecurityPerformancesRequest is the request envelope of ListRelatedSecurityPerformances.
type ListRelatedSecurityPerformancesRequest struct {
	Security Security `json:"security"`
}

// ListRelatedSecurityPerformances calls the native list_related_security_performances method (category securities).
// Required fields: security.
// Returns a JSON array of RelatedSecurityPerformance objects.
func (c *Client) ListRelatedSecurityPerformances(ctx context.Context, req ListRelatedSecurityPerformancesRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_related_security_performances")
	}
	return c.request(ctx, "list_related_security_performances", &req)
}

// ListSecuritiesByRealtimeSignalRequest is the request envelope of ListSecuritiesByRealtimeSignal.
type ListSecuritiesByRealtimeSignalRequest struct {
	Markets []string `json:"markets"`
	Signal  string   `json:"signal"`
}

// ListSecuritiesByRealtimeSignal calls the native list_securities_by_realtime_signal method (category rankings).
// Required fields: markets, signal.
// Returns a JSON array of RealtimeSignalSecurity objects.
func (c *Client) ListSecuritiesByRealtimeSignal(ctx context.Context, req ListSecuritiesByRealtimeSignalRequest) (RawData, error) {
	if len(req.Markets) == 0 {
		return RawData{}, invalidArgument("markets", "is required by list_securities_by_realtime_signal")
	}
	if req.Signal == "" {
		return RawData{}, invalidArgument("signal", "is required by list_securities_by_realtime_signal")
	}
	return c.request(ctx, "list_securities_by_realtime_signal", &req)
}

// ListSecurityAHRelationsRequest is the request envelope of ListSecurityAHRelations.
type ListSecurityAHRelationsRequest struct {
	Security Security `json:"security"`
}

// ListSecurityAHRelations calls the native list_security_ah_relations method (category securities).
// Required fields: security.
// Returns a JSON array of SecurityAHRelation objects.
func (c *Client) ListSecurityAHRelations(ctx context.Context, req ListSecurityAHRelationsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_ah_relations")
	}
	return c.request(ctx, "list_security_ah_relations", &req)
}

// ListSecurityBlockMembershipsRequest is the request envelope of ListSecurityBlockMemberships.
type ListSecurityBlockMembershipsRequest struct {
	Security Security `json:"security"`
}

// ListSecurityBlockMemberships calls the native list_security_block_memberships method (category securities).
// Required fields: security.
// Returns a JSON array of SecurityBlockMembership objects.
func (c *Client) ListSecurityBlockMemberships(ctx context.Context, req ListSecurityBlockMembershipsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_block_memberships")
	}
	return c.request(ctx, "list_security_block_memberships", &req)
}

// ListSecurityCallAuctionQuotesRequest is the request envelope of ListSecurityCallAuctionQuotes.
type ListSecurityCallAuctionQuotesRequest struct {
	Security Security `json:"security"`
	Phase    string   `json:"phase"`
	Date     Opt      `json:"date,omitempty"`
	Window   any      `json:"window"`
}

// ListSecurityCallAuctionQuotes calls the native list_security_call_auction_quotes method (category market_data).
// Required fields: security, phase.
// Returns a JSON array of CallAuctionQuote objects.
func (c *Client) ListSecurityCallAuctionQuotes(ctx context.Context, req ListSecurityCallAuctionQuotesRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_call_auction_quotes")
	}
	if req.Phase == "" {
		return RawData{}, invalidArgument("phase", "is required by list_security_call_auction_quotes")
	}
	return c.request(ctx, "list_security_call_auction_quotes", &req)
}

// ListSecurityCorporateActionsRequest is the request envelope of ListSecurityCorporateActions.
type ListSecurityCorporateActionsRequest struct {
	Security Security `json:"security"`
	Start    int      `json:"start"`
	End      int      `json:"end"`
	Adjust   Adjust   `json:"adjust"`
	Period   Period   `json:"period"`
}

// ListSecurityCorporateActions calls the native list_security_corporate_actions method (category market_data).
// Required fields: security, start, end, adjust, period.
// Returns a JSON array of CorporateAction objects.
func (c *Client) ListSecurityCorporateActions(ctx context.Context, req ListSecurityCorporateActionsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_corporate_actions")
	}
	return c.request(ctx, "list_security_corporate_actions", &req)
}

// ListSecurityDailyCapitalFlowsRequest is the request envelope of ListSecurityDailyCapitalFlows.
type ListSecurityDailyCapitalFlowsRequest struct {
	Securities []Security `json:"securities"`
}

// ListSecurityDailyCapitalFlows calls the native list_security_daily_capital_flows method (category market_data).
// Required fields: securities.
// Returns a JSON array of DailyCapitalFlow objects.
func (c *Client) ListSecurityDailyCapitalFlows(ctx context.Context, req ListSecurityDailyCapitalFlowsRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by list_security_daily_capital_flows")
	}
	return c.request(ctx, "list_security_daily_capital_flows", &req)
}

// ListSecurityDailyKLinesWithPreviousCloseRequest is the request envelope of ListSecurityDailyKLinesWithPreviousClose.
type ListSecurityDailyKLinesWithPreviousCloseRequest struct {
	Security Security `json:"security"`
	Start    int      `json:"start"`
	End      int      `json:"end"`
	Adjust   Adjust   `json:"adjust"`
}

// ListSecurityDailyKLinesWithPreviousClose calls the native list_security_daily_k_lines_with_previous_close method (category market_data).
// Required fields: security, start, end, adjust.
// Returns a JSON array of DailyKLineWithPreviousClose objects.
func (c *Client) ListSecurityDailyKLinesWithPreviousClose(ctx context.Context, req ListSecurityDailyKLinesWithPreviousCloseRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_daily_k_lines_with_previous_close")
	}
	return c.request(ctx, "list_security_daily_k_lines_with_previous_close", &req)
}

// ListSecurityExtendedHoursIntradayBarsRequest is the request envelope of ListSecurityExtendedHoursIntradayBars.
type ListSecurityExtendedHoursIntradayBarsRequest struct {
	Security Security `json:"security"`
	Session  string   `json:"session"`
	Date     string   `json:"date"`
}

// ListSecurityExtendedHoursIntradayBars calls the native list_security_extended_hours_intraday_bars method (category market_data).
// Required fields: security, session, date.
// Returns a JSON array of IntradayBar objects.
func (c *Client) ListSecurityExtendedHoursIntradayBars(ctx context.Context, req ListSecurityExtendedHoursIntradayBarsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_extended_hours_intraday_bars")
	}
	if req.Session == "" {
		return RawData{}, invalidArgument("session", "is required by list_security_extended_hours_intraday_bars")
	}
	if req.Date == "" {
		return RawData{}, invalidArgument("date", "is required by list_security_extended_hours_intraday_bars")
	}
	return c.request(ctx, "list_security_extended_hours_intraday_bars", &req)
}

// ListSecurityExtendedHoursTicksRequest is the request envelope of ListSecurityExtendedHoursTicks.
type ListSecurityExtendedHoursTicksRequest struct {
	Security Security `json:"security"`
	Session  string   `json:"session"`
	Window   any      `json:"window"`
}

// ListSecurityExtendedHoursTicks calls the native list_security_extended_hours_ticks method (category market_data).
// Required fields: security, session, window.
// Returns a JSON array of Tick objects.
func (c *Client) ListSecurityExtendedHoursTicks(ctx context.Context, req ListSecurityExtendedHoursTicksRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_extended_hours_ticks")
	}
	if req.Session == "" {
		return RawData{}, invalidArgument("session", "is required by list_security_extended_hours_ticks")
	}
	if req.Window == nil {
		return RawData{}, invalidArgument("window", "is required by list_security_extended_hours_ticks")
	}
	return c.request(ctx, "list_security_extended_hours_ticks", &req)
}

// ListSecurityFinancialSnapshotsRequest is the request envelope of ListSecurityFinancialSnapshots.
type ListSecurityFinancialSnapshotsRequest struct {
	Securities []Security `json:"securities"`
	Fields     any        `json:"fields"`
}

// ListSecurityFinancialSnapshots calls the native list_security_financial_snapshots method (category market_data).
// Required fields: securities, fields.
// Returns a JSON array of SecurityFinancialSnapshot objects.
func (c *Client) ListSecurityFinancialSnapshots(ctx context.Context, req ListSecurityFinancialSnapshotsRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by list_security_financial_snapshots")
	}
	if req.Fields == nil {
		return RawData{}, invalidArgument("fields", "is required by list_security_financial_snapshots")
	}
	return c.request(ctx, "list_security_financial_snapshots", &req)
}

// ListSecurityFuturesRelations calls the native list_security_futures_relations method (category securities).
// Required fields: securities.
// Returns a JSON array of SecurityFuturesRelation objects.
func (c *Client) ListSecurityFuturesRelations(ctx context.Context, securities []Security) (RawData, error) {
	list, err := securityList(securities)
	if err != nil {
		return RawData{}, err
	}
	return c.request(ctx, "list_security_futures_relations", list)
}

// ListSecurityIndustryMappingsRequest is the request envelope of ListSecurityIndustryMappings.
type ListSecurityIndustryMappingsRequest struct {
	Securities []Security `json:"securities"`
}

// ListSecurityIndustryMappings calls the native list_security_industry_mappings method (category securities).
// Required fields: securities.
// Returns a JSON array of SecurityIndustryMapping objects.
func (c *Client) ListSecurityIndustryMappings(ctx context.Context, req ListSecurityIndustryMappingsRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by list_security_industry_mappings")
	}
	return c.request(ctx, "list_security_industry_mappings", &req)
}

// ListSecurityIntradayBarsRequest is the request envelope of ListSecurityIntradayBars.
type ListSecurityIntradayBarsRequest struct {
	Security Security `json:"security"`
	Date     string   `json:"date"`
}

// ListSecurityIntradayBars calls the native list_security_intraday_bars method (category market_data).
// Required fields: security, date.
// Returns a JSON array of IntradayBar objects.
func (c *Client) ListSecurityIntradayBars(ctx context.Context, req ListSecurityIntradayBarsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_intraday_bars")
	}
	if req.Date == "" {
		return RawData{}, invalidArgument("date", "is required by list_security_intraday_bars")
	}
	return c.request(ctx, "list_security_intraday_bars", &req)
}

// ListSecurityKLinesRequest is the request envelope of ListSecurityKLines.
type ListSecurityKLinesRequest struct {
	Security Security `json:"security"`
	Start    int      `json:"start"`
	End      int      `json:"end"`
	Adjust   Adjust   `json:"adjust"`
	Period   Period   `json:"period"`
}

// ListSecurityKLines calls the native list_security_k_lines method (category market_data).
// Required fields: security, start, end, adjust, period.
// Returns a JSON array of KLine objects.
func (c *Client) ListSecurityKLines(ctx context.Context, req ListSecurityKLinesRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_k_lines")
	}
	return c.request(ctx, "list_security_k_lines", &req)
}

// ListSecurityLinkRelations calls the native list_security_link_relations method (category securities).
// Required fields: link_key.
// Returns a JSON array of SecurityLinkRelation objects.
func (c *Client) ListSecurityLinkRelations(ctx context.Context, value string) (RawData, error) {
	if err := requireText("link_key", value); err != nil {
		return RawData{}, err
	}
	return c.request(ctx, "list_security_link_relations", value)
}

// ListSecurityNewsEventsRequest is the request envelope of ListSecurityNewsEvents.
type ListSecurityNewsEventsRequest struct {
	Security Security `json:"security"`
	Timeline string   `json:"timeline"`
	Window   any      `json:"window"`
}

// ListSecurityNewsEvents calls the native list_security_news_events method (category news).
// Required fields: security, timeline, window.
// Returns SecurityNewsItemsResult; the row set is under `items`.
func (c *Client) ListSecurityNewsEvents(ctx context.Context, req ListSecurityNewsEventsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_news_events")
	}
	if req.Timeline == "" {
		return RawData{}, invalidArgument("timeline", "is required by list_security_news_events")
	}
	if req.Window == nil {
		return RawData{}, invalidArgument("window", "is required by list_security_news_events")
	}
	return c.request(ctx, "list_security_news_events", &req)
}

// ListSecurityNewsMarkersRequest is the request envelope of ListSecurityNewsMarkers.
type ListSecurityNewsMarkersRequest struct {
	Security Security `json:"security"`
	Timeline string   `json:"timeline"`
}

// ListSecurityNewsMarkers calls the native list_security_news_markers method (category news).
// Required fields: security, timeline.
// Returns SecurityNewsItemsResult; the row set is under `items`.
func (c *Client) ListSecurityNewsMarkers(ctx context.Context, req ListSecurityNewsMarkersRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_news_markers")
	}
	if req.Timeline == "" {
		return RawData{}, invalidArgument("timeline", "is required by list_security_news_markers")
	}
	return c.request(ctx, "list_security_news_markers", &req)
}

// ListSecurityOptionGreeksRequest is the request envelope of ListSecurityOptionGreeks.
type ListSecurityOptionGreeksRequest struct {
	Securities []Security `json:"securities"`
	Metrics    any        `json:"metrics"`
}

// ListSecurityOptionGreeks calls the native list_security_option_greeks method (category market_data).
// Required fields: securities.
// Returns a JSON array of OptionGreeksRow objects.
func (c *Client) ListSecurityOptionGreeks(ctx context.Context, req ListSecurityOptionGreeksRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by list_security_option_greeks")
	}
	return c.request(ctx, "list_security_option_greeks", &req)
}

// ListSecurityOrderBooksRequest is the request envelope of ListSecurityOrderBooks.
type ListSecurityOrderBooksRequest struct {
	Securities []Security     `json:"securities"`
	Depth      OrderBookDepth `json:"depth"`
}

// ListSecurityOrderBooks calls the native list_security_order_books method (category market_data).
// Required fields: securities, depth.
// Returns a JSON array of SecurityOrderBook objects.
func (c *Client) ListSecurityOrderBooks(ctx context.Context, req ListSecurityOrderBooksRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by list_security_order_books")
	}
	return c.request(ctx, "list_security_order_books", &req)
}

// ListSecurityPriceVolumeLevelsRequest is the request envelope of ListSecurityPriceVolumeLevels.
type ListSecurityPriceVolumeLevelsRequest struct {
	Security Security `json:"security"`
}

// ListSecurityPriceVolumeLevels calls the native list_security_price_volume_levels method (category market_data).
// Required fields: security.
// Returns a JSON array of PriceVolumeLevel objects.
func (c *Client) ListSecurityPriceVolumeLevels(ctx context.Context, req ListSecurityPriceVolumeLevelsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_price_volume_levels")
	}
	return c.request(ctx, "list_security_price_volume_levels", &req)
}

// ListSecurityResearchReportsRequest is the request envelope of ListSecurityResearchReports.
type ListSecurityResearchReportsRequest struct {
	Security  Security `json:"security"`
	Traversal any      `json:"traversal"`
}

// ListSecurityResearchReports calls the native list_security_research_reports method (category news).
// Required fields: security, traversal.
// Returns SecurityResearchReportsResult; the row set is under `items`.
func (c *Client) ListSecurityResearchReports(ctx context.Context, req ListSecurityResearchReportsRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_research_reports")
	}
	if req.Traversal == nil {
		return RawData{}, invalidArgument("traversal", "is required by list_security_research_reports")
	}
	return c.request(ctx, "list_security_research_reports", &req)
}

// ListSecurityTicksRequest is the request envelope of ListSecurityTicks.
type ListSecurityTicksRequest struct {
	Security Security `json:"security"`
	Window   any      `json:"window"`
}

// ListSecurityTicks calls the native list_security_ticks method (category market_data).
// Required fields: security, window.
// Returns a JSON array of Tick objects.
func (c *Client) ListSecurityTicks(ctx context.Context, req ListSecurityTicksRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_ticks")
	}
	if req.Window == nil {
		return RawData{}, invalidArgument("window", "is required by list_security_ticks")
	}
	return c.request(ctx, "list_security_ticks", &req)
}

// ListSecurityTimeAndSalesTicksRequest is the request envelope of ListSecurityTimeAndSalesTicks.
type ListSecurityTimeAndSalesTicksRequest struct {
	Security Security `json:"security"`
	Count    int      `json:"count"`
}

// ListSecurityTimeAndSalesTicks calls the native list_security_time_and_sales_ticks method (category market_data).
// Required fields: security, count.
// Returns a JSON array of Tick objects.
func (c *Client) ListSecurityTimeAndSalesTicks(ctx context.Context, req ListSecurityTimeAndSalesTicksRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by list_security_time_and_sales_ticks")
	}
	return c.request(ctx, "list_security_time_and_sales_ticks", &req)
}

// ListWencaiExpressionFields calls the native list_wencai_expression_fields method (category wencai).
// Required fields: query.
// Returns a JSON array of WencaiField objects.
func (c *Client) ListWencaiExpressionFields(ctx context.Context, value string) (RawData, error) {
	if err := requireText("query", value); err != nil {
		return RawData{}, err
	}
	return c.request(ctx, "list_wencai_expression_fields", value)
}

// ListWencaiHotBlocksRequest is the request envelope of ListWencaiHotBlocks.
type ListWencaiHotBlocksRequest struct {
	Kind string `json:"kind"`
}

// ListWencaiHotBlocks calls the native list_wencai_hot_blocks method (category wencai).
// Required fields: kind.
// Returns WencaiHotBlocksResult; the row set is under `items`.
func (c *Client) ListWencaiHotBlocks(ctx context.Context, req ListWencaiHotBlocksRequest) (RawData, error) {
	if req.Kind == "" {
		return RawData{}, invalidArgument("kind", "is required by list_wencai_hot_blocks")
	}
	return c.request(ctx, "list_wencai_hot_blocks", &req)
}

// QueryWencaiRequest is the request envelope of QueryWencai.
type QueryWencaiRequest struct {
	Query   string   `json:"query"`
	Markets []string `json:"markets"`
	Limit   Opt      `json:"limit,omitempty"`
}

// QueryWencai calls the native query_wencai method (category wencai).
// Required fields: query.
// Returns WencaiQueryResult; the row set is under `rows`.
func (c *Client) QueryWencai(ctx context.Context, req QueryWencaiRequest) (RawData, error) {
	if req.Query == "" {
		return RawData{}, invalidArgument("query", "is required by query_wencai")
	}
	return c.request(ctx, "query_wencai", &req)
}

// QueryWencaiRealtimeFieldsRequest is the request envelope of QueryWencaiRealtimeFields.
type QueryWencaiRealtimeFieldsRequest struct {
	Securities []Security `json:"securities"`
	Fields     any        `json:"fields"`
}

// QueryWencaiRealtimeFields calls the native query_wencai_realtime_fields method (category wencai).
// Required fields: securities, fields.
// Returns WencaiRealtimeFieldsResult; the row set is under `rows`.
func (c *Client) QueryWencaiRealtimeFields(ctx context.Context, req QueryWencaiRealtimeFieldsRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by query_wencai_realtime_fields")
	}
	if req.Fields == nil {
		return RawData{}, invalidArgument("fields", "is required by query_wencai_realtime_fields")
	}
	return c.request(ctx, "query_wencai_realtime_fields", &req)
}

// QueryWencaiSecuritiesRequest is the request envelope of QueryWencaiSecurities.
type QueryWencaiSecuritiesRequest struct {
	Query   string   `json:"query"`
	Markets []string `json:"markets"`
	Limit   Opt      `json:"limit,omitempty"`
}

// QueryWencaiSecurities calls the native query_wencai_securities method (category wencai).
// Required fields: query.
// Returns WencaiSecuritiesResult; the row set is under `items`.
func (c *Client) QueryWencaiSecurities(ctx context.Context, req QueryWencaiSecuritiesRequest) (RawData, error) {
	if req.Query == "" {
		return RawData{}, invalidArgument("query", "is required by query_wencai_securities")
	}
	return c.request(ctx, "query_wencai_securities", &req)
}

// RankBlockSecuritiesRequest is the request envelope of RankBlockSecurities.
type RankBlockSecuritiesRequest struct {
	Block             Block      `json:"block"`
	SortBegin         Opt        `json:"sort_begin,omitempty"`
	SortCount         Opt        `json:"sort_count,omitempty"`
	SortID            Opt        `json:"sort_id,omitempty"`
	SortOrder         Opt        `json:"sort_order,omitempty"`
	ExcludeSecurities []Security `json:"exclude_securities"`
	ExcludeBlocks     []Block    `json:"exclude_blocks"`
	FuncPeriod        Opt        `json:"func_period,omitempty"`
}

// RankBlockSecurities calls the native rank_block_securities method (category rankings).
// Required fields: block.
// Returns SortedSecuritiesResult; the row set is under `securities`.
func (c *Client) RankBlockSecurities(ctx context.Context, req RankBlockSecuritiesRequest) (RawData, error) {
	if req.Block.ID == nil && req.Block.Name == "" {
		return RawData{}, invalidArgument("block", "is required by rank_block_securities")
	}
	return c.request(ctx, "rank_block_securities", &req)
}

// RankBlockSecuritiesByIndustryRequest is the request envelope of RankBlockSecuritiesByIndustry.
type RankBlockSecuritiesByIndustryRequest struct {
	Block     Block `json:"block"`
	SortBegin Opt   `json:"sort_begin,omitempty"`
	SortCount Opt   `json:"sort_count,omitempty"`
	SortOrder Opt   `json:"sort_order,omitempty"`
}

// RankBlockSecuritiesByIndustry calls the native rank_block_securities_by_industry method (category rankings).
// Required fields: block.
// Returns SortedSecuritiesResult; the row set is under `securities`.
func (c *Client) RankBlockSecuritiesByIndustry(ctx context.Context, req RankBlockSecuritiesByIndustryRequest) (RawData, error) {
	if req.Block.ID == nil && req.Block.Name == "" {
		return RawData{}, invalidArgument("block", "is required by rank_block_securities_by_industry")
	}
	return c.request(ctx, "rank_block_securities_by_industry", &req)
}

// RankBlockSecuritiesByWencaiFieldRequest is the request envelope of RankBlockSecuritiesByWencaiField.
type RankBlockSecuritiesByWencaiFieldRequest struct {
	Block     Block `json:"block"`
	Field     any   `json:"field"`
	Offset    Opt   `json:"offset,omitempty"`
	Limit     Opt   `json:"limit,omitempty"`
	SortOrder Opt   `json:"sort_order,omitempty"`
}

// RankBlockSecuritiesByWencaiField calls the native rank_block_securities_by_wencai_field method (category wencai).
// Required fields: block, field.
// Returns WencaiBlockRankingResult; the row set is under `items`.
func (c *Client) RankBlockSecuritiesByWencaiField(ctx context.Context, req RankBlockSecuritiesByWencaiFieldRequest) (RawData, error) {
	if req.Block.ID == nil && req.Block.Name == "" {
		return RawData{}, invalidArgument("block", "is required by rank_block_securities_by_wencai_field")
	}
	if req.Field == nil {
		return RawData{}, invalidArgument("field", "is required by rank_block_securities_by_wencai_field")
	}
	return c.request(ctx, "rank_block_securities_by_wencai_field", &req)
}

// RankConstituentsByPerformanceContributionRequest is the request envelope of RankConstituentsByPerformanceContribution.
type RankConstituentsByPerformanceContributionRequest struct {
	Code       string `json:"code"`
	Market     string `json:"market"`
	Target     Opt    `json:"target,omitempty"`
	SortID     Opt    `json:"sort_id,omitempty"`
	SortOrder  Opt    `json:"sort_order,omitempty"`
	SortBegin  Opt    `json:"sort_begin,omitempty"`
	SortCount  Opt    `json:"sort_count,omitempty"`
	ValidBegin Opt    `json:"valid_begin,omitempty"`
	ValidEnd   Opt    `json:"valid_end,omitempty"`
}

// RankConstituentsByPerformanceContribution calls the native rank_constituents_by_performance_contribution method (category rankings).
// Required fields: code, market.
// Returns a JSON array of PerformanceContributionItem objects.
func (c *Client) RankConstituentsByPerformanceContribution(ctx context.Context, req RankConstituentsByPerformanceContributionRequest) (RawData, error) {
	if req.Code == "" {
		return RawData{}, invalidArgument("code", "is required by rank_constituents_by_performance_contribution")
	}
	if req.Market == "" {
		return RawData{}, invalidArgument("market", "is required by rank_constituents_by_performance_contribution")
	}
	return c.request(ctx, "rank_constituents_by_performance_contribution", &req)
}

// RankRelatedSecuritiesRequest is the request envelope of RankRelatedSecurities.
type RankRelatedSecuritiesRequest struct {
	Security   Security `json:"security"`
	SortBegin  Opt      `json:"sort_begin,omitempty"`
	SortCount  Opt      `json:"sort_count,omitempty"`
	SortID     Opt      `json:"sort_id,omitempty"`
	SortOrder  Opt      `json:"sort_order,omitempty"`
	FuncPeriod Opt      `json:"func_period,omitempty"`
}

// RankRelatedSecurities calls the native rank_related_securities method (category rankings).
// Required fields: security.
// Returns SortedSecuritiesResult; the row set is under `securities`.
func (c *Client) RankRelatedSecurities(ctx context.Context, req RankRelatedSecuritiesRequest) (RawData, error) {
	if req.Security.Market == "" || req.Security.Code == "" {
		return RawData{}, invalidArgument("security", "is required by rank_related_securities")
	}
	return c.request(ctx, "rank_related_securities", &req)
}

// RankSecuritiesByPopularityRequest is the request envelope of RankSecuritiesByPopularity.
type RankSecuritiesByPopularityRequest struct {
	Type      Opt `json:"type,omitempty"`
	Benchmark any `json:"benchmark"`
}

// RankSecuritiesByPopularity calls the native rank_securities_by_popularity method (category rankings).
// Returns PopularityRankTopResult; the row set is under `items`.
func (c *Client) RankSecuritiesByPopularity(ctx context.Context, req RankSecuritiesByPopularityRequest) (RawData, error) {
	return c.request(ctx, "rank_securities_by_popularity", &req)
}

// RankSecuritiesByRealtimeStatisticRequest is the request envelope of RankSecuritiesByRealtimeStatistic.
type RankSecuritiesByRealtimeStatisticRequest struct {
	Markets           []string   `json:"markets"`
	Securities        []Security `json:"securities"`
	SortBy            int        `json:"sort_by"`
	SortDir           Opt        `json:"sort_dir,omitempty"`
	SortBegin         Opt        `json:"sort_begin,omitempty"`
	SortCount         Opt        `json:"sort_count,omitempty"`
	ExcludeSecurities []Security `json:"exclude_securities"`
}

// RankSecuritiesByRealtimeStatistic calls the native rank_securities_by_realtime_statistic method (category rankings).
// Required fields: sort_by.
// Returns RealtimeStatsResult; the row set is under `rows`.
func (c *Client) RankSecuritiesByRealtimeStatistic(ctx context.Context, req RankSecuritiesByRealtimeStatisticRequest) (RawData, error) {
	return c.request(ctx, "rank_securities_by_realtime_statistic", &req)
}

// RemoveAccountWatchlistGroupSecuritiesRequest is the request envelope of RemoveAccountWatchlistGroupSecurities.
type RemoveAccountWatchlistGroupSecuritiesRequest struct {
	GroupID    int        `json:"group_id"`
	Securities []Security `json:"securities"`
}

// RemoveAccountWatchlistGroupSecurities calls the native remove_account_watchlist_group_securities method (category account).
// Required fields: group_id, securities.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) RemoveAccountWatchlistGroupSecurities(ctx context.Context, req RemoveAccountWatchlistGroupSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by remove_account_watchlist_group_securities")
	}
	return c.request(ctx, "remove_account_watchlist_group_securities", &req)
}

// RemoveAccountWatchlistSecuritiesRequest is the request envelope of RemoveAccountWatchlistSecurities.
type RemoveAccountWatchlistSecuritiesRequest struct {
	Securities []Security `json:"securities"`
}

// RemoveAccountWatchlistSecurities calls the native remove_account_watchlist_securities method (category account).
// Required fields: securities.
// Mutates account state and returns AccountWatchlistUpdate.
func (c *Client) RemoveAccountWatchlistSecurities(ctx context.Context, req RemoveAccountWatchlistSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by remove_account_watchlist_securities")
	}
	return c.request(ctx, "remove_account_watchlist_securities", &req)
}

// RenameAccountWatchlistGroupRequest is the request envelope of RenameAccountWatchlistGroup.
type RenameAccountWatchlistGroupRequest struct {
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
}

// RenameAccountWatchlistGroup calls the native rename_account_watchlist_group method (category account).
// Required fields: group_id, name.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) RenameAccountWatchlistGroup(ctx context.Context, req RenameAccountWatchlistGroupRequest) (RawData, error) {
	if req.Name == "" {
		return RawData{}, invalidArgument("name", "is required by rename_account_watchlist_group")
	}
	return c.request(ctx, "rename_account_watchlist_group", &req)
}

// ReplaceAccountWatchlistGroupSecuritiesRequest is the request envelope of ReplaceAccountWatchlistGroupSecurities.
type ReplaceAccountWatchlistGroupSecuritiesRequest struct {
	GroupID    int        `json:"group_id"`
	Securities []Security `json:"securities"`
}

// ReplaceAccountWatchlistGroupSecurities calls the native replace_account_watchlist_group_securities method (category account).
// Required fields: group_id, securities.
// Mutates account state and returns AccountWatchlistGroups.
func (c *Client) ReplaceAccountWatchlistGroupSecurities(ctx context.Context, req ReplaceAccountWatchlistGroupSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by replace_account_watchlist_group_securities")
	}
	return c.request(ctx, "replace_account_watchlist_group_securities", &req)
}

// ReplaceAccountWatchlistSecuritiesRequest is the request envelope of ReplaceAccountWatchlistSecurities.
type ReplaceAccountWatchlistSecuritiesRequest struct {
	Version    int        `json:"version"`
	Securities []Security `json:"securities"`
}

// ReplaceAccountWatchlistSecurities calls the native replace_account_watchlist_securities method (category account).
// Required fields: version, securities.
// Mutates account state and returns AccountWatchlistUpdate.
func (c *Client) ReplaceAccountWatchlistSecurities(ctx context.Context, req ReplaceAccountWatchlistSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by replace_account_watchlist_securities")
	}
	return c.request(ctx, "replace_account_watchlist_securities", &req)
}

// ResolveBlockRequest is the request envelope of ResolveBlock.
type ResolveBlockRequest struct {
	Name string `json:"name"`
}

// ResolveBlock calls the native resolve_block method (category securities).
// Required fields: name.
// Returns Block.
func (c *Client) ResolveBlock(ctx context.Context, req ResolveBlockRequest) (RawData, error) {
	if req.Name == "" {
		return RawData{}, invalidArgument("name", "is required by resolve_block")
	}
	return c.request(ctx, "resolve_block", &req)
}

// ResolveSecuritiesRequest is the request envelope of ResolveSecurities.
type ResolveSecuritiesRequest struct {
	Codes []string `json:"codes"`
}

// ResolveSecurities calls the native resolve_securities method (category securities).
// Required fields: codes.
// Returns a JSON array of ResolvedSecurity objects.
func (c *Client) ResolveSecurities(ctx context.Context, req ResolveSecuritiesRequest) (RawData, error) {
	if len(req.Codes) == 0 {
		return RawData{}, invalidArgument("codes", "is required by resolve_securities")
	}
	return c.request(ctx, "resolve_securities", &req)
}

// ResolveWencaiSecuritiesRequest is the request envelope of ResolveWencaiSecurities.
type ResolveWencaiSecuritiesRequest struct {
	Query       string `json:"query"`
	Limit       Opt    `json:"limit,omitempty"`
	StockSuffix Opt    `json:"stock_suffix,omitempty"`
}

// ResolveWencaiSecurities calls the native resolve_wencai_securities method (category wencai).
// Required fields: query.
// Returns ResolvedWencaiSecuritiesResult; the row set is under `items`.
func (c *Client) ResolveWencaiSecurities(ctx context.Context, req ResolveWencaiSecuritiesRequest) (RawData, error) {
	if req.Query == "" {
		return RawData{}, invalidArgument("query", "is required by resolve_wencai_securities")
	}
	return c.request(ctx, "resolve_wencai_securities", &req)
}

// SearchSecuritiesRequest is the request envelope of SearchSecurities.
type SearchSecuritiesRequest struct {
	Pattern string `json:"pattern"`
	Market  Opt    `json:"market,omitempty"`
}

// SearchSecurities calls the native search_securities method (category securities).
// Required fields: pattern.
// Returns a JSON array of SecurityCandidate objects.
func (c *Client) SearchSecurities(ctx context.Context, req SearchSecuritiesRequest) (RawData, error) {
	if req.Pattern == "" {
		return RawData{}, invalidArgument("pattern", "is required by search_securities")
	}
	return c.request(ctx, "search_securities", &req)
}

// SortSecuritiesRequest is the request envelope of SortSecurities.
type SortSecuritiesRequest struct {
	Securities []Security `json:"securities"`
	SortBegin  Opt        `json:"sort_begin,omitempty"`
	SortCount  Opt        `json:"sort_count,omitempty"`
	SortID     Opt        `json:"sort_id,omitempty"`
	SortOrder  Opt        `json:"sort_order,omitempty"`
	FuncPeriod Opt        `json:"func_period,omitempty"`
}

// SortSecurities calls the native sort_securities method (category rankings).
// Required fields: securities.
// Returns SortedSecuritiesResult; the row set is under `securities`.
func (c *Client) SortSecurities(ctx context.Context, req SortSecuritiesRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by sort_securities")
	}
	return c.request(ctx, "sort_securities", &req)
}

// SortSecuritiesByIndustryRequest is the request envelope of SortSecuritiesByIndustry.
type SortSecuritiesByIndustryRequest struct {
	Securities []Security `json:"securities"`
	SortBegin  Opt        `json:"sort_begin,omitempty"`
	SortCount  Opt        `json:"sort_count,omitempty"`
	SortOrder  Opt        `json:"sort_order,omitempty"`
}

// SortSecuritiesByIndustry calls the native sort_securities_by_industry method (category rankings).
// Required fields: securities.
// Returns SortedSecuritiesResult; the row set is under `securities`.
func (c *Client) SortSecuritiesByIndustry(ctx context.Context, req SortSecuritiesByIndustryRequest) (RawData, error) {
	if len(req.Securities) == 0 {
		return RawData{}, invalidArgument("securities", "is required by sort_securities_by_industry")
	}
	return c.request(ctx, "sort_securities_by_industry", &req)
}
