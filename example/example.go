package main

import (
	"github.com/eanson023/sxcrawler"
)

func main() {
	defer sxcrawler.Done()
	_, err := sxcrawler.Login("账号", "密码")
	if err != nil {
		panic(err)
	}
}
