package parameters_test

import (
	"net/http"
	"testing"
	"yatter-backend-go/app/handler/handler_test_setup"
	"yatter-backend-go/app/handler/parameters"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
				_, err = parameters.Parse(req)
				if err != nil {
					return true
				}
				return false
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
				_, err = parameters.Parse(req)
				if err != nil {
					return true
				}
				return false
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
				_, err = parameters.Parse(req)
				if err != nil {
					return true
				}
				return false
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
