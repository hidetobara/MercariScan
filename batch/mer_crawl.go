package batch

import (
	"fmt"
	"net/http"
	"os"
	"gopkg.in/xmlpath.v2"
	"strings"
	"strconv"
)

func main() {
	response, err := http.Get("https://www.mercari.com/jp/search/?sort_order=&keyword=switch&price_min=30000&price_max=60000&status_all=1&status_trading_sold_out=1");
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("status:", response.Status)

	//fmt.Println(string(body))
	root, err := xmlpath.ParseHTML(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	xSection := xmlpath.MustCompile(`//section[@class='items-box']`)
	//path := xmlpath.MustCompile(`//div[@class='items-box-price font-5']`)
	xHref := xmlpath.MustCompile( `./a[@href]/@href`)
	xPrice := xmlpath.MustCompile( `.//div[@class='items-box-price font-5']`)
	xName := xmlpath.MustCompile( `.//h3[@class='items-box-name font-2']`)

	iSection := xSection.Iter(root)
	for iSection.Next() { // イテレータ回せ
		section := iSection.Node()

		iHref := xHref.Iter(section)
		iHref.Next()
		href := iHref.Node()
		cells := strings.Split(href.String(), "/")
		id := cells[len(cells) - 2]

		iPrice := xPrice.Iter(section)
		iPrice.Next()
		price := iPrice.Node()
		tmp := strings.Trim(price.String(), "¥ ")
		tmp = strings.Replace(tmp, ",", "", -1)
		yen, err := strconv.Atoi(tmp)
		if err != nil {
			continue
		}

		iName := xName.Iter(section)
		iName.Next()
		name := iName.Node().String()

		fmt.Println(id, yen, name)
	}
}
