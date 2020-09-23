package markdown

import ()
import "os"
import "errors"
import "strings"
import "fmt"
import "bytes"

type Heading byte

const (
	First Heading = 1 + iota
	Second
	Third
	Fourth
	Fifth
	Sixth
)

type Handler interface {
	// 创建
	Build(buffer *bytes.Buffer) error
	// 添加数据
	Add(data string)
}

type Markdown struct {
	filename string
	buf      *bytes.Buffer
	hanlders []Handler
}

// Table 表格
type Table struct {
	row   int
	col   int
	texts []*Text
	// 实际存储大小
	size int
}

// Title 标题
type Title struct {
	heading Heading
	text    *Text
}

// Text 基本文本
type Text struct {
	line string
}

// List Ul Ol实现
type List interface {
}

// Ul 无序列表
type Ul struct {
	lies []*Li
}

// Block 区块
type Block struct {
	text *Text
}

// Code 代码块
type Code struct {
	language string
	code     *Text
}

// Ol 有序列表
type Ol struct {
	lies []*Li
}

// Li 每个小节点
type Li struct {
	parent, child  List
	prevLi, nextLi *Li
	text           Text
}

func New(filename string) *Markdown {
	return &Markdown{
		filename: filename,
		buf:      &bytes.Buffer{},
	}
}

func NewText() *Text {
	return new(Text)
}

func NewTitle(heading Heading) *Title {
	return &Title{
		heading: heading,
	}
}

func NewTable(row int, col int) *Table {
	return &Table{
		row:   row,
		col:   col,
		texts: make([]*Text, row*col),
		size:  0,
	}
}

// 添加存储链
func (md *Markdown) Join(handlers ...Handler) *Markdown {
	for _, handle := range handlers {
		md.hanlders = append(md.hanlders, handle)
	}
	return md
}

// 进行存储
func (md *Markdown) Store() error {
	for _, handler := range md.hanlders {
		err := handler.Build(md.buf)
		if err != nil {
			return err
		}
	}
	file, err := os.Create(md.filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = md.buf.WriteTo(file)
	return err
}

func (tt *Text) Add(data string) {
	tt.line = data
}

func (title *Title) Add(data string) {
	if title.text == nil {
		title.text = NewText()
	}
	title.text.line = data
}

func writeLine(data string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(data + "\r\n")
	return err
}

func (tt *Text) Build(buf *bytes.Buffer) error {
	return writeLine(tt.line, buf)
}

func (title *Title) Build(buf *bytes.Buffer) error {
	orig := title.text.line
	title.text.line = fmt.Sprintf("%s %s", strings.Repeat("#", int(title.heading)), orig)
	return title.text.Build(buf)
}

// 多了将会自动丢弃
func (table *Table) Add(data string) {
	length := table.row * table.col
	if table.size >= length {
		return
	}
	index := table.size
	table.texts[index].line = data
	table.size++
}

func (table *Table) Update(data string, rowIdx int, colIdx int) error {
	if rowIdx < 0 || rowIdx >= table.row || colIdx < 0 || colIdx >= table.col {
		return errors.New("The width or length index of the table must be within the specified range of the table")
	}
	table.texts[rowIdx*table.row+colIdx].line = data
	return nil
}
