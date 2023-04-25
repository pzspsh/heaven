package requests

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptrace"
	"os"
	"path/filepath"
)

type (
	Value map[string]string
	File  struct {
		Path   string
		Name   string
		Extras Value
	}
)

type dialOptions struct {
	err     error
	debug   bool
	client  *http.Client
	headers Value
	cookies []*http.Cookie
	query   string
	body    io.Reader
	session bool
	middles []Middleware
	retry   Middleware
	trace   *httptrace.ClientTrace
}

func (opts *dialOptions) setContentType(contentType string) {
	opts.setHeader("Content-Type", contentType)
}

func (opts *dialOptions) setHeader(k, v string) {
	if opts.headers == nil {
		opts.headers = make(Value)
	}
	opts.headers[k] = v
}

type DialOption func(opts *dialOptions)

//添加中间件
//中间件中可以获取到client,Request, Response对象，所以可以对请求做很多的操作
func WithMiddleware(middles ...Middleware) DialOption {
	return func(opts *dialOptions) {
		opts.middles = append(opts.middles, middles...)
	}
}

//是否开启debug模式
//可以在初始化的时候统一设置，也可以给每个请求单独设置
func WithDebug(debug bool) DialOption {
	return func(opts *dialOptions) {
		opts.debug = debug
	}
}

func WithParam(query Value) DialOption {
	return func(opts *dialOptions) {
		opts.query = mapToValues(query).Encode()
	}
}

//直接传递一个结构体指针作为query参数
//注意：
//1. 支持map、struct，其他的类型会直接panic，struct使用`form`指定字段名称，未指定的使用默认值
//2. 支持匿名嵌套，但不支持命名嵌套，内容不会解析，直接变成一个字符串
func WithQuery(i interface{}) DialOption {
	return func(opts *dialOptions) {
		opts.query = structToValues(i).Encode()
	}
}

func WithForm(form Value) DialOption {
	return func(opts *dialOptions) {
		s := mapToValues(form).Encode()
		opts.body = bytes.NewBufferString(s)
		opts.setContentType("application/x-www-form-urlencoded")
	}
}

func WithJSON(data interface{}) DialOption {
	return func(opts *dialOptions) {
		buf, err := json.Marshal(data)
		if err != nil {
			opts.err = err
			return
		}
		opts.body = bytes.NewBuffer(buf)
		opts.setContentType("application/json")
	}
}

func WithXML(data interface{}) DialOption {
	return func(opts *dialOptions) {
		buf, err := xml.Marshal(data)
		if err != nil {
			opts.err = err
			return
		}
		opts.body = bytes.NewBuffer(buf)
		opts.setContentType("application/xml")
	}
}

//直接设置一个请求body
func WithBody(body io.Reader) DialOption {
	return func(opts *dialOptions) {
		opts.body = body
	}
}

//上传文件
func WithFile(file *File) DialOption {
	return func(opts *dialOptions) {
		f, err := os.Open(file.Path)
		if err != nil {
			opts.err = err
			return
		}
		defer f.Close()

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile(file.Name, filepath.Base(file.Path))
		if err != nil {
			opts.err = err
			return
		}
		if _, err := io.Copy(part, f); err != nil {
			opts.err = err
			return
		}

		for k, v := range file.Extras {
			_ = writer.WriteField(k, v)
		}

		if err := writer.Close(); err != nil {
			opts.err = err
			return
		}

		opts.body = &body
		opts.setContentType(writer.FormDataContentType())
	}
}

//设置请求client
//一般在创建一个requests的时候才使用
//中间件中也可以直接替换一个client
func WithClient(clients ...*http.Client) DialOption {
	return func(opts *dialOptions) {
		if len(clients) > 0 {
			opts.client = clients[0]
		} else {
			opts.client = http.DefaultClient
		}
	}
}

//设置cookie
func WithCookies(cookies ...*http.Cookie) DialOption {
	return func(opts *dialOptions) {
		opts.middles = append(opts.middles, cookieMiddleware(cookies...))
	}
}

//是否清空cookies
//如果设置成true，后续的请求都会带上前面请求返回的cookie，所以不要随便设置，只有在确实需要的时候再设置
func WithSession(session bool) DialOption {
	return func(opts *dialOptions) {
		opts.session = session
	}
}

//设置请求重试
//自带一个默认的重试实现，可以自定义实现
func WithRetry(retries ...Retry) DialOption {
	return func(opts *dialOptions) {
		var retry Retry
		if len(retries) > 0 {
			retry = retries[0]
		} else {
			retry = &defaultRetry{}
		}
		opts.retry = retryMiddleware(newRetry(retry))
	}
}

//添加请求追踪
//trace需要自定义
func WithTrace(trace *httptrace.ClientTrace) DialOption {
	return func(opts *dialOptions) {
		opts.trace = trace
	}
}

//添加请求头
//这里设置的headers会
func WithHeaders(headers Value) DialOption {
	return func(opts *dialOptions) {
		opts.middles = append(opts.middles, headerMiddleware(headers))
	}
}
