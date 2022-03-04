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
	//file creation
	file, err := os.Create("cards.csv")
	if err != nil {
		log.Fatalf("Could not create file, error is %q", err)
	}
	//allows writer to write in the correlating file
	writer := csv.NewWriter(file)

	defer file.Close()

	defer writer.Flush()

	linksNcosts := make([][]string, 0)

	c := colly.NewCollector(
		colly.AllowedDomains("www.ebay.com"),
	)

	//finds the link associated with each card, adds it to a slice of strings and prints it to a link file
	c.OnHTML(".s-item__info.clearfix", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")

		csvLinks := make([]string, 0)
		csvLinks = append(csvLinks, link)

		linksNcosts = append(linksNcosts, csvLinks)
	})

	count := 0
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

		//filters out the ghost text behind the HTML element indicating $0 for html console, as well as a ghost 200 item
		if newprice == "" || newprice == "200" {
		} else {
			//range over the [][]strings that already has the links, and append costs to the proper []strings
			for k, v := range linksNcosts {
				if k == count+1 {
					linksNcosts[count] = append(v, newprice)
				} else {
					continue
				}
			}
			count++
		}
	})

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String()+" and getting video card prices")
	})

	//original search page
	ogpage := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=gtx+3080+graphics+card&_sacat=0&rt=nc&LH_BIN=1"
	pagecounter := 1

	//loop through respective pages
	for i := 0; i < 17; i++ {
		if i == 0 {
			c.Visit(ogpage)
		} else {
			count++
			pagecounter++
			newpage := ogpage + "&_pgn=" + strconv.Itoa(pagecounter)
			c.Visit(newpage)
		}
	}

	//takes the values([]strings) and writes them directly to CSV
	for _, cLink := range linksNcosts {
		if err := writer.Write(cLink); err != nil {
			log.Fatalln("Failed printing to csv file")
		}
	}
}
