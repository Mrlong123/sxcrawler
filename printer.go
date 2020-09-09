package sxcrawler

import (
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"
)

//为了更好的爬取，写一个HTTP格式化输出语句 没写对说明没封装好

type encoding byte

const (
	//UTF8 UTF8字符
	UTF8 encoding = 1 + iota
	//GBK 国标码
	GBK
)

// PrintResponse 打印响应信息
func PrintResponse(response *http.Response, charset encoding) {
	//HTTP/1.1 302 Found
	fmt.Printf("%s %d %s\r\n", response.Proto, response.StatusCode, response.Status)
	printHeader(response.Header)
	switch charset {
	case UTF8:
		printBody(response.Body)
	case GBK:
		printBodyGBK(response.Body)
	}
}

// PrintRequest 打印请求信息
func PrintRequest(request *http.Request, charset encoding) {
	// POST /default2.aspx HTTP/1.1
	fmt.Printf("%s %s %s\r\n", request.Method, request.RequestURI, request.Proto)
	printHeader(request.Header)
	switch charset {
	case UTF8:
		printBody(request.Body)
	case GBK:
		printBodyGBK(request.Body)
	}
}

func printHeader(header http.Header) {
	for key, values := range header {
		fmt.Printf("%s: %s\r\n", key, strings.Join(values, ","))
	}
}

func printlnCookie(cookies []*http.Cookie) {
	for _, cookie := range cookies {
		fmt.Println(cookie.String())
	}
}

func printBody(body io.ReadCloser) {
	defer body.Close()
	sb := readBody(body)
	if sb == nil {
		return
	}
	fmt.Printf("\r\n%s\r\n", sb.String())
}

func printBodyGBK(body io.ReadCloser) {
	reader, err := charset.NewReaderLabel("GBK", body)
	if err != nil {
		panic(err)
	}
	sb := readBody(reader)
	if sb == nil {
		return
	}
	fmt.Printf("\r\n%s\r\n", sb.String())
}

func readBody(body io.Reader) *strings.Builder {
	if body == nil {
		return nil
	}
	var sb strings.Builder
	buf := make([]byte, 1024)
	for {
		n, err := body.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		sb.Write(buf[:n])
		if err == io.EOF {
			break
		}
	}
	return &sb
}
