package shttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

// 默认Client
var (
	defaultClient     *Client
	defaultClientOnce sync.Once
	defaultClientMtx  sync.Mutex
)

// Client
type Client struct {
	*http.Client
}

// NewClient
func NewClient(opts ...ClientOption) *Client {
	options := newClientOptions(opts...)

	c := &Client{Client: &http.Client{
		Timeout: options.Timeout,
	}}
	c.Transport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: options.DialTimeout,
		}).Dial,
		MaxIdleConns:        options.MaxIdleConnsPerHost,
		MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
		IdleConnTimeout:     options.MaxIdleConnTimeout,
	}

	return c
}

// 获取默认的Client
func DefaultClient() *Client {
	defaultClientMtx.Lock()
	defer defaultClientMtx.Unlock()

	// 创建DefaultClient
	defaultClientOnce.Do(func() {
		defaultClient = NewClient()
	})

	return defaultClient
}

// 封装HTTP回复
type HTTPResponse struct {
	*http.Response
	Err      error
	cancelFn context.CancelFunc
	cache    []byte // 缓存
	DoTimes  int    // 执行次数
}

// 构造http回复
func NewHTTPResponse(response *http.Response, err error, cancelFn context.CancelFunc) *HTTPResponse {
	res := &HTTPResponse{
		Response: response,
		Err:      err,
		cancelFn: cancelFn,
		DoTimes:  1,
	}

	if err != nil {
		cancelFn()
	}

	return res
}

// 读取Body
func (res *HTTPResponse) ReadAll() ([]byte, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	return res.readAll()
}

func (res *HTTPResponse) readAll() ([]byte, error) {
	// 读取缓存
	if len(res.cache) > 0 {
		return res.cache, nil
	}

	if res.Response == nil {
		return nil, errors.New("res is nil")
	}

	if res.Response.Body == nil {
		return nil, errors.New("body is nil")
	}

	if res.cancelFn != nil {
		defer res.cancelFn()
	}

	defer func() {
		res.Response.Body.Close()
		res.Response.Body = nil
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		res.Err = err
		return nil, err
	}

	res.cache = body

	return res.cache, nil
}

// JSON解析数据
func (res *HTTPResponse) JSONUnmarshal(obj interface{}) error {
	data, err := res.readAll()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}

// 关闭Body
func (res *HTTPResponse) Close() {
	if res.Response == nil {
		return
	}

	if res.Response.Body == nil {
		return
	}

	if res.cancelFn != nil {
		defer res.cancelFn()
	}

	// 读完所有的数据，用于HTTP的连接重用
	body, _ := io.ReadAll(res.Response.Body)
	res.cache = body

	res.Response.Body.Close()
	res.Response.Body = nil
}

// 关闭后，返回请求错误
func (res *HTTPResponse) CloseReturnErr() error {
	res.Close()

	return res.Err
}

// 关闭后，返回请求错误和内容
func (res *HTTPResponse) CloseReturnErrBody() error {
	res.Close()

	if len(res.cache) != 0 {
		return fmt.Errorf("%w, msg: %s", res.Err, res.cache)
	}

	return res.Err
}

// 封装HTTP请求
type HTTPRequest struct {
	options  *HTTPOptions
	ctx      context.Context
	cancelFn context.CancelFunc
	url      string
	method   string
	body     []byte
}

// 构造HTTPRequest
func (req *HTTPRequest) BuildHTTPRequest() (*http.Request, error) {
	// 设置过期时间
	requestCtx, cancelFn := context.WithDeadline(req.ctx, time.Now().Add(req.options.Timeout))

	req.cancelFn = cancelFn

	request, err := http.NewRequestWithContext(requestCtx, req.method, req.url, bytes.NewBuffer(req.body))
	if err != nil {
		return nil, err
	}

	// 设置Header
	for k, v := range req.options.Header {
		request.Header.Set(k, v)
	}

	return request, nil
}

// 构造http请求
func NewHTTPRequest(ctx context.Context, method, url string, body []byte, opts ...HTTPOption) *HTTPRequest {
	options := NewHTTPOptions(opts...)

	return &HTTPRequest{
		options: options,
		ctx:     ctx,
		method:  method,
		url:     url,
		body:    body,
	}
}

// 执行http请求
func (req *HTTPRequest) Do() *HTTPResponse {
	return req.DoWithClient(DefaultClient())
}

// 用指定Client执行
func (req *HTTPRequest) DoWithClient(c *Client) *HTTPResponse {
	res := req.do(c)
	if res.Err == nil {
		return res
	}

	// 重试
	for i := 0; i < req.options.Retry; i++ {
		// 关闭前一个Response
		res.Close()

		// 进行一次新的请求
		res = req.do(c)
		res.DoTimes += i + 1
		if res.Err == nil {
			return res
		}
	}

	return res
}

func (req *HTTPRequest) do(c *Client) *HTTPResponse {
	request, err := req.BuildHTTPRequest()
	if err != nil {
		return NewHTTPResponse(nil, err, req.cancelFn)
	}

	httpRes, err := c.Do(request)
	if err != nil {
		// 可能拿到的空闲连接已经断开，会收到EOF错误，针对这种情况自动重试一次
		if errors.Is(err, io.EOF) {
			// 关闭旧的httpRes
			if httpRes != nil && httpRes.Body != nil {
				httpRes.Body.Close()
			}
			request, _ := req.BuildHTTPRequest()
			//nolint
			httpRes, err = c.Do(request)
			if err != nil {
				return NewHTTPResponse(httpRes, err, req.cancelFn)
			}
		} else {
			return NewHTTPResponse(httpRes, err, req.cancelFn)
		}
	}

	// 判断返回的状态码是否在白名单内
	_, find := req.options.StatusWhiteList[httpRes.StatusCode]
	if !find {
		// 如果不在白名单中，则报错
		return NewHTTPResponse(httpRes, errors.New(httpRes.Status), req.cancelFn)
	}

	return NewHTTPResponse(httpRes, nil, req.cancelFn)
}

// 构造Get请求
func NewGetRequest(ctx context.Context, url string, opts ...HTTPOption) *HTTPRequest {
	return NewHTTPRequest(ctx, http.MethodGet, url, nil, opts...)
}

// 构造Post请求
func NewPostRequest(ctx context.Context, url string, contentType string, body []byte, opts ...HTTPOption) *HTTPRequest {
	opts = append(opts, WithHeader("Content-Type", contentType))
	request := NewHTTPRequest(ctx, http.MethodPost, url, body, opts...)

	return request
}

// Post
func Post(ctx context.Context, url string, contentType string, body []byte, opts ...HTTPOption) *HTTPResponse {
	return NewPostRequest(ctx, url, contentType, body, opts...).Do()
}

// Post
func (c *Client) Post(ctx context.Context, url string, contentType string, body []byte, opts ...HTTPOption) *HTTPResponse {
	return NewPostRequest(ctx, url, contentType, body, opts...).DoWithClient(c)
}

// PostJSON
func PostJSON(ctx context.Context, url string, obj interface{}, opts ...HTTPOption) *HTTPResponse {
	data, err := json.Marshal(obj)
	if err != nil {
		return &HTTPResponse{
			Err: err,
		}
	}

	return Post(context.Background(), url, "application/json", data, opts...)
}

// PostJSON
func (c *Client) PostJSON(ctx context.Context, url string, obj interface{}, opts ...HTTPOption) *HTTPResponse {
	data, err := json.Marshal(obj)
	if err != nil {
		return &HTTPResponse{
			Err: err,
		}
	}

	return c.Post(context.Background(), url, "application/json", data, opts...)
}

// JSON解码
func (res *HTTPResponse) JSONDecode(obj interface{}) error {
	data, err := res.ReadAll()
	if err != nil {
		return fmt.Errorf("read all error: %w, data: %s", err, data)
	}

	// fmt.Println(string(data))

	if err := json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("json decode error: %w, data: %s", err, data)
	}

	return nil
}

// Get
func Get(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodGet, url, nil, opts...).Do()
}

// GetJSON
func GetJSON(ctx context.Context, url string, obj interface{}, opts ...HTTPOption) *HTTPResponse {
	data, err := json.Marshal(obj)
	if err != nil {
		return &HTTPResponse{
			Err: err,
		}
	}

	opts = append(opts, WithHeader("Content-Type", "application/json"))

	return NewHTTPRequest(ctx, http.MethodGet, url, data, opts...).Do()
}

// Get
func (c *Client) Get(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodGet, url, nil, opts...).DoWithClient(c)
}

// Head
func Head(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodHead, url, nil, opts...).Do()
}

// Head
func (c *Client) Head(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodHead, url, nil, opts...).DoWithClient(c)
}

// Delete
func Delete(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodDelete, url, nil, opts...).Do()
}

// Delete
func (c *Client) Delete(ctx context.Context, url string, opts ...HTTPOption) *HTTPResponse {
	return NewHTTPRequest(ctx, http.MethodDelete, url, nil, opts...).DoWithClient(c)
}
