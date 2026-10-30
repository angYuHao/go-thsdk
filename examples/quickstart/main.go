package main

import (
	"context"
	"fmt"
	"os"

	thsdk "github.com/angYuHao/go-thsdk"
)

func main() {
	if err := thsdk.Auth(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer thsdk.Close()

	client, err := thsdk.Default()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	ctx := context.Background()

	// Convenience layer: recent daily KLines, forward adjusted.
	rows, err := client.GetPrice(ctx, thsdk.PriceQuery{
		Security:  thsdk.MustSecurity("USHA600519"),
		Count:     5,
		Frequency: "daily",
		FQ:        "pre",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, row := range rows {
		fmt.Printf("KLine %s: close=%s volume=%s\n", row["time"], row["close"], row["volume"])
	}

	// Standard layer: typed request, typed rows.
	// The live snapshot is not always served — the upstream Python SDK fails the
	// same way — so report it as unavailable instead of aborting.
	book, err := client.Depth(ctx, thsdk.MustSecurity("USHA600519"), thsdk.MustSecurity("USZA300033"))
	if err != nil {
		fmt.Println("depth unavailable right now:", err)
		return
	}
	fmt.Printf("depth rows: %d, best bid %v\n", len(book), book[0])
}
