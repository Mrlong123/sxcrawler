# sx-crawler

> 三院教务管理系统爬虫

编写中...，完成度60%

功能:爬取你的所有课程信息、四六级成绩存储到Excel或Markdown

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
	rg.GetAllCourseInfo()
}
```

