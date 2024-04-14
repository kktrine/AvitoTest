package server_tests

import (
	"banner/models"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func TestPostUserBanner200_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner 1, status 200",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{1, 2, 3, 4},
		FeatureId: 2000,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "admin_token").
		Expect().Status(http.StatusCreated).JSON().
		Object().ContainsKey("banner_id")
}

func TestPostUserBanner400_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner, status 400 (tag = 0)",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{0, 2, 3, 4},
		FeatureId: 2001,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "admin_token").
		Expect().Status(http.StatusBadRequest)
}

func TestPostUserBanner400_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner, status 400 (feature = 0)",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{1, 2, 3, 4},
		FeatureId: 0,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "admin_token").
		Expect().Status(http.StatusBadRequest)
}

func TestPostUserBanner400_Test_3(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner, status 400 (banner already exists)",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{1, 2, 3, 4},
		FeatureId: 3,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "admin_token").
		Expect().Status(http.StatusBadRequest)
}

func TestPostUserBanner401_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner, status 401 (wrong_token)",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{1, 2, 3, 4},
		FeatureId: 2,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "wrong_token").
		Expect().Status(http.StatusUnauthorized)
}

func TestPostUserBanner403_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "POST /banner, status 401 (user_token)",
	})

	exp.POST("/banner").WithJSON(models.BannerGetRequest{
		TagIds:    []int32{1, 2, 3, 4},
		FeatureId: 2,
		Content: map[string]interface{}{
			"title": "record from E2E test",
			"text":  "expect status 200",
		},
		IsActive: true,
	}).WithHeader("token", "user_token").
		Expect().Status(http.StatusForbidden)
}
