package main

import (
	"fmt"
	"net/http"
	"os"
	"gopkg.in/xmlpath.v2"
	"strings"
	"strconv"
	"time"
	"container/list"
	"io"
	"bufio"
	"../include/option"
)


type Buy struct {
	Date string
	ID string
	Price int
	Title string
}

func (b *Buy) ToCsv() string {
	return fmt.Sprintf("%s,%s,%d,%s", b.Date, b.ID, b.Price, b.Title)
}

func downloadPage(url string) io.Reader {
	response, err := http.Get(url);
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("status:", response.Status)

	return response.Body
}

func retrieveSales(pager io.Reader) *list.List {
	root, err := xmlpath.ParseHTML(pager)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	xSection := xmlpath.MustCompile(`//section[@class='items-box']`)
	//path := xmlpath.MustCompile(`//div[@class='items-box-price font-5']`)
	xHref := xmlpath.MustCompile( `./a[@href]/@href`)
	xPrice := xmlpath.MustCompile( `.//div[@class='items-box-price font-5']`)
	xName := xmlpath.MustCompile( `.//h3[@class='items-box-name font-2']`)

	now := time.Now()
	//nowStr := fmt.Sprintf("%04d-%02d-%02d", now.Year(), now.Month(), now.Day())
	nowStr := now.Format("2006-01-02 15:04:05")

	sales := list.New()
	iSection := xSection.Iter(root)
	for iSection.Next() {
		b := new(Buy)
		b.Date = nowStr
		section := iSection.Node()

		iHref := xHref.Iter(section)
		iHref.Next()
		href := iHref.Node()
		cells := strings.Split(href.String(), "/")
		b.ID = cells[len(cells) - 2]

		iPrice := xPrice.Iter(section)
		iPrice.Next()
		price := iPrice.Node()
		tmp := strings.Trim(price.String(), "¥ ")
		tmp = strings.Replace(tmp, ",", "", -1)
		b.Price, err = strconv.Atoi(tmp)
		if err != nil {
			continue
		}

		iName := xName.Iter(section)
		iName.Next()
		b.Title= iName.Node().String()

		sales.PushBack(b)
	}
	return sales
}

func main() {
	o := option.Load()
	pager := downloadPage(o.Get("url"))
	sales := retrieveSales(pager)

	now := time.Now()
	filename := fmt.Sprintf("%04d%02d%02d_%02d.csv", now.Year(), now.Month(), now.Day(), now.Hour())
	file, err := os.Create(o.Get("data_dir") + filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	writer := bufio.NewWriter(file)
	for i := sales.Front(); i != nil; i = i.Next() {
		writer.WriteString(i.Value.(*Buy).ToCsv() + "\n")
	}
	writer.Flush()
}
