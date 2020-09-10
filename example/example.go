package main

import (
	"github.com/eanson023/sxcrawler"
)

func main() {
	defer sxcrawler.Done()
	rg, err := sxcrawler.Login("", "")
	if err != nil {
		panic(err)
	}
	rg.GetScore()
}
