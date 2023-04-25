package requests

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
)

type Handler func(client *http.Client, request *http.Request) (response *http.Response, err error)
type Middleware func(next Handler) Handler

//compress middleware
func Chain(middles ...Middleware) Middleware {
	return func(next Handler) Handler {
		for i := len(middles) - 1; i >= 0; i-- {
			next = middles[i](next)
		}
		return next
	}
}

func loggerMiddleware() Middleware {
	return func(next Handler) Handler {
		return func(client *http.Client, request *http.Request) (*http.Response, error) {
			if buf, err := httputil.DumpRequest(request, true); err == nil {
				log.Printf("[request]\n%s\n", buf)
			}
			resp, err := next(client, request)
			if err != nil {
				log.Printf("[response] error: %+v\n", err)
			} else {
				if buf, err := httputil.DumpResponse(resp, true); err == nil {
					log.Printf("[response]\n%s\n", buf)
				}
			}
			return resp, err
		}
	}
}

func retryMiddleware(retry *retry) Middleware {
	return func(next Handler) Handler {
		return func(client *http.Client, request *http.Request) (response *http.Response, err error) {
			return retry.backoff(func() (response *http.Response, err error) {
				return next(client, request)
			})
		}
	}
}

func headerMiddleware(headers Value) Middleware {
	return func(next Handler) Handler {
		return func(client *http.Client, request *http.Request) (response *http.Response, err error) {
			for k, v := range headers {
				request.Header.Set(k, v)
			}
			return next(client, request)
		}
	}
}

func cookieMiddleware(cookies ...*http.Cookie) Middleware {
	return func(next Handler) Handler {
		return func(client *http.Client, request *http.Request) (response *http.Response, err error) {
			if len(cookies) > 0 {
				if client.Jar == nil {
					jar, _ := cookiejar.New(nil)
					client.Jar = jar
				}
				client.Jar.SetCookies(request.URL, cookies)
			}
			return next(client, request)
		}
	}
}
