package sxcrawler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"
)

//_URL 网址
const _URL = "http://jwgl.sanxiau.edu.cn"

var fileName string = "checkcode.gif"

// 发送请求的客户端
var client *http.Client = new(http.Client)

//RequestGeneral ...
type RequestGeneral struct {
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

func newRequestGeneral() *RequestGeneral {
	return &RequestGeneral{
		cookies: []*http.Cookie{},
		headers: map[string]string{
			"User-Agent":      " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
			"Accept-Language": " zh-CN,zh;q=0.9,en;q=0.8",
		},
	}
}

//AddCookie 添加cookie进ReqHeader
func (rg *RequestGeneral) AddCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		rg.cookies = append(rg.cookies, cookie)
	}
}

//AddHeader 添加Header
func (rg *RequestGeneral) AddHeader(key string, value string) {
	rg.headers[key] = value
}

// 将RequestHeader导入到http请求中
func (rg *RequestGeneral) headerIntoRequest(request *http.Request) {
	// 添加requestHeader
	for key, value := range rg.headers {
		request.Header.Add(key, value)
	}
}

//将cookie加到请求中
func (rg *RequestGeneral) cookieIntoRequest(request *http.Request) {
	// 添加cookie
	for _, cookie := range rg.cookies {
		request.AddCookie(cookie)
	}
}

// fmt.Printf("", var)
func (rg *RequestGeneral) newEmptyBodyRequest(method string, url string) *http.Request {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	// 将基本请求头加入到请求中
	rg.headerIntoRequest(request)
	// 将cookie加入请求中
	rg.cookieIntoRequest(request)
	return request
}

func (rg *RequestGeneral) newRequest(method string, url string, body io.Reader) *http.Request {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	// 将基本请求头加入到请求中
	rg.headerIntoRequest(request)
	// 将cookie加入请求中
	rg.cookieIntoRequest(request)
	return request
}

func (lf *loginForm) generateForm(reqGeneral *RequestGeneral) *bytes.Buffer {
	// x-www-form-urlencoded
	payload := &bytes.Buffer{}
	getValue := reflect.ValueOf(lf).Elem()
	getType := reflect.TypeOf(lf).Elem()
	// 反射添加进条件中
	params := []string{}
	for i := 0; i < getValue.NumField(); i++ {
		str := getType.Field(i).Name + "=" + url.QueryEscape(getValue.Field(i).String())
		params = append(params, str)
	}
	ret := strings.Join(params, "&")
	payload.ReadFrom(strings.NewReader(ret))
	reqGeneral.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	return payload
}

//newLoginForm 创建新的登录表单
func newLoginForm(username string, password string, checkcode string) *loginForm {
	return &loginForm{
		// 固定值
		__VIEWSTATE: "dDw3OTkxMjIwNTU7Oz6bmpbeSO1k01TBeZU9nxNbmYM4aw==",
		TextBox1:    username,
		TextBox2:    password,
		TextBox3:    checkcode,
		// 学生
		RadioButtonList1: "学生",
	}
}

//Login ... 程序入口
func Login(username string, password string) (*RequestGeneral, error) {
	fmt.Println("请注意不同IDE Console接收键盘输入问题😁")
	rg := newRequestGeneral()
	cookies, err := getJwglCookies(rg)
	if err != nil {
		panic(err)
	}
	// 添加全局cookie
	rg.AddCookie(cookies)
	// 睡眠
	time.Sleep(time.Millisecond * 500)
	// 输入验证码
	checkcode := inputCheckCode(rg)
	// 制造表单
	loginForm := newLoginForm(username, password, checkcode)
	// 再次睡眠
	time.Sleep(time.Millisecond * 500)
	// 登录
	err = login(rg, loginForm)
	if err != nil {
		return nil, err
	}
	return rg, nil
}

// 获取教务系统cookie 里面包含sessionID
func getJwglCookies(reqGeneral *RequestGeneral) ([]*http.Cookie, error) {
	request := reqGeneral.newEmptyBodyRequest("GET", _URL)
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	return resp.Cookies(), nil
}

// 1. 获取cookie,通过cookie获取该验证码
func inputCheckCode(reqGeneral *RequestGeneral) string {
	var checkcode string
	checkCodeURL := "http://jwgl.sanxiau.edu.cn/CheckCode.aspx?"
	request := reqGeneral.newEmptyBodyRequest("GET", checkCodeURL)
	if resp, err := client.Do(request); err != nil {
		panic(err)
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		err = ioutil.WriteFile(fileName, body, os.ModePerm)
		if err != nil {
			fmt.Println("写入验证码文件出错")
			panic(err)
		}
		fmt.Println("请查看文件夹目录里的" + fileName + "文件")
		for checkcode == "" {
			fmt.Print("请输入验证码(回车结束):")
			fmt.Scanln(&checkcode)
		}
	}
	return checkcode
}

// 2. 模拟表单登录，随后根据Location跳转获取信息
// return 是否登录成功和响应结果
func login(reqGeneral *RequestGeneral, form *loginForm) error {
	loginURL := "http://jwgl.sanxiau.edu.cn/default2.aspx"
	// 添加
	param := form.generateForm(reqGeneral)
	request := reqGeneral.newRequest("POST", loginURL, param)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	reader, err := charset.NewReaderLabel("GBK", resp.Body)
	if err != nil {
		panic(err)
	}
	// 用goquery解析html 判断是否登录成功
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	// 如果登录失败了 网页上会有错误信息
	scriptText := doc.Find("script").Last().Text()
	wrong := strings.Contains(scriptText, "alert")
	if wrong {
		// 返回错误信息
		errorMsg := scriptText[strings.Index(scriptText, "(")+1 : strings.LastIndex(scriptText, ")")]
		return errors.New(errorMsg)
	}
	return nil
}

//Done 结束删除验证码文件
func Done() {
	os.Remove(fileName)
}
