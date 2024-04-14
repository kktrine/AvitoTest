package server_tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestGerManyBanners200_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner 1, status 200 (with tag only)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array()
}

func TestGerManyBanners200_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner 1, status 200 (with tag, limit)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)
}

func TestGerManyBanners200_Test_3(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner 3, status 200 (with tag, limit, offset)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithQuery("offset", 2).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)
}

func TestGerManyBanners200_Test_4(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner 4, status 200 (with feature only)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array()
}

func TestGerManyBanners200_Test_5(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner 5, status 200 (with feature, limit)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)
}
