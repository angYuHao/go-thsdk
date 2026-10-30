package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
)

func main() {
	client, err := thsdk.New(thsdk.Options{Verbose: true})
	must(err)
	defer client.Close()

	must(client.Auth(context.Background(), thsdk.Credentials{}))

	ctx := context.Background()

	mout, err := json.Marshal(thsdk.MustSecurity("USHA600519"))
	must(err)
	fmt.Println("security json:", string(mout))

	// Typed standard method: KLine rows.
	klineData, err := client.ListSecurityKLines(ctx, thsdk.ListSecurityKLinesRequest{
		Security: thsdk.MustSecurity("USHA600519"),
		Start:    -10,
		End:      0,
		Adjust:   thsdk.AdjustForward,
		Period:   thsdk.PeriodDay,
	})
	must(err)
	var klines []thsdk.KLine
	must(klineData.Rows(&klines))
	fmt.Printf("klines: %d rows, last=%+v\n", len(klines), klines[len(klines)-1])
	date, err := klines[len(klines)-1].TradeDate()
	must(err)
	fmt.Println("last trade date:", date.Format("2006-01-02"))

	// Order book with typed decode.
	bookData, err := client.ListSecurityOrderBooks(ctx, thsdk.ListSecurityOrderBooksRequest{
		Securities: []thsdk.Security{thsdk.MustSecurity("USHA600519")},
		Depth:      thsdk.DepthFive,
	})
	if err != nil {
		// The live snapshot is not always served; the upstream Python SDK fails
		// the same way, so keep going instead of aborting.
		fmt.Println("order book unavailable right now:", err)
	} else {
		var books []thsdk.OrderBook
		must(bookData.Rows(&books))
		fmt.Printf("order book best bid/ask: %v / %v\n", books[0].Bids[0], books[0].Asks[0])
	}

	// Envelope result with metadata.
	flash, err := client.ListNewsFlashItems(ctx, thsdk.ListNewsFlashItemsRequest{
		Traversal: thsdk.FlashPage, Page: thsdk.Num(1), PageSize: thsdk.Num(3),
	})
	must(err)
	var items []map[string]any
	must(flash.Rows(&items))
	meta, err := flash.Metadata()
	must(err)
	fmt.Printf("flash: %d items, metadata keys=%d\n", len(items), len(meta))

	// Search and resolve.
	found, err := client.SearchSecurities(ctx, thsdk.SearchSecuritiesRequest{Pattern: "茅台"})
	must(err)
	var candidates []thsdk.SecurityRef
	must(found.Rows(&candidates))
	fmt.Printf("search: %s -> %s\n", candidates[0].Name, candidates[0].FullCode)

	// Context deadline becomes the component-side timeout.
	short, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()
	if _, err := client.ListMarketSecurities(short, thsdk.ListMarketSecuritiesRequest{
		Market: "USHA", SortBegin: thsdk.Num(0), SortCount: thsdk.Num(5),
	}); err != nil {
		fmt.Println("deadline path reported:", err)
	} else {
		fmt.Println("deadline path unexpectedly succeeded")
	}

	// Rate limiting: five back-to-back calls with no manual pause.
	for i := 0; i < 5; i++ {
		if _, err := client.AccountPermissions(ctx); err != nil {
			must(err)
		}
	}
	fmt.Println("five unpaced calls survived the internal 50ms spacing")
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
