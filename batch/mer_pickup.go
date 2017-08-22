package main

import (
	"../include/option"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
	"container/list"
	"bufio"
	"strings"
	"strconv"
)

type Buy struct {
	Date string
	ID string
	Price int
	Title string
}

func main() {
	o := option.Load()
	dataDir := o.Get("data_dir")

	infos, err := ioutil.ReadDir(dataDir)
	if err != nil {
		os.Exit(1)
	}
	table := map[string]*Buy{}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		path := filepath.Join(dataDir, info.Name())
		fmt.Println("path:" + path)
		list := load(path)
		for i := list.Front(); i != nil; i = i.Next() {
			buy := i.Value.(*Buy)
			if _, ok := table[buy.ID]; !ok {
				table[buy.ID] = buy
			}
		}
	}

	dates := map[string]map[string]*Buy{}
	for _, v := range table {
		if _, ok := dates[v.Date]; !ok {
			dates[v.Date] = map[string]*Buy{}
		}
		dates[v.Date][v.ID] = v;
	}
	for k, m := range dates {
		fmt.Println("date:" + k)
		for _, b := range m {
			fmt.Println("\t" + b.ID, b.Price, b.Title)
		}
	}
}

func load(path string) *list.List {
	list := new(list.List)
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return list
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cells := strings.Split(scanner.Text(), ",")
		if len(cells) < 4 {
			continue
		}
		b := new(Buy)
		b.Date = cells[0]
		b.ID = cells[1]
		b.Price, _ = strconv.Atoi(cells[2])
		b.Title = cells[3]
		list.PushBack(b)
	}
	return list
}