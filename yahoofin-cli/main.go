package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sklarsa/yahoofin"
)

func main() {
	tickerPtr := flag.String("ticker", "", "Ticker to pull prices for")
	startDatePtr := flag.String("startDate", "", "Start date [YYYY-MM-DD]")
	endDatePtr := flag.String("endDate", "", "End date [YYYY-MM-DD]")
	flag.Parse()

	if *tickerPtr == "" {
		fmt.Println("Must provide a ticker")
		flag.Usage()
		os.Exit(1)
	}

	const dateFmt = "2006-01-02"

	startDate, err := time.Parse(dateFmt, *startDatePtr)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	endDate, err := time.Parse(dateFmt, *endDatePtr)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	client, err := yahoofin.NewClient()
	if err != nil {
		panic(err)
	}
	resp, err := client.GetSecurityDataString(*tickerPtr, startDate, endDate, yahoofin.History)
	fmt.Println(resp)
}
