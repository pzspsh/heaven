package requests

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	resp *http.Response
}

func (r *Response) Close() error {
	return r.resp.Body.Close()
}

//原始的Response对象
func (r *Response) Response() *http.Response {
	return r.resp
}

//json内容
func (r *Response) JSON(i interface{}) error {
	return r.withClose(func() error {
		return json.NewDecoder(r.resp.Body).Decode(i)
	})
}

//XML格式内容
func (r *Response) XML(i interface{}) error {
	return r.withClose(func() error {
		return xml.NewDecoder(r.resp.Body).Decode(i)
	})
}

//二进制内容
func (r *Response) Content() ([]byte, error) {
	defer r.resp.Body.Close()
	return ioutil.ReadAll(r.resp.Body)
}

//字符串内容
func (r *Response) Text() (string, error) {
	buf, err := r.Content()
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

//读取字节流
//读取完或者遇到错误，会关闭channel
func (r *Response) Raw() chan []byte {
	ch := make(chan []byte)
	buf := make([]byte, 1024*1024)
	go func() {
		defer func() {
			close(ch)
			_ = r.resp.Body.Close()
		}()

		for {
			n, err := r.resp.Body.Read(buf)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}
			ch <- buf[:n]
		}
	}()
	return ch
}

func (r *Response) withClose(f func() error) error {
	defer r.resp.Body.Close()
	return f()
}
