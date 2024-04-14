package server_tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestGetManyBanners200_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with tag only)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array()
}

func TestGetManyBanners200_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with tag, limit)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)
}

func TestGetManyBanners200_Test_3(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with tag, limit, offset)",
	})

	exp.GET("/banner").
		WithQuery("tag_id", 8).
		WithQuery("offset", 2).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(2)
}

func TestGetManyBanners200_Test_4(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with feature only)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array()
}

func TestGetManyBanners200_Test_5(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with feature, limit)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(1)
}

func TestGetManyBanners200_Test_6(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 200 (with feature, tag)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithQuery("tag_id", 2).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Array().Length().IsEqual(1)
}

func TestGetManyBanners400_Test_4(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 400 (wrong_token)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "wrong_token").
		Expect().Status(http.StatusUnauthorized)
}

func TestGetManyBanners403_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /banner, status 403 (user_token)",
	})

	exp.GET("/banner").
		WithQuery("feature_id", 8).
		WithQuery("limit", 2).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusForbidden)
}
