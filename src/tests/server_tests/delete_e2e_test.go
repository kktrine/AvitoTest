package server_tests

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func TestDelete200_Test_1(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "DELETE /banner/{id}, status 204",
	})
	e.DELETE("/banner/{id}").WithPath("id", 1).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusNoContent)
}

func TestDelete400_Test_1(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "DELETE /banner/{id}, status 400",
	})
	e.DELETE("/banner/{id}").WithPath("id", -1).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusBadRequest)
}

func TestDelete401_Test_1(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "DELETE /banner/{id}, status 401 (wrong_token)",
	})
	e.DELETE("/banner/{id}").WithPath("id", 1).
		WithHeader("token", "wrong_token").
		Expect().Status(http.StatusUnauthorized)
}

func TestDelete403_Test_1(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "DELETE /banner/{id}, status 401 (user_token)",
	})
	e.DELETE("/banner/{id}").WithPath("id", 1).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusForbidden)
}

func TestDelete404_Test_1(t *testing.T) {
	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "DELETE /banner/{id}, status 404 ",
	})
	e.DELETE("/banner/{id}").WithPath("id", 100000000).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusNotFound)
}
