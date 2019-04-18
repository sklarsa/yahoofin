package yahoofin

import (
	"encoding/csv"
	"strings"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Error(err)
	}

	ticker := "AAPL"
	startDate := time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2016, time.February, 1, 0, 0, 0, 0, time.UTC)

	fields := []Field{History, Dividend}

	for _, field := range fields {

		csvString, err := client.GetSecurityDataString(ticker, startDate, endDate, field)
		if err != nil {
			t.Error(err)
		}

		rdr := strings.NewReader(csvString)
		csvRdr := csv.NewReader(rdr)

		records, err := csvRdr.ReadAll()
		if err != nil {
			t.Error(err)
		}
		if len(records) == 0 {
			t.Error("No records parsed")
		}

		prices, err := client.GetSecurityData(ticker, startDate, endDate, History)
		if err != nil {
			t.Error(err)
		}

		if len(prices) == 0 {
			t.Error("No records parsed")
		}

	}

}
