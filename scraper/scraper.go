package scraper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func Scraper() {
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
		if newprice == "" {
		} else {
			parsePrice, err := strconv.Atoi(newprice)
			if err == nil {
				fmt.Println(parsePrice)
			}
		}
	})
	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Currently Visiting", request.URL.String())
	})
	c.Visit("https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2380057.m570.l1312&_nkw=gtx+3080+graphics+card&_sacat=0")
}
