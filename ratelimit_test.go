package rate_limit

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTick(t *testing.T) {
	rate := New(10)
	testingStart := time.Now()
	for i:=0;i < 10; i++ {
		rate.Tick()
	}

	if time.Now().Sub(testingStart).Round(time.Second) != 1 * time.Second {
		t.Errorf("expect to execute this in one second")
	}
}

func TestMiddleware(t *testing.T){
	pl := "hello"
	ts := httptest.NewServer(New(1).Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, pl)
	})))

	defer ts.Close()

	start := time.Now()
	res, err := http.Get(ts.URL)
	if time.Now().Sub(start).Round(time.Second) != 1 *time.Second {
		t.Errorf("expect to execute this in one second")
	}

	if err != nil {
		t.Errorf(err.Error())
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !strings.EqualFold(pl, string(b)) {
		t.Errorf("expect %s and %s are equal", pl, string(b))
	}
}