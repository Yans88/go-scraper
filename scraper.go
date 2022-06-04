package main

import (
	"bytes"
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

	content, err := JSONMarshal(items, true)
	if err != nil {
		fmt.Println(err.Error())
	}
	os.WriteFile("products.json", content, 0644)
}

func JSONMarshal(v interface{}, safeEncoding bool) ([]byte, error) {
    b, err := json.Marshal(v)

    if safeEncoding {
        b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
        b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
        b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
    }
    return b, err
}