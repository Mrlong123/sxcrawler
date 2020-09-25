package main

import (
	"github.com/eanson023/sxcrawler"
)

func main() {
	defer sxcrawler.Done()
	rg, err := sxcrawler.Login("学号", "密码")
	if err != nil {
		panic(err)
	}
	rg.GetAllCourseInfo().StoreToMarkdown("信息.md")
}
