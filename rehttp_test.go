package rehttp

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoDelay(t *testing.T) {
	fn := NoDelay()
	want := time.Duration(0)
	for i := 0; i < 5; i++ {
		delay := fn(nil, nil, i, nil)
		assert.Equal(t, want, delay, "%d", i)
	}
}

func TestConstDelay(t *testing.T) {
	want := 2 * time.Second
	fn := ConstDelay(want)
	for i := 0; i < 5; i++ {
		delay := fn(nil, nil, i, nil)
		assert.Equal(t, want, delay, "%d", i)
	}
}

func TestLinearDelay(t *testing.T) {
	initial := 2 * time.Second
	fn := LinearDelay(initial)
	want := []time.Duration{2 * time.Second, 4 * time.Second, 6 * time.Second, 8 * time.Second, 10 * time.Second}
	for i := 0; i < len(want); i++ {
		got := fn(nil, nil, i, nil)
		assert.Equal(t, want[i], got, "%d", i)
	}
}

func TestExponentialDelay(t *testing.T) {
	initial := 2 * time.Second
	fn := ExponentialDelay(initial, time.Second)
	want := []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second, 16 * time.Second, 32 * time.Second}
	for i := 0; i < len(want); i++ {
		got := fn(nil, nil, i, nil)
		assert.Equal(t, want[i], got, "%d", i)
	}

	initial = 100 * time.Millisecond
	fn = ExponentialDelay(initial, 10*time.Millisecond)
	want = []time.Duration{100 * time.Millisecond, time.Second, 10 * time.Second}
	for i := 0; i < len(want); i++ {
		got := fn(nil, nil, i, nil)
		assert.Equal(t, want[i], got, "%d", i)
	}
}

func TestRetryHTTPMethods(t *testing.T) {
	cases := []struct {
		retries int
		meths   []string
		inMeth  string
		att     int
		want    bool
	}{
		{retries: 1, meths: nil, inMeth: "GET", att: 0, want: false},
		{retries: 0, meths: nil, inMeth: "GET", att: 1, want: false},
		{retries: 1, meths: []string{"get"}, inMeth: "GET", att: 0, want: true},
		{retries: 1, meths: []string{"GET"}, inMeth: "GET", att: 0, want: true},
		{retries: 1, meths: []string{"GET"}, inMeth: "POST", att: 0, want: false},
		{retries: 2, meths: []string{"GET", "POST"}, inMeth: "POST", att: 0, want: true},
		{retries: 2, meths: []string{"GET", "POST"}, inMeth: "POST", att: 1, want: true},
		{retries: 2, meths: []string{"GET", "POST"}, inMeth: "POST", att: 2, want: false},
		{retries: 2, meths: []string{"GET", "POST"}, inMeth: "put", att: 0, want: false},
		{retries: 2, meths: []string{"GET", "POST", "PUT"}, inMeth: "put", att: 0, want: true},
	}

	for i, tc := range cases {
		fn := RetryHTTPMethods(tc.retries, tc.meths...)
		req, err := http.NewRequest(tc.inMeth, "", nil)
		require.Nil(t, err)
		got := fn(req, nil, tc.att, nil)
		assert.Equal(t, tc.want, got, "%d", i)
	}
}

func TestRetryStatus500(t *testing.T) {

}
