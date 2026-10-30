package main

import (
	"context"
	"fmt"
	"os"
	"time"

	thsdk "github.com/angYuHao/go-thsdk"
)

func main() {
	client, err := thsdk.New(thsdk.Options{})
	if err != nil {
		fail(err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Auth(ctx, thsdk.Credentials{}); err != nil {
		fail(err)
	}

	daily, err := client.GetPrice(ctx, thsdk.PriceQuery{
		Security:  thsdk.MustSecurity("USHA600519"),
		Frequency: "daily",
		Start:     time.Now().AddDate(-1, 0, 0).In(thsdk.Shanghai),
		End:       time.Now().In(thsdk.Shanghai),
		Fields:    []string{"time", "open", "high", "low", "close", "volume"},
	})
	if err != nil {
		fail(err)
	}
	fmt.Printf("last year candles: %d\n", len(daily))

	wencai, err := client.Wencai(ctx, "市盈率小于20，总市值大于500亿", 4)
	if err != nil {
		fail(err)
	}
	fmt.Printf("wencai %q matched %d rows\n", wencai.Query, len(wencai.Rows))
	for index, row := range wencai.Rows {
		if index == 3 {
			break
		}
		fmt.Printf("  %v\n", row)
	}

	page, err := client.News(ctx, 1, 5, "")
	if err != nil {
		fail(err)
	}
	fmt.Printf("flash page: %d items, %d metadata keys\n", len(page.Items), len(page.Metadata))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
