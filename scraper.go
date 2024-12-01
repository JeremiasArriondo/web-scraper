package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Product struct {
	Url, Image, Name, Price string
}

func main() {
	var products []Product

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	c.Visit("https://www.scrapingcourse.com/ecommerce")

	c.OnHTML(
		"li.product", func(e *colly.HTMLElement) {
			product := Product{}
			// scrape the target data
			product.Url = e.ChildAttr("a", "href")
			product.Image = e.ChildAttr("img", "src")
			product.Name = e.ChildText(".product-name")
			product.Price = e.ChildText(".price")

			// add the product instance with scraped data to the list of products
			products = append(products, product)
		},
	)

	c.OnScraped(func(r *colly.Response) {
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file: ", err)
		}
		defer file.Close()

		// initialize a file writer
		writer := csv.NewWriter(file)

		// write headers to the CSV file
		headers := []string{
			"Url",
			"Image",
			"Name",
			"Price",
		}
		writer.Write(headers)

		// write each product as a CSV row
		for _, product := range products {
			// convert a Product to an array of strings
			record := []string{
				product.Url,
				product.Image,
				product.Name,
				product.Price,
			}

			// add a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	},
	)

}
