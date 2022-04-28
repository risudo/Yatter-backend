package parameters_test

import (
	"math"
	"net/http"
	"testing"
	"yatter-backend-go/app/domain/object"
	"yatter-backend-go/app/handler/handler_test_setup"
	"yatter-backend-go/app/handler/parameters"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParseErr(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name          string
		actual_is_err bool
		expect_is_err bool
	}{
		{
			name: "string in max_id",
			actual_is_err: func() bool {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("max_id", "a")
				req.URL.RawQuery = params.Encode()
				_, err = parameters.ParseAll(req)
				return err != nil
			}(),
			expect_is_err: true,
		},
		{
			name: "overflow max_id",
			actual_is_err: func() bool {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("max_id", "9999999999999999999")
				req.URL.RawQuery = params.Encode()
				_, err = parameters.ParseAll(req)
				return err != nil
			}(),
			expect_is_err: true,
		},
		{
			name: "overflow since_id",
			actual_is_err: func() bool {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("since_id", "9999999999999999999")
				req.URL.RawQuery = params.Encode()
				_, err = parameters.ParseAll(req)
				return err != nil
			}(),
			expect_is_err: true,
		},
		{
			name: "OverMaxLimit",
			actual_is_err: func() bool {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "81")
				req.URL.RawQuery = params.Encode()
				_, err = parameters.ParseAll(req)
				return err != nil
			}(),
			expect_is_err: true,
		},
		{
			name: "LessMinLimit",
			actual_is_err: func() bool {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "-1")
				req.URL.RawQuery = params.Encode()
				_, err = parameters.ParseAll(req)
				return err != nil
			}(),
			expect_is_err: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !assert.Equal(t, tt.expect_is_err, tt.actual_is_err) {
				return
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	m := handler_test_setup.MockSetup()
	defer m.Close()

	tests := []struct {
		name   string
		actual *object.Parameters
		expect *object.Parameters
	}{
		{
			name: "only_media",
			actual: func() *object.Parameters {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("only_media", "1")
				req.URL.RawQuery = params.Encode()
				param, err := parameters.ParseAll(req)
				return param
			}(),
			expect: &object.Parameters{
				OnlyMedia: true,
				MaxID:     math.MaxInt64,
				SinceID:   0,
				Limit:     parameters.DefaultLimit,
			},
		},
		{
			name: "max_id",
			actual: func() *object.Parameters {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("max_id", "90")
				req.URL.RawQuery = params.Encode()
				param, err := parameters.ParseAll(req)
				return param
			}(),
			expect: &object.Parameters{
				OnlyMedia: false,
				MaxID:     90,
				SinceID:   0,
				Limit:     parameters.DefaultLimit,
			},
		},
		{
			name: "since_id",
			actual: func() *object.Parameters {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("since_id", "10")
				req.URL.RawQuery = params.Encode()
				param, err := parameters.ParseAll(req)
				return param
			}(),
			expect: &object.Parameters{
				OnlyMedia: false,
				MaxID:     math.MaxInt64,
				SinceID:   10,
				Limit:     parameters.DefaultLimit,
			},
		},
		{
			name: "limit",
			actual: func() *object.Parameters {
				req, err := http.NewRequest("GET", m.AsURL("/v1/timelines/public"), nil)
				if err != nil {
					t.Fatal(err)
				}
				params := req.URL.Query()
				params.Add("limit", "10")
				req.URL.RawQuery = params.Encode()
				param, err := parameters.ParseAll(req)
				return param
			}(),
			expect: &object.Parameters{
				OnlyMedia: false,
				MaxID:     math.MaxInt64,
				SinceID:   0,
				Limit:     10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if d := cmp.Diff(tt.actual, tt.expect); len(d) != 0 {
				t.Fatalf("differs: (-got +want)\n%s", d)
			}
		})
	}
}
