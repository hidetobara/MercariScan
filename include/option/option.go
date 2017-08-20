package option

import (
	"os"
	"path/filepath"
	"fmt"
	"bufio"
	"strings"
)

type Paramters struct {
	Path string
	Table map[string]string
}

var single *Paramters

func (p *Paramters) Initialize() {
	p.Table = map[string]string{}
	p.Table["url"] = "https://www.mercari.com/jp/search/?sort_order=&keyword=switch&price_min=30000&price_max=60000&status_all=1&status_trading_sold_out=1&page=2"
	p.Table["data_dir"] = "C:/obara/MercariScan/data/"
}
func (p *Paramters) Get(key string) string {
	return p.Table[key]
}

func Load() *Paramters {
	if single != nil {
		return single
	}

	p := new(Paramters)
	p.Initialize()
	single = p
	fullpath, _ := filepath.Abs(os.Args[0])
	dirpath := filepath.Dir(fullpath)
	optpath := dirpath + "/go.def"
	p.Path = optpath

	file, err := os.Open(optpath)
	if err != nil {
		fmt.Println("option is missed.")
		return p
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cells := strings.Split(scanner.Text(), ":=")
		if len(cells) != 2 {
			continue
		}
		p.Table[cells[0]] = cells[1]
		fmt.Println(cells[0], "=", cells[1])
	}
	return p
}