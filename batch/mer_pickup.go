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
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		path := filepath.Join(dataDir, info.Name())
		list := load(path)
		for i := list.Front(); i != nil; i = i.Next() {
			buy := i.Value.(*Buy)
			fmt.Println(buy.ID)
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