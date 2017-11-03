package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
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

	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	resp, err := client.GetSecurityDataString(*tickerPtr, startDate, endDate)
	fmt.Println(resp)
}

// NewClient creates a new Yahoo Finance client
func NewClient() (*Client, error) {
	var seedTickers = []string{"AAPL", "GOOG", "MSFT"}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar}
	c := Client{httpClient: httpClient}

	i := rand.Intn(len(seedTickers))
	ticker := seedTickers[i]
	crumb, err := getCrumb(c.httpClient, ticker)
	if err != nil {
		return nil, err
	}

	c.crumb = crumb

	return &c, nil

}

// Client is a struct that represents a Yahoo Finance client
type Client struct {
	// crumb is sent along with each request and is needed to make successful requests directly to the historical prices endpoint
	crumb string
	// httpClient is a persistent client used to store cookies after the initial request is sent
	httpClient *http.Client
}

// Price represents a single datapoint returned by the yahoo api
type Price struct {
	Date     DateTime `csv:"Date"`
	Open     float64  `csv:"Open"`
	High     float64  `csv:"High"`
	Low      float64  `csv:"Low"`
	Close    float64  `csv:"Close"`
	AdjClose float64  `csv:"Adj Close"`
	Volume   float64  `csv:"Volume"`
}

func (c *Client) makeRequest(ticker string, startDate, endDate time.Time) (*http.Response, error) {
	urlFmtStr := "https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=1d&events=history&crumb=%s"
	url := fmt.Sprintf(urlFmtStr, ticker, startDate.Unix(), endDate.Unix(), c.crumb)
	return c.httpClient.Get(url)
}

// GetSecurityDataString returns the raw response data from the yahoo endpoint.
// This string will be CSV formatted if the request succeeds.
// In the event of a failed request, this string will be JSON-formatted
func (c *Client) GetSecurityDataString(ticker string, startDate, endDate time.Time) (string, error) {
	resp, err := c.makeRequest(ticker, startDate, endDate)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// GetSecurityData returns a slice of pointers to Price structs, based on the data received from yahoo
func (c *Client) GetSecurityData(ticker string, startDate, endDate time.Time) ([]*Price, error) {
	prices := []*Price{}
	resp, err := c.makeRequest(ticker, startDate, endDate)
	if err != nil {
		return prices, err
	}

	if err := gocsv.Unmarshal(resp.Body, &prices); err != nil {
		return prices, err
	}

	return prices, nil

}

// DateTime is a custom implementation of time.Time used to unmarshal yahoo csv data
type DateTime struct {
	time.Time
}

// UnmarshalCSV sonverts the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02", csv)
	if err != nil {
		return err
	}
	return nil
}

func getCrumb(client *http.Client, ticker string) (string, error) {
	crumb := ""
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", ticker)
	resp, err := client.Get(url)
	if err != nil {
		return crumb, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return crumb, err
	}

	crumbRe, err := regexp.Compile(`"CrumbStore":{"crumb":"([^"]+)"\}`)
	if err != nil {
		return crumb, err
	}

	matches := crumbRe.FindAllStringSubmatch(string(body), 1)

	if len(matches) > 0 {
		crumb = matches[0][1]
	}

	if len(crumb) > 0 {
		crumb, err = strconv.Unquote(`"` + crumb + `"`)
		if err != nil {
			return crumb, err
		}
	}

	return crumb, nil

}
