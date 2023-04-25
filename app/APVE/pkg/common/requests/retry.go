package requests

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type (
	RetryCondition func(resp *http.Response, err error) bool
	RetryPolicy    func(resp *http.Response, min, max time.Duration, attempt int) (time.Duration, error)
	Retry          interface {
		MaxEntries() int
		WaitTime() time.Duration
		MaxWaitTime() time.Duration
		RetryConditions() []RetryCondition
		RetryPolicy() RetryPolicy
	}
)

type retry struct {
	maxRetries  int
	waitTime    time.Duration
	maxWaitTime time.Duration
	conditions  []RetryCondition
	policy      RetryPolicy
}

func newRetry(iretry Retry) *retry {
	return &retry{
		maxRetries:  iretry.MaxEntries(),
		waitTime:    iretry.WaitTime(),
		maxWaitTime: iretry.MaxWaitTime(),
		conditions:  iretry.RetryConditions(),
		policy:      iretry.RetryPolicy(),
	}
}

func (r *retry) backoff(fn func() (*http.Response, error)) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		fmt.Printf("retry attempt: %d\n", attempt)
		var (
			ctx       context.Context
			needRetry bool
		)
		resp, err = fn()
		if resp != nil {
			ctx = resp.Request.Context()
		}
		if ctx == nil {
			ctx = context.Background()
		}
		// context异常直接返回无需重试
		if ctx.Err() != nil {
			return nil, err
		}
		for _, condition := range r.conditions {
			if condition(resp, err) {
				needRetry = true
				break
			}
		}
		if !needRetry {
			return resp, err
		}
		//计算重试间隔时间
		waitTime, err := r.policy(resp, r.waitTime, r.maxWaitTime, attempt)
		if err != nil {
			return nil, err
		}
		timer := time.NewTimer(waitTime)
		select {
		case <-timer.C:
		case <-ctx.Done():
			//close timer and drain the channel
			if !timer.Stop() {
				<-timer.C
			}
			return nil, ctx.Err()
		}
	}
	return resp, err
}

type defaultRetry struct {
}

func (r *defaultRetry) MaxEntries() int {
	return 3
}
func (r *defaultRetry) WaitTime() time.Duration {
	return 100 * time.Microsecond
}
func (r *defaultRetry) MaxWaitTime() time.Duration {
	return 10 * time.Second
}
func (r *defaultRetry) RetryConditions() []RetryCondition {
	return []RetryCondition{
		func(resp *http.Response, err error) bool {
			if err != nil {
				return true
			}
			return false
		},
	}
}
func (r *defaultRetry) RetryPolicy() RetryPolicy {
	return func(resp *http.Response, min, max time.Duration, attempt int) (time.Duration, error) {
		waitTime := min + time.Duration(attempt)*min
		if waitTime > max {
			waitTime = max
		}
		return waitTime, nil
	}
}
