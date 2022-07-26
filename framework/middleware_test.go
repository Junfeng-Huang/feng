package framework

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func middleware1(c *Context) {
	c.Text("%s", "go into middleware1\n")
	c.Next()
	c.Text("%s", "leave middleware1\n")
}

func middleware2(c *Context) {
	c.Text("%s", "go into middleware2\n")
	c.Next()
	c.Text("%s", "leave middleware2\n")
}

func middleware3(c *Context) {
	c.Text("%s", "go into middleware3\n")
	c.Next()
	c.Text("%s", "leave middleware3\n")
}

func TestMiddleware(t *testing.T) {
	core := NewCore()
	core.Use(middleware1)
	core.GET("/testMiddleware1", func(c *Context) {
	})
	core.GET("/testMiddleware23", middleware2, middleware3, func(c *Context) {

	})
	srv := http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	defer srv.Shutdown(context.Background())
	time.Sleep(time.Second)

	resp, err := http.Get("http://127.0.0.1:8080/testMiddleware1")
	assert.Nil(t, err)
	body := bytes.NewBuffer(nil)
	io.Copy(body, resp.Body)
	assert.Equal(t, "go into middleware1\nleave middleware1\n", body.String())

	resp, err = http.Get("http://127.0.0.1:8080/testMiddleware23")
	assert.Nil(t, err)
	body.Reset()
	io.Copy(body, resp.Body)
	assert.Equal(t, "go into middleware1\ngo into middleware2\ngo into middleware3\nleave middleware3\nleave middleware2\nleave middleware1\n",
		body.String())
}
