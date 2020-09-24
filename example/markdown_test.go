package main

import (
	"github.com/eanson023/sxcrawler/markdown"
	"testing"
)

func TestMarkdown(t *testing.T) {
	md := markdown.New("test.md")
	title := markdown.NewTitle(markdown.Heading1)
	title.SetTitle("测试")
	text := markdown.NewText("ssss")
	md.Join(title, text)
	table := markdown.NewTable(2, 2)
	table.AddIgnoreError("1").AddIgnoreError("2222222").AddIgnoreError("3333333").AddIgnoreError("4444444")
	md.Join(table).Store()
}
