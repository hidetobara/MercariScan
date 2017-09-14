package main

import (
	"../include/option"
	"io/ioutil"
	"os"
	"fmt"
	"path/filepath"
	"bufio"
	"strings"
	"strconv"
	"time"
	"sort"
	"math"
)

type Buy struct {
	Date time.Time
	ID string
	Price int
	Title string
}
type Average struct {
	Count int
	Amount int
	Distribution int
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
		for _, buy := range list {
			if _, ok := table[buy.ID]; !ok {
				table[buy.ID] = buy
			}
		}
	}

	keys := []string{}
	dates := map[string]map[string]*Buy{}
	averages := map[string]*Average{}
	for _, v := range table {
		key := fmt.Sprintf("%02d/%02d", v.Date.Month(), v.Date.Day())
		if _, ok := dates[key]; !ok {
			dates[key] = map[string]*Buy{}
			keys = append(keys, key)
		}
		dates[key][v.ID] = v

		if _, ok := averages[key]; !ok {
			averages[key] = new(Average)
		}
		averages[key].Amount += v.Price
		averages[key].Distribution += v.Price * v.Price
		averages[key].Count += 1
	}
	sort.Strings(keys)
/*
	fmt.Println("-----dates")
	for _, key := range keys {
		for _, b := range dates[key] {
			fmt.Printf("%s,%s,%d,%s\n", key, b.ID, b.Price, b.Title )
		}
	}
*/
	fmt.Println("-----average")
	for _, key := range keys {
		average := averages[key]
		c := average.Count
		a := average.Amount / c
		d := average.Distribution / c - a * a
		fmt.Printf("%s,%d,%d,%f\n", key, c, a, math.Sqrt(float64(d)))
	}
}

func load(path string) []*Buy {
	list := []*Buy{}
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
		date, err := time.Parse("2006-01-02", cells[0])
		if err != nil {
			date, err = time.Parse("2006-01-02 15:04:05", cells[0])
			if err != nil {
				continue
			}
		}
		b := new(Buy)
		b.Date = date
		b.ID = cells[1]
		b.Price, _ = strconv.Atoi(cells[2])
		b.Title = cells[3]
		list = append(list, b)
	}
	return list
}