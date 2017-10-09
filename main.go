package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"regexp"
	"time"
	// "github.com/gocarina/gocsv"
)

func main() {
	tickerPtr := flag.String("ticker", "", "Ticker to pull prices for")
	flag.Parse()

	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	resp, err := client.GetSecurityData(*tickerPtr, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC), time.Date(2017, time.March, 30, 0, 0, 0, 0, time.UTC))
	fmt.Println(resp)
}

func NewClient() (*Client, error) {
	var seedTickers = []string{"AAPL", "GOOG", "MSFT"}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	httpClient := &http.Client{Jar: jar}
	c := Client{HttpClient: httpClient}

	i := rand.Intn(len(seedTickers))
	ticker := seedTickers[i]
	crumb, err := getCrumb(c.HttpClient, ticker)
	if err != nil {
		return nil, err
	}

	c.Crumb = crumb

	return &c, nil

}

type Client struct {
	Crumb      string
	HttpClient *http.Client
}

type Price struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func (c *Client) GetSecurityData(ticker string, startDate, endDate time.Time) (string, error) {
	urlFmtStr := "https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%d&period2=%d&interval=1d&events=history&crumb=%s"
	url := fmt.Sprintf(urlFmtStr, ticker, startDate.Unix(), endDate.Unix(), c.Crumb)
	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

func getCrumb(client *http.Client, ticker string) (string, error) {
	url := fmt.Sprintf("https://finance.yahoo.com/quote/%s", ticker)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	crumbRe, err := regexp.Compile(`"CrumbStore":{"crumb":"([^"]+)"\}`)
	if err != nil {
		return "", err
	}

	matches := crumbRe.FindAllStringSubmatch(string(body), 1)

	if len(matches) > 0 {
		return matches[0][1], nil
	}

	return "", nil

}
