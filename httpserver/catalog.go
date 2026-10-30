package httpserver

// methodInfo describes one exposed runtime method.
type methodInfo struct {
	Category string
	Write    bool
}

// methodCatalog lists every method this server exposes: the 70 read-only calls
// plus the 10 account mutations, which stay gated behind Config.AllowWrites.
func methodCatalog() map[string]methodInfo {
	catalog := map[string]methodInfo{
		"account_permissions":                             {Category: "account"},
		"analyze_security_limit_up":                       {Category: "news"},
		"calculate_security_realtime_statistics":          {Category: "rankings"},
		"check_market_timeline_readiness":                 {Category: "market_data"},
		"get_account_watchlist":                           {Category: "account"},
		"get_account_watchlist_groups":                    {Category: "account"},
		"get_market_metadata":                             {Category: "securities"},
		"get_market_security_names":                       {Category: "securities"},
		"get_security_concept_tags":                       {Category: "securities"},
		"get_security_industry":                           {Category: "securities"},
		"get_security_popularity_rank":                    {Category: "rankings"},
		"get_security_price_limit_events":                 {Category: "news"},
		"get_security_short_term_highlights":              {Category: "news"},
		"get_security_trading_timeline":                   {Category: "market_data"},
		"list_block_constituents":                         {Category: "securities"},
		"list_block_descriptions":                         {Category: "securities"},
		"list_commodity_stock_linkage_news":               {Category: "news"},
		"list_futures_related_securities":                 {Category: "securities"},
		"list_hot_block_calendar_events":                  {Category: "news"},
		"list_hot_event_news":                             {Category: "news"},
		"list_index_call_auction_quotes":                  {Category: "market_data"},
		"list_industry_children":                          {Category: "securities"},
		"list_market_call_auction_unusual_events":         {Category: "news"},
		"list_market_change_events":                       {Category: "news"},
		"list_market_securities":                          {Category: "securities"},
		"list_market_short_term_events":                   {Category: "news"},
		"list_news_flash_items":                           {Category: "news"},
		"list_news_items":                                 {Category: "news"},
		"list_related_security_performances":              {Category: "rankings"},
		"list_securities_by_realtime_signal":              {Category: "rankings"},
		"list_security_ah_relations":                      {Category: "securities"},
		"list_security_block_memberships":                 {Category: "securities"},
		"list_security_call_auction_quotes":               {Category: "market_data"},
		"list_security_corporate_actions":                 {Category: "market_data"},
		"list_security_daily_capital_flows":               {Category: "market_data"},
		"list_security_daily_k_lines_with_previous_close": {Category: "market_data"},
		"list_security_extended_hours_intraday_bars":      {Category: "market_data"},
		"list_security_extended_hours_ticks":              {Category: "market_data"},
		"list_security_financial_snapshots":               {Category: "market_data"},
		"list_security_futures_relations":                 {Category: "securities"},
		"list_security_industry_mappings":                 {Category: "securities"},
		"list_security_intraday_bars":                     {Category: "market_data"},
		"list_security_k_lines":                           {Category: "market_data"},
		"list_security_link_relations":                    {Category: "securities"},
		"list_security_news_events":                       {Category: "news"},
		"list_security_news_markers":                      {Category: "news"},
		"list_security_option_greeks":                     {Category: "market_data"},
		"list_security_order_books":                       {Category: "market_data"},
		"list_security_price_volume_levels":               {Category: "market_data"},
		"list_security_research_reports":                  {Category: "news"},
		"list_security_ticks":                             {Category: "market_data"},
		"list_security_time_and_sales_ticks":              {Category: "market_data"},
		"list_wencai_expression_fields":                   {Category: "wencai"},
		"list_wencai_hot_blocks":                          {Category: "wencai"},
		"query_wencai":                                    {Category: "wencai"},
		"query_wencai_realtime_fields":                    {Category: "wencai"},
		"query_wencai_securities":                         {Category: "wencai"},
		"rank_block_securities":                           {Category: "rankings"},
		"rank_block_securities_by_industry":               {Category: "rankings"},
		"rank_block_securities_by_wencai_field":           {Category: "rankings"},
		"rank_constituents_by_performance_contribution":   {Category: "rankings"},
		"rank_related_securities":                         {Category: "rankings"},
		"rank_securities_by_popularity":                   {Category: "rankings"},
		"rank_securities_by_realtime_statistic":           {Category: "rankings"},
		"resolve_block":                                   {Category: "securities"},
		"resolve_securities":                              {Category: "securities"},
		"resolve_wencai_securities":                       {Category: "wencai"},
		"search_securities":                               {Category: "securities"},
		"sort_securities":                                 {Category: "securities"},
		"sort_securities_by_industry":                     {Category: "securities"},
	}
	for name, info := range writeCatalog() {
		catalog[name] = info
	}
	return catalog
}

// writeCatalog holds the account mutations, which would report Write: true.
func writeCatalog() map[string]methodInfo {
	return map[string]methodInfo{
		"add_account_watchlist_group_securities":     {Category: "account", Write: true},
		"add_account_watchlist_securities":           {Category: "account", Write: true},
		"clear_account_watchlist":                    {Category: "account", Write: true},
		"create_account_watchlist_group":             {Category: "account", Write: true},
		"delete_account_watchlist_group":             {Category: "account", Write: true},
		"remove_account_watchlist_group_securities":  {Category: "account", Write: true},
		"remove_account_watchlist_securities":        {Category: "account", Write: true},
		"rename_account_watchlist_group":             {Category: "account", Write: true},
		"replace_account_watchlist_group_securities": {Category: "account", Write: true},
		"replace_account_watchlist_securities":       {Category: "account", Write: true},
	}
}
