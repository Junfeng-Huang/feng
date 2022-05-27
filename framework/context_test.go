package framework

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// IRequest接口测试
func TestContextIRequestQuery(t *testing.T) {
	req := httptest.NewRequest("",
		"http://test.com?int=1&int64=1&float64=1.1&float32=1.1&bool=true&string=string&stringslice=a&stringslice=b",
		nil)
	testContext := Context{request: req}

	intQuery, ok := testContext.QueryInt("int", 0)
	assert.True(t, ok)
	assert.Equal(t, int(1), intQuery)

	int64Query, ok := testContext.QueryInt64("int64", 0)
	assert.True(t, ok)
	assert.Equal(t, int64(1), int64Query)

	float64Query, ok := testContext.QueryFloat64("float64", 0)
	assert.True(t, ok)
	assert.Equal(t, float64(1.1), float64Query)

	float32Query, ok := testContext.QueryFloat32("float32", 0)
	assert.True(t, ok)
	assert.Equal(t, float32(1.1), float32Query)

	boolQuery, ok := testContext.QueryBool("bool", true)
	assert.True(t, ok)
	assert.Equal(t, true, boolQuery)

	stringQuery, ok := testContext.QueryString("string", "")
	assert.True(t, ok)
	assert.Equal(t, "string", stringQuery)

	stringSliceQuery, ok := testContext.QueryStringSlice("stringslice", nil)
	assert.True(t, ok)
	assert.ElementsMatch(t, []string{"a", "b"}, stringSliceQuery)
}

func TestContextIRequestParam(t *testing.T) {
	core := NewCore()
	core.GET("/test/:intParam/:int64Param/:float64Param/:float32Param/:boolParam/:stringParam")
	req := httptest.NewRequest("", "http://test.com/test/1/1/1.1/1.1/true/string", nil)
	node := core.FindRouteNodeByRequest(req)
	params := node.parseParamsFromEndNode(req.URL.Path)
	ctx := Context{params: params}

	intParam, ok := ctx.ParamInt("intParam", 0)
	assert.True(t, ok)
	assert.Equal(t, int(1), intParam)

	int64Param, ok := ctx.ParamInt64("int64Param", 0)
	assert.True(t, ok)
	assert.Equal(t, int64(1), int64Param)

	float64Param, ok := ctx.ParamFloat64("float64Param", 0)
	assert.True(t, ok)
	assert.Equal(t, float64(1.1), float64Param)

	float32Param, ok := ctx.ParamFloat32("float32Param", 0)
	assert.True(t, ok)
	assert.Equal(t, float32(1.1), float32Param)

	boolParam, ok := ctx.ParamBool("boolParam", false)
	assert.True(t, ok)
	assert.True(t, true, boolParam)

	stringParam, ok := ctx.ParamString("stringParam", "")
	assert.True(t, ok)
	assert.Equal(t, "string", stringParam)

}

func TestContextIRequestForm(t *testing.T) {
	body := bytes.NewBufferString("int=1&int64=1&float64=1.1&float32=1.1&bool=true&string=string&stringslice=a&stringslice=b")
	req := httptest.NewRequest("POST", "http://test.com/test", body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := Context{request: req}

	formInt, ok := ctx.FormInt64("int64", 0)
	assert.True(t, ok)
	assert.Equal(t, int64(1), formInt)

	formFloat64, ok := ctx.FormFloat64("float64", 0)
	assert.True(t, ok)
	assert.Equal(t, float64(1.1), formFloat64)

	formFloat32, ok := ctx.FormFloat32("float32", 0)
	assert.True(t, ok)
	assert.Equal(t, float32(1.1), formFloat32)

	formBool, ok := ctx.FormBool("bool", false)
	assert.True(t, ok)
	assert.Equal(t, true, formBool)

	formStringSlice, ok := ctx.FormStringSlice("stringslice", nil)
	assert.True(t, ok)
	assert.ElementsMatch(t, []string{"a", "b"}, formStringSlice)

	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, err := mw.CreateFormFile("file", "test")
	if assert.NoError(t, err) {
		_, err = w.Write([]byte("test"))
		assert.NoError(t, err)
	}
	mw.Close()
	req = httptest.NewRequest("POST", "http://test.com/test", buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	ctx.request = req
	fh, err := ctx.FormFile("file")
	if assert.NoError(t, err) {
		assert.Equal(t, "test", fh.Filename)
	}
	f, err := fh.Open()
	if assert.NoError(t, err) {
		content, _ := ioutil.ReadAll(f)
		t.Log(string(content))
		assert.Equal(t, "test", string(content))
	}
}

func TestContextIRequestOthers(t *testing.T) {
	body := bytes.NewBufferString("")
	body.Write([]byte(`{"Foo":"bar"}`))
	req := httptest.NewRequest("POST", "http://test.com/test", body)
	ctx := Context{request: req}
	type formTest struct {
		Foo string
	}
	j := formTest{}
	err := ctx.BindJson(&j)
	assert.NoError(t, err)
	assert.Equal(t, "bar", j.Foo)

	body = bytes.NewBuffer(nil)
	body.Write([]byte("<formTest><Foo>bar</Foo></formTest>"))
	ctx.request.Body = io.NopCloser(body)
	x := formTest{}
	err = ctx.BindXml(&x)
	assert.NoError(t, err)
	assert.Equal(t, "bar", x.Foo)

	uri := ctx.Uri()
	assert.Equal(t, req.RequestURI, uri)

	method := ctx.Method()
	assert.Equal(t, "POST", method)

	host := ctx.Host()
	assert.Equal(t, req.Host, host)

	clientIP := ctx.ClientIp()
	assert.Equal(t, "192.0.2.1:1234", clientIP)

	req = httptest.NewRequest("", "http://test.com/test", nil)
	cookie := &http.Cookie{Name: "foo", Value: "bar"}
	req.AddCookie(cookie)
	ctx.request = req
	testCookie, ok := ctx.Cookie("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", testCookie)
}

func TestContextIResponseJson(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com", nil)
	ctx := Context{responseWriter: w, request: req}
	type testOBJ struct {
		Foo string `json:foo`
	}
	obj := testOBJ{Foo: "bar"}
	ctx.Json(&obj)
	assert.Equal(t, "{\"Foo\":\"bar\"}", w.Body.String())
}

func TestContextIResponseXml(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com", nil)
	ctx := Context{responseWriter: w, request: req}
	type testOBJ struct {
		Foo string `json:foo`
	}
	obj := testOBJ{Foo: "bar"}
	w.Body = bytes.NewBuffer(nil)
	ctx.Xml(&obj)
	assert.Equal(t, "<testOBJ><Foo>bar</Foo></testOBJ>", w.Body.String())
}

func TestContextIResponseJsonp(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com/?callback=call", nil)
	ctx := Context{responseWriter: w, request: req}
	type testOBJ struct {
		Foo string `json:foo`
	}
	obj := testOBJ{Foo: "bar"}
	ctx.Jsonp(&obj)
	assert.Equal(t, "call({\"Foo\":\"bar\"});", w.Body.String())
}

func TestContextIResponseHtml(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := Context{responseWriter: w}
	type htmlTestOBJ struct {
		Name string
	}
	htlmTest := htmlTestOBJ{Name: "test"}
	w.Body = bytes.NewBuffer(nil)
	ctx.responseWriter = w
	ctx.Html("../template/test.html", htlmTest)
	assert.Equal(t, "Hello test", w.Body.String())
}

func TestContextIResponseText(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := Context{responseWriter: w}
	ctx.Text("this is %s", "test")
	assert.Equal(t, "this is test", w.Body.String())
}

func TestContextIResponseRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com/", nil)
	ctx := Context{responseWriter: w, request: req}
	ctx.Redirect("http://test.com/tests")
	assert.Equal(t, "http://test.com/tests", ctx.GetResponse().Header().Get("Location"))
}

func TestContextIResponseSetHeader(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := Context{responseWriter: w}
	ctx.SetHeader("Content-Type", "multipart/form-data")
	assert.Equal(t, "multipart/form-data", w.Header().Get("Content-Type"))
}

func TestContextIResponseSetCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := Context{responseWriter: w}
	ctx.SetCookie("a", "1", 1, "/", "test.com", false, true)
	assert.Equal(t, "a=1; Path=/; Domain=test.com; Max-Age=1; HttpOnly", w.Header().Get("Set-Cookie"))
}

func TestContextWithValue(t *testing.T) {
	var testNum int = 5
	var rec = make(chan interface{}, 5)
	req := httptest.NewRequest("", "/", nil)
	ctx := &Context{request: req}

	for i := 0; i < testNum; i++ {
		k := strconv.Itoa(i)
		testC := context.WithValue(ctx, k, i)
		go func(c context.Context, key interface{}) {
			v := c.Value(key)
			rec <- v
		}(testC, k)
	}
	count := 0
	for i := 0; i < testNum; i++ {
		v, ok := (interface{}(<-rec)).(int)
		assert.True(t, ok)
		count += v
	}
	assert.Equal(t, 0+1+2+3+4, count)
}

func TestContextWithCancel(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	ctx := &Context{request: req}
	cancelC, cancel := context.WithCancel(context.Background())
	ctx.request = ctx.request.WithContext(cancelC)
	var testNum int = 5
	for i := 0; i < testNum; i++ {
		go func(c context.Context) {
			<-ctx.Done()
		}(ctx)
	}
	cancel()
	assert.NotNil(t, <-ctx.Done())

}

func TestContextWithDeadline(t *testing.T) {
	req := httptest.NewRequest("", "/", nil)
	ctx := Context{request: req}

	d := time.Now().Add(1 * time.Second)
	deadlineC, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	ctx.request = ctx.request.WithContext(deadlineC)

	for i := 0; i < 5; i++ {
		go func(c context.Context) {
			<-c.Done()
		}(&ctx)
	}
	deadline, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, d, deadline)
	time.Sleep(2 * time.Second)
	assert.NotNil(t, <-ctx.Done())
}
