package sxcrawler

import "net/http"

//RequestGeneral ...
type RequestGeneral struct {
	stu     *student
	cookies []*http.Cookie
	headers map[string]string
}

//loginForm 登录表单
type loginForm struct {
	// 不知道是啥
	__VIEWSTATE string
	// 学号
	TextBox1 string
	// 密码
	TextBox2 string
	// 验证码
	TextBox3 string
	//角色
	RadioButtonList1 string
	Button1          string
}

//student 学生信息
type student struct {
	xh       string
	password string
}

//成绩查询表单
type scoreForm struct {
	__VIEWSTATE string
	// 学年
	ddlXN string
	// 学期
	ddlXQ string
	// 安学期查询
	Button1 string
}

type studentInfo struct {
	// 学号
	stuID string
	// 姓名
	name string
	// 学院
	college string
	// 专业
	major string
	// 班级
	grade     string
	semesters []*semester
}

// 每学期信息
type semester struct {
	// 学年
	year string
	// 学期
	semester int
	// 所选学分
	selectedCredit float32
	// 所获学分
	gainCredit float32
	// 重修学分
	retakeCredit float32
	// 每学期的课程
	courses []*course
}

// course 课程信息
type course struct {
	// 课程代码
	courseCode string
	// 课程名称
	courseName string
	// 课程性质
	courseNature string
	// 课程归属
	belong string
	// 学分
	credit string
	// 绩点
	gradePoint string
	// 成绩
	score string
	// 辅修标记
	minorMark string
	// 补考成绩
	retestScore string
	// 重修成绩
	retakeScore string
	// 学院名称
	collegeName string
	// 备注
	remarks string
	// 重修标记
	retakeMark string
}
