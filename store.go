package sxcrawler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
)

// 存储接口
type Store interface {
	StoreToExcel(path string)
	StoreToMarkdown(path string)
}

func (si *studentInfo) StoreToExcel(path string) *studentInfo {
	header1 := map[string]string{"A1": "姓名", "B1": "学号", "C1": "学院", "D1": "专业", "E1": "班级"}
	f := excelize.NewFile()
	for k, v := range header1 {
		f.SetCellValue("Sheet1", k, v)
	}
	basic := map[string]string{"A2": si.name, "B2": si.stuID, "C2": si.college, "D2": si.major, "E2": si.grade}
	for k, v := range basic {
		f.SetCellValue("Sheet1", k, v)
	}
	idx := 4
	for _, v := range si.semesters {
		str := fmt.Sprintf("%v学年第%v学期学习成绩\t\t所选学分%v;获得学分%v;重修学分%v", v.year, v.semester, v.selectedCredit, v.gainCredit, v.retakeCredit)
		idxStr := strconv.Itoa(idx)
		interval := fmt.Sprintf("A%v:E%v", idxStr, idxStr)
		f.SetCellValue("Sheet1", interval, str)
		idx += 2
	}
	if err := f.SaveAs(path); err != nil {
		panic(err)
	}
	return si
}

func (si *studentInfo) StoreToMarkdown(path string) *studentInfo {
	return nil
}
