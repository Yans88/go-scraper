package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgUrl string `json:"imgurl"`
}

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)
	var items []item
	c.OnHTML("div.col-sm-9 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
		item := item{
			Name:   h.ChildText("h2.product-title"),
			Price:  h.ChildText("div.sale-price"),
			ImgUrl: h.ChildAttr("img", "src"),
		}
		items = append(items, item)
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		next_page := h.Request.AbsoluteURL(h.Attr("href"))
		c.Visit(next_page)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	c.Visit("https://j2store.net/demo/index.php/shop")
	content, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("products.json", content, 0644)
}

// Patreon: https://www.patreon.com/johnwatsonrooney (NEW)
// Oxylabs: https://oxylabs.go2cloud.org/aff_c?of... - code JR15
// Amazon UK: https://amzn.to/2OYuMwo
// Hosting: Digital Ocean: https://m.do.co/c/c7c90f161ff6
// Gear Used: https://jhnwr.com/gear/ (NEW)
