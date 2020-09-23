package main

import (
	"github.com/eanson023/sxcrawler/markdown"
	"testing"
)

func TestMarkdown(t *testing.T) {
	md := markdown.New("test.md")
	title := markdown.NewTitle(markdown.First)
	title.Add("测试")
	text := markdown.NewText()
	text.Add("ssss")
	md.Join(title, text)
	md.Store()
}
