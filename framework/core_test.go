package framework

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCoreHTTPMethods(t *testing.T) {
	core := NewCore()
	core.GET("/test", func(c *Context) {
		method := c.GetRequest().Method
		c.Text("%s", method)
	})
	core.POST("/test", func(c *Context) {
		method := c.GetRequest().Method
		c.Text("%s", method)
	})
	core.PUT("/test", func(c *Context) {
		method := c.GetRequest().Method
		c.Text("%s", method)
	})
	core.DELETE("/test", func(c *Context) {
		method := c.GetRequest().Method
		c.Text("%s", method)
	})
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	time.Sleep(1 * time.Second)
	resp, err := http.Get("http://localhost:8080/test")
	assert.Nil(t, err)
	buffer := bytes.NewBuffer(nil)
	io.Copy(buffer, resp.Body)
	assert.Equal(t, "GET", buffer.String())

	resp, err = http.Post("http://localhost:8080/test", "", nil)
	assert.Nil(t, err)
	buffer = bytes.NewBuffer(nil)
	io.Copy(buffer, resp.Body)
	assert.Equal(t, "POST", buffer.String())

	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/test", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	buffer = bytes.NewBuffer(nil)
	io.Copy(buffer, resp.Body)
	assert.Equal(t, "PUT", buffer.String())

	req, err = http.NewRequest(http.MethodDelete, "http://localhost:8080/test", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.Nil(t, err)
	buffer = bytes.NewBuffer(nil)
	io.Copy(buffer, resp.Body)
	assert.Equal(t, "DELETE", buffer.String())
}

func TestCoreStaticRouting(t *testing.T) {
	core := NewCore()
	core.GET("/test", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	g := core.Group("/user")
	g.GET("/login", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	g.GET("/register", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	time.Sleep(time.Second)
	resp, err := http.Get("http://127.0.0.1:8080/test")
	assert.Nil(t, err)
	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/test", body.String())

	resp, err = http.Get("http://127.0.0.1:8080/user/login")
	assert.Nil(t, err)
	body = bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/user/login", body.String())

	resp, err = http.Get("http://127.0.0.1:8080/user/register")
	assert.Nil(t, err)
	body = bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/user/register", body.String())
}

func TestCoreDynamicRouting(t *testing.T) {
	core := NewCore()
	id := 1
	core.GET("/test/:id", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	srv := http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.TODO())
	time.Sleep(time.Second)
	resp, err := http.Get("http://127.0.0.1:8080/test/" + strconv.Itoa(id))
	assert.Nil(t, err)
	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/test/"+strconv.Itoa(id), body.String())
}

func TestCoreGroup(t *testing.T) {
	core := NewCore()
	g := core.Group("/subject")
	g.GET("/test", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	id := 1
	g.GET("/test/:id", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	gg := g.Group("/subsubject")
	gg.GET("/test2", func(c *Context) {
		c.Text("%s", c.Uri())
	})
	srv := &http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	time.Sleep(time.Second)
	defer srv.Shutdown(context.Background())
	resp, err := http.Get("http://127.0.0.1:8080/subject/test")
	assert.Nil(t, err)
	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/subject/test", body.String())

	resp, err = http.Get("http://127.0.0.1:8080/subject/test/" + strconv.Itoa(id))
	assert.Nil(t, err)
	body = bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "/subject/test/"+strconv.Itoa(id), body.String())

	resp, err = http.Get("http://127.0.0.1:8080/subject/subsubject/test2")
	assert.Nil(t, err)
	body.Reset()
	io.Copy(body, resp.Body)
	assert.Equal(t, "/subject/subsubject/test2", body.String())
}
