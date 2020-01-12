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

func TestAvailable(t *testing.T) {
	rate := New(2)
	c := 0
	for i := 0; i < 10; i++ {
		if rate.Available() {
			c = c + 1
		}
		time.Sleep(110 * time.Millisecond)
	}

	if c != 2 {
		t.Errorf("exepect %d actual %d", 2, c)
	}
}

func TestTick(t *testing.T) {
	rate := New(10)
	testingStart := time.Now()
	for i := 0; i < 10; i++ {
		rate.PromiseTick()
	}

	if time.Now().Sub(testingStart).Round(time.Second) != 1*time.Second {
		t.Errorf("expect to execute this in one second")
	}
}

func TestTimeoutMiddleware(t *testing.T) {
	pl := "hello"
	ts := httptest.NewServer(New(1).TimeoutMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, pl)
	}), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})))

	time.Sleep(1 * time.Second)

	res1, err := http.Get(ts.URL)
	defer res1.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}

	res2, err := http.Get(ts.URL)
	defer res2.Body.Close()
	if err != nil {
		t.Error(err.Error())
	}

	b, err := ioutil.ReadAll(res1.Body)

	if !strings.EqualFold(pl, string(b)) {
		t.Errorf("expect %s actual %s", pl, string(b))
	}

	if res2.StatusCode != http.StatusNotFound {
		t.Errorf("expect %d actual %d", http.StatusNotFound, res2.StatusCode)
	}
}

func TestPromiseMiddleware(t *testing.T) {
	pl := "hello"
	ts := httptest.NewServer(New(1).PromiseMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, pl)
	})))

	defer ts.Close()

	start := time.Now()
	res, err := http.Get(ts.URL)
	if time.Now().Sub(start).Round(time.Second) != 1*time.Second {
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
