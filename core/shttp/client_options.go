package shttp

import (
	"net/http"
	"time"
)

// DefaultClientOptions 默认客户端参数
var DefaultClientOptions = ClientOptions{
	Timeout:             time.Minute,      // 默认请求超时时间
	DialTimeout:         10 * time.Second, // 默认建立连接超时时间
	MaxIdleConnsPerHost: 200,              // 默认单个Host的空闲连接数量
	MaxIdleConns:        0,                // 默认空闲连接数量，不限制
	MaxIdleConnTimeout:  10 * time.Second, // 设置空闲连接的超时时间，避免被服务器提前断开
}

// 客户端参数
type ClientOptions struct {
	Timeout             time.Duration
	DialTimeout         time.Duration
	MaxIdleConnsPerHost int
	MaxIdleConns        int
	MaxIdleConnTimeout  time.Duration
}

type ClientOption func(*ClientOptions)

func newClientOptions(opts ...ClientOption) *ClientOptions {
	// 默认配置
	options := &ClientOptions{
		Timeout:             DefaultClientOptions.Timeout,
		DialTimeout:         DefaultClientOptions.DialTimeout,
		MaxIdleConnsPerHost: DefaultClientOptions.MaxIdleConnsPerHost,
		MaxIdleConns:        DefaultClientOptions.MaxIdleConns,
		MaxIdleConnTimeout:  DefaultClientOptions.MaxIdleConnTimeout,
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}

// 设置客户端请求超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

// 设置客户端建立连接超时时间
func WithClientDialTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.DialTimeout = timeout
	}
}

// 设置每个Host最大空闲连接数
func WithMaxIdleConnsPerHost(num int) ClientOption {
	return func(opts *ClientOptions) {
		opts.MaxIdleConnsPerHost = num
	}
}

// 设置最大空闲连接数
func WithMaxIdleConns(num int) ClientOption {
	return func(opts *ClientOptions) {
		opts.MaxIdleConns = num
	}
}

// 设置空闲连接超时时间
func WithIdleConnTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.MaxIdleConnTimeout = timeout
	}
}

// http选项
type HTTPOptions struct {
	Header          map[string]string
	Timeout         time.Duration    // 过期时间，默认1分钟
	StatusWhiteList map[int]struct{} // 状态白名单，默认非200的状态就是错误
	Retry           int              // 尝试次数，默认0，只有当前一个错误的时候才会尝试
}

type HTTPOption func(*HTTPOptions)

// 超时
func WithHTTPTimeout(timeout time.Duration) HTTPOption {
	return func(opts *HTTPOptions) {
		opts.Timeout = timeout
	}
}

// 设置状态白名单
func WithHTTPStatusWhiteList(whiteList []int) HTTPOption {
	return func(opts *HTTPOptions) {
		for _, status := range whiteList {
			opts.StatusWhiteList[status] = struct{}{}
		}
	}
}

// 设置重试次数
func WithHTTPRetry(retry int) HTTPOption {
	return func(opts *HTTPOptions) {
		opts.Retry = retry
	}
}

// 设置Header
func WithHeader(key, val string) HTTPOption {
	return func(opts *HTTPOptions) {
		opts.Header[key] = val
	}
}

// http选项
func NewHTTPOptions(opts ...HTTPOption) *HTTPOptions {
	options := &HTTPOptions{
		Header:  make(map[string]string),
		Timeout: time.Minute,
		StatusWhiteList: map[int]struct{}{
			http.StatusOK: {},
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
