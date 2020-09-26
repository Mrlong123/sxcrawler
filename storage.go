package sxcrawler

import (
	"fmt"
	"github.com/eanson023/mkdown"
	"io/ioutil"
	"strconv"
)

func (si *studentInfo) StoreToMarkdown(fileName string) {
	md := mkdown.New(fileName)
	si.writeTitle(md)
	si.writeUsualInfo(md)
	si.writeSemesterInfo(md)
	md.Store()
	bytes, _ := ioutil.ReadFile("banner.txt")
	fmt.Println(string(bytes))
}

func (si *studentInfo) writeTitle(md *mkdown.Markdown) {
	title := mkdown.NewTitleWithText(mkdown.Heading1, "sxcrawler")
	titleSubBlock1 := mkdown.NewBlock("厚德  博学  自强  创新")
	intro := fmt.Sprintf("author: %s , 请为我star(#^.^#)", mkdown.NewLink("eanson", "https://github.com/eanson023").String())
	titleSubBlock2 := mkdown.NewBlock(intro)
	md.Join(title, titleSubBlock1, titleSubBlock2)
}

// writeUsualInfo 写入学生基本信息
func (si *studentInfo) writeUsualInfo(md *mkdown.Markdown) {
	title := mkdown.NewTitleWithText(mkdown.Heading2, "基本信息")
	table := mkdown.NewTable(2, 6)
	table.Add("学号").Add("姓名").Add("学院").Add("专业").Add("班级").Add("所选学分")
	table.Add(si.stuID).Add(si.name).Add(si.college).Add(si.major).Add(si.grade)
	var score float32 = 0
	for _, semester := range si.semesters {
		score += semester.selectedCredit
	}
	// 添加总学分 float转string
	table.Add(strconv.FormatFloat(float64(score), 'f', -1, 32))
	md.Join(title, table)
}

// writeSemesterInfo 写入每学期的信息
func (si *studentInfo) writeSemesterInfo(md *mkdown.Markdown) {
	title1 := mkdown.NewTitleWithText(mkdown.Heading2, "学业信息")
	md.Join(title1)
	for _, semester := range si.semesters {
		title2 := mkdown.NewTitleWithText(mkdown.Heading3, fmt.Sprintf("%s 学年 第 %d 学期", semester.year, semester.semester))
		text := mkdown.NewText(fmt.Sprintf("%s:%v\t%s:%v\t%s:%v", mkdown.NewStrongString("所选学分"), semester.selectedCredit, mkdown.NewStrongString("所获学分"), semester.gainCredit, mkdown.NewStrongString("重修学分"), semester.retakeCredit))
		tableRow := len(semester.courses) + 1
		table := mkdown.NewTable(tableRow, 13)
		table.Add("课程代码").Add("课程名称").Add("课程性质").Add("课程归属").Add("学分").Add("绩点").Add("成绩").Add("辅修标记").Add("补考成绩").Add("重修成绩").Add("学院名称").Add("备注").Add("重修标记")
		for _, course := range semester.courses {
			table.Add(course.courseCode).Add(course.courseName).Add(course.courseNature).Add(course.belong).Add(course.credit).Add(course.gradePoint)
			table.Add(course.score).Add(course.minorMark).Add(course.retestScore).Add(course.retakeScore).Add(course.collegeName).Add(course.remarks).Add(course.retakeMark)
		}
		line := mkdown.NewText("\r\n")
		md.Join(title2, text, table, line)
	}
}
