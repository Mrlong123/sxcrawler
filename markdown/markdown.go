package markdown

import ()
import "strconv"
import "os"
import "errors"
import "strings"
import "fmt"
import "bytes"

type Heading byte

const (
	Heading1 Heading = 1 + iota
	Heading2
	Heading3
	Heading4
	Heading5
	Heading6
)

// Hanlder 接口用于将内容写入bytes.Buffer
type Handler interface {
	// 创建
	Build(buffer *bytes.Buffer) error
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
	// 表格中内容最大长度大小
	maxLength int
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

// List Ul Ol接口
type List interface {
	AppendOne(data string)
}

// Ul 无序列表
type UnList struct {
	lies []*Li
}

// Ol 有序列表
type OrderList struct {
	lies []*Li
}

// Li 每个小节点
type Li struct {
	parent, child  List
	prevLi, nextLi *Li
	text           Text
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

// New 创建新的markdown文档
func New(filename string) *Markdown {
	return &Markdown{
		filename: filename,
		buf:      &bytes.Buffer{},
	}
}

// NewText 创建新的文本
func NewText(data string) *Text {
	return &Text{
		line: data,
	}
}

// NewTitle 创建新的标题
func NewTitle(heading Heading) *Title {
	return &Title{
		heading: heading,
	}
}

// NewTable 创建新的表格
func NewTable(row int, col int) *Table {
	return &Table{
		row:       row,
		col:       col,
		texts:     make([]*Text, row*col),
		size:      0,
		maxLength: 0,
	}
}

// Join 添加handler到存储链
func (md *Markdown) Join(handlers ...Handler) *Markdown {
	for _, handle := range handlers {
		md.hanlders = append(md.hanlders, handle)
	}
	return md
}

// Store 进行存储
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

// Append 追加文字
func (tt *Text) Append(data string) {
	tt.line += data
}

// SetTitle 设置文本
func (title *Title) SetTitle(data string) {
	if title.text == nil {
		title.text = NewText(data)
	}
	title.text.line = data
}

func writeLine(data string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(data + "\r\n")
	return err
}

// Build 写入文本
func (tt *Text) Build(buf *bytes.Buffer) error {
	return writeLine(tt.line, buf)
}

// Build 写入标题
func (title *Title) Build(buf *bytes.Buffer) error {
	orig := title.text.line
	title.text.line = fmt.Sprintf("%s %s", strings.Repeat("#", int(title.heading)), orig)
	return title.text.Build(buf)
}

// Build 创建表格
func (table *Table) Build(buf *bytes.Buffer) (err error) {
	//header
	//eg:%-99s
	var format string = "| %-" + strconv.Itoa(table.maxLength) + "s "
	for i := 0; i < table.col; i++ {
		_, err = buf.WriteString(fmt.Sprintf(format, table.texts[i].line))
	}
	if err != nil {
		return err
	}
	writeLine("|", buf)
	// 分割线
	buf.WriteString(strings.Repeat(fmt.Sprintf("| %s ", strings.Repeat("-", table.maxLength)), table.col) + "|\r\n")
	// 内容
	for r := 1; r < table.row; r++ {
		for c := 0; c < table.col; c++ {
			if _, err = buf.WriteString(fmt.Sprintf(format, table.texts[r*table.col+c].line)); err != nil {
				return err
			}
		}
		writeLine("|", buf)
	}
	return nil
}

// Add 添加到表格得格子中 多了将会返回异常
func (table *Table) Add(data string) error {
	length := table.row * table.col
	if table.size >= length {
		return errors.New("the table size will overflow")
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	index := table.size
	table.texts[index] = NewText(data)
	table.size++
	return nil
}

// AddIgnoreError 直接链式编程但忽略了错误 多了会返回nil
func (table *Table) AddIgnoreError(data string) *Table {
	length := table.row * table.col
	if table.size >= length {
		return nil
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	index := table.size
	table.texts[index] = NewText(data)
	table.size++
	return table
}

// Update 更新某个格子内容
func (table *Table) Update(data string, rowIdx int, colIdx int) error {
	if rowIdx < 0 || rowIdx >= table.row || colIdx < 0 || colIdx >= table.col {
		return errors.New("The width or length index of the table must be within the specified range of the table")
	}
	if len(data) > table.maxLength {
		table.maxLength = len(data)
	}
	table.texts[rowIdx*table.row+colIdx].line = data
	return nil
}
