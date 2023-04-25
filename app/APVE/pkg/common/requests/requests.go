package requests

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptrace"
)

type Http interface {
	Get(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Post(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Put(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Patch(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Delete(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Head(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Connect(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Options(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Trace(ctx context.Context, url string, opts ...DialOption) (*Response, error)
	Do(ctx context.Context, method, url string, opts ...DialOption) (*Response, error)
}

var _ Http = New()

var defaultReq = New(WithClient(http.DefaultClient))

func Get(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodGet, url, opts...)
}

func Post(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodPost, url, opts...)
}

func Put(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodPut, url, opts...)
}

func Patch(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodPatch, url, opts...)
}

func Delete(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodDelete, url, opts...)
}

func Head(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodHead, url, opts...)
}

func Connect(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodConnect, url, opts...)
}

func Options(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodOptions, url, opts...)
}

func Trace(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, http.MethodTrace, url, opts...)
}

func Do(ctx context.Context, method, url string, opts ...DialOption) (*Response, error) {
	return defaultReq.do(ctx, method, url, opts...)
}

type Client struct {
	opts    dialOptions
	Request *http.Request
}

/*
init request client
eg: New(WithClient(http.DefaultClient))
*/
func New(opts ...DialOption) *Client {
	req := &Client{}
	for _, opt := range opts {
		opt(&req.opts)
	}
	return req
}

//session
func Session(opts ...DialOption) *Client {
	opts = append(opts, WithSession(true))
	return New(opts...)
}

func (req *Client) Get(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodGet, url, opts...)
}

func (req *Client) Post(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodPost, url, opts...)
}

func (req *Client) Put(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodPut, url, opts...)
}

func (req *Client) Patch(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodPatch, url, opts...)
}

func (req *Client) Delete(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodDelete, url, opts...)
}

func (req *Client) Head(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodHead, url, opts...)
}

func (req *Client) Connect(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodConnect, url, opts...)
}

func (req *Client) Options(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodOptions, url, opts...)
}

func (req *Client) Trace(ctx context.Context, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, http.MethodTrace, url, opts...)
}

func (req *Client) Do(ctx context.Context, method, url string, opts ...DialOption) (*Response, error) {
	return req.do(ctx, method, url, opts...)
}

func (req Client) do(ctx context.Context, method string, url string, opts ...DialOption) (*Response, error) {
	for _, opt := range opts {
		opt(&req.opts)
	}

	// set params error
	if req.opts.err != nil {
		return nil, req.opts.err
	}

	if req.opts.query != "" {
		url += "?" + req.opts.query
	}

	if req.opts.trace != nil {
		ctx = httptrace.WithClientTrace(ctx, req.opts.trace)
	}

	request, err := http.NewRequestWithContext(ctx, method, url, req.opts.body)
	if err != nil {
		return nil, err
	}

	//set request headers
	for k, v := range req.opts.headers {
		request.Header.Set(k, v)
	}

	//set http client
	client := req.opts.client
	if client == nil {
		client = http.DefaultClient
		req.opts.client = client
	}

	//set cookies
	if !req.opts.session {
		client.Jar = nil
	}

	//debug
	if req.opts.debug {
		//debug的日志中间件放在最外层
		req.opts.middles = append([]Middleware{loggerMiddleware()}, req.opts.middles...)
	}
	//将重试的中间件加到最后一个，这样重试的时候执行的代码就是最后请求的方法
	if req.opts.retry != nil {
		req.opts.middles = append(req.opts.middles, req.opts.retry)
	}
	r, err := Chain(req.opts.middles...)(exec)(client, request)
	resp := &Response{resp: r}

	return resp, err
}

func exec(client *http.Client, request *http.Request) (*http.Response, error) {
	c := make(chan error)
	var (
		resp *http.Response
		err  error
	)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("http requests panic: \n%+v\n", r)
				switch x := r.(type) {
				case string:
					c <- errors.New(x)
				case error:
					c <- x
				default:
					c <- errors.New("unknown panic")
				}
			}
		}()

		resp, err = client.Do(request)
		c <- err
	}()

	ctx := request.Context()
	select {
	case <-ctx.Done():
		<-c
		return nil, ctx.Err()
	case err := <-c:
		return resp, err
	}
}
