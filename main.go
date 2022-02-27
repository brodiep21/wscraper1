package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func main() {

	file, err := os.Create("costs.csv")
	if err != nil {
		log.Fatalf("Could not create file, error is %q", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("www.ebay.com"),
	)

	//finds the price of each individual card, parses the value to be an int in the format of up to 4 digits with no decimal and prints it to the terminal
	c.OnHTML(".s-item__detail.s-item__detail--primary", func(e *colly.HTMLElement) {
		prices := e.ChildText("span.s-item__price")
		prices = strings.ReplaceAll(prices, "$", "")
		prices = strings.ReplaceAll(prices, ",", "")
		newprice := ""
		for k, v := range prices {
			if k < 4 && string(v) != "." {
				newprice += string(v)
			}
		}
		csvCosts := make([]string, 0)
		//filters out the ghost text behind the HTML element indicating $0
		if newprice == "" {
		} else {
			//formats string to int and compares price, then converts back to string in order to create a slice of strings
			parsePrice, err := strconv.Atoi(newprice)
			if err == nil && parsePrice <= 1500 {
				parsePrice2 := strconv.Itoa(parsePrice)
				csvCosts = append(csvCosts, parsePrice2)
				writer.Write(csvCosts)
			}
		}
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String()+" and getting video card prices")
	})

	c.Visit("https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1312&_nkw=gtx+3080+graphics+card&_sacat=0")
}
