package framework

import (
	"bytes"
	"context"
	"feng/framework/container"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const serviceTestKey = "serviceTest"

var testout = bytes.NewBuffer(nil)

type serviceTest interface {
	use()
}

type service struct{}

func (s service) use() {
	fmt.Fprint(testout, "service is using\n")
}

type serviceProvider struct{}

func (s *serviceProvider) IsDefer() bool {
	return false
}

func (s *serviceProvider) Boot(container container.Container) error {
	return nil
}

func (s *serviceProvider) Register(container container.Container) container.NewInstance {
	fmt.Fprintf(testout, "%s", "register service\n")
	return func(params ...interface{}) (interface{}, error) {
		return service{}, nil
	}
}

func (s *serviceProvider) Name() string {
	return serviceTestKey
}

func (s *serviceProvider) Params(container container.Container) []interface{} {
	return nil
}

func TestContainer(t *testing.T) {
	core := NewCore()
	assert.False(t, core.IsBind(serviceTestKey))
	sp := serviceProvider{}
	err := core.Bind(&sp)
	assert.Nil(t, err)
	assert.Equal(t, "register service\n", testout.String())

	core.GET("/testService", func(c *Context) {
		s := c.MustMake(serviceTestKey).(serviceTest)
		s.use()
	})
	srv := http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go srv.ListenAndServe()
	time.Sleep(time.Second)
	_, err = http.Get("http://127.0.0.1:8080/testService")
	assert.Nil(t, err)
	assert.Equal(t, "register service\nservice is using\n", testout.String())
	srv.Shutdown(context.Background())

	core.GET("/testServiceMakeNew", func(c *Context) {
		s, _ := c.MakeNew(serviceTestKey, nil)
		q, ok := s.(serviceTest)
		assert.True(t, ok)
		q.use()
	})
	testout.Reset()
	srv = http.Server{Addr: "0.0.0.0:8080", Handler: core}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			t.Log(err)
		}
	}()
	time.Sleep(time.Second)
	_, err = http.Get("http://127.0.0.1:8080/testServiceMakeNew")
	assert.Nil(t, err)
	assert.Equal(t, "register service\nservice is using\n", testout.String())
	srv.Shutdown(context.Background())

}
