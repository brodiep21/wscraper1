package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	//file creater for costs
	file, err := os.Create("costs.csv")
	if err != nil {
		log.Fatalf("Could not create file, error is %q", err)
	}
	//file creater for links
	file2, err := os.Create("links.csv")
	if err != nil {
		log.Fatalf("Could not create file, error is %q", err)
	}
	//allows writer to write in the correlating file
	writer := csv.NewWriter(file)
	writer2 := csv.NewWriter(file2)

	defer file.Close()
	defer file2.Close()

	defer writer.Flush()
	defer writer2.Flush()

	c := colly.NewCollector(
		colly.AllowedDomains("www.ebay.com"),
	)

	//finds the link associated with each card, adds it to a slice of strings and prints it to a link file
	c.OnHTML(".s-item__info.clearfix", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		csvLinks := make([]string, 0)
		csvLinks = append(csvLinks, link)
		writer2.Write(csvLinks)
	})
	//finds the price of each individual card, parses the value to be an int in the format of up to 4 digits with no decimal and prints it to the terminal
	c.OnHTML(".s-item__detail.s-item__detail--primary", func(e *colly.HTMLElement) {
		prices := e.ChildText("span.s-item__price")
		csvCosts := make([]string, 0)
		prices = strings.ReplaceAll(prices, "$", "")
		prices = strings.ReplaceAll(prices, ",", "")
		newprice := ""
		for k, v := range prices {
			if k < 4 && string(v) != "." {
				newprice += string(v)
			}
		}
		//filters out the ghost text behind the HTML element indicating $0
		if newprice == "" {
		} else {
			//formats string to int and compares price, then converts back to string in order to create a slice of strings
			csvCosts = append(csvCosts, newprice)
			writer.Write(csvCosts)
		}
	})
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String()+" and getting video card prices")
	})

	c.Visit("https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1312&_nkw=gtx+3080+graphics+card&_sacat=0")
}
