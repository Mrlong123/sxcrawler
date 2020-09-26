# sx-crawler

> 三院(重庆三峡学院)教务管理系统爬虫

功能:爬取你的所有课程信息、四六级成绩(太懒了，不做了)存储到Markdown

## Installation

```bash
go get github.com/eanson023/sxcrawler
```

## Usage

```go
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

```

