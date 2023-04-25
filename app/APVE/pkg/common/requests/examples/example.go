package main

import (
	"bytes"
	"context"
	"fmt"
	"heaven/app/APVE/pkg/common/requests"
	"net/http"
)

type User struct {
	Id   int    `form:"id" xml:"id"`
	Name string `form:"name" xml:"name"`
}

func main() {
	endpoint := "https://www.baidu.com/"
	//直接请求
	_, _ = requests.Get(context.Background(), endpoint)
	//使用字符串map作为请求参数
	_, _ = requests.Get(context.Background(), endpoint, requests.WithParam(map[string]string{
		"id":  "1",
		"key": "abc",
	}))
	//使用结构体作为请求参数
	//参数key使用form的tag,如果不指定form直接使用field name
	query := User{
		Id:   1,
		Name: "hello",
	}
	_, _ = requests.Get(context.Background(), endpoint, requests.WithQuery(&query))
	//空body的POST请求
	_, _ = requests.Post(context.Background(), endpoint)

	//使用json格式数据作为请求body
	jsonData := User{
		Id:   2,
		Name: "json",
	}
	_, _ = requests.Post(context.Background(), endpoint, requests.WithJSON(jsonData))

	//使用xml格式数据作为请求body
	xmlData := User{
		Id:   3,
		Name: "xml",
	}
	_, _ = requests.Post(context.Background(), endpoint, requests.WithXML(xmlData))

	//直接设置请求body
	_, _ = requests.Post(context.Background(), endpoint, requests.WithBody(bytes.NewReader([]byte("hello world!"))))

	//设置请求头
	headers := map[string]string{
		"Content-Type":  "application/xml",
		"custom-header": "custom",
	}
	_, _ = requests.Post(context.Background(), endpoint, requests.WithXML(xmlData), requests.WithHeaders(headers))

	//自定义本地请求的client
	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	_, _ = requests.Post(context.Background(), endpoint, requests.WithClient(&client))

	//自定义一个请求client
	newRequests := requests.New(requests.WithClient(&client))
	_, _ = newRequests.Post(context.Background(), endpoint)

	//处理请求结果
	resp, err := requests.Get(context.Background(), endpoint)
	if err != nil {
		return
	}
	//返回值只能读取一次
	//返回内容字符串
	fmt.Println(resp.Text())
	//json
	var jsonResp interface{}
	_ = resp.JSON(&jsonResp)
	fmt.Println(jsonResp)
	//xml
	var xmlResp interface{}
	_ = resp.XML(&xmlResp)
	fmt.Println(xmlResp)
	//读取字节流，resp.Raw()返回的是一个channel
	var buffs bytes.Buffer
	for buf := range resp.Raw() {
		buffs.Write(buf)
	}
	fmt.Println(buffs.Bytes())
}

func middleware() {
	endpoint := "https://www.baidu.com/"
	var m1 = func() requests.Middleware {
		return func(next requests.Handler) requests.Handler {
			return func(client *http.Client, request *http.Request) (response *http.Response, err error) {
				ctx := context.WithValue(request.Context(), "token", "1234")
				request = request.Clone(ctx)
				fmt.Printf("method=%s, url=%s, body=%s\n", request.Method, request.URL, request.Body)
				resp, err := next(client, request)
				if err != nil {
					return resp, err
				}
				ctx = resp.Request.Context()
				fmt.Printf("token=%s\n", ctx.Value("token"))
				return resp, err
			}
		}
	}
	var m2 = func() requests.Middleware {
		return func(next requests.Handler) requests.Handler {
			return func(client *http.Client, request *http.Request) (response *http.Response, err error) {
				fmt.Println("m2 before...")
				resp, err := next(client, request)
				fmt.Println("m2 after")
				return resp, err
			}
		}
	}
	resp, err := requests.Get(context.Background(), endpoint,
		requests.WithParam(map[string]string{"id": "1"}),
		requests.WithDebug(true),
		requests.WithRetry(),
		requests.WithMiddleware(m1(), m2()))
	if err != nil {
		return
	}
	fmt.Println(resp.Text())
}
