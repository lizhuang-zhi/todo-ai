package shttp

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPGet(t *testing.T) {
	Convey("HTTPGet", t, func() {
		res := Get(context.Background(), "http://www.baidu.com")
		So(res.Err, ShouldBeNil)

		// 第一次读取
		body, err := res.ReadAll()
		So(err, ShouldBeNil)
		So(body, ShouldNotBeNil)

		// 第二次读取
		body2, err2 := res.ReadAll()
		So(err2, ShouldBeNil)
		So(body2, ShouldResemble, body)

		// 直接关闭，不获取内容
		err = Get(context.Background(), "http://www.baidu.com").CloseReturnErr()
		So(err, ShouldBeNil)
	})
}

func TestHTTPHead(t *testing.T) {
	Convey("HTTPHead", t, func() {
		// 获取头信息
		res := Head(context.Background(), "http://www.baidu.com")
		So(res.Err, ShouldBeNil)
		body, err := res.ReadAll()
		So(err, ShouldBeNil)
		So(body, ShouldNotBeNil)

		// 访问一个不存在的网页
		res = Head(context.Background(), "http://www.baidu.com/1.txt", WithHTTPRetry(1))
		So(res.Err.Error(), ShouldEqual, "404 Not Found")
		body, err = res.ReadAll()
		So(err, ShouldEqual, res.Err)
		So(body, ShouldBeNil)
		So(res.DoTimes, ShouldEqual, 2)
	})
}

func TestHTTPPost(t *testing.T) {
	Convey("HTTPPost", t, func() {
		res := Post(context.Background(), "http://www.baidu.com", "application/x-www-form-urlencoded", nil)
		So(res.Err, ShouldBeNil)

		// 第一次读取
		body, err := res.ReadAll()
		So(err, ShouldBeNil)
		So(body, ShouldNotBeNil)

		// 第二次读取
		body2, err2 := res.ReadAll()
		So(err2, ShouldBeNil)
		So(body2, ShouldResemble, body)

		So(res.CloseReturnErr(), ShouldBeNil)
	})
}

func TestHTTPClientGet(t *testing.T) {
	Convey("HTTPGet", t, func() {
		client := NewClient(WithIdleConnTimeout(time.Millisecond))

		res := client.Get(context.Background(), "http://www.baidu.com")
		So(res.Err, ShouldBeNil)

		var body []byte
		var err error
		for i := 0; i < 100; i++ {
			body, err = res.ReadAll()
			So(err, ShouldBeNil)
			So(body, ShouldNotBeNil)
		}
	})
}
