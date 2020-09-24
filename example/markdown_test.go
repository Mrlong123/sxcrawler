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
	table.Add("1")
	table.Add("2222222")
	table.Add("323213132")
	table.Add("444444")
	md.Join(table)
	md.Store()
}
