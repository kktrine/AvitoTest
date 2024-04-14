package server_tests

import (
	"banner/models"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

func TestPatch200_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "PATCH /banner/{id}, status 200",
	})
	var featureID int32 = 1
	exp.PATCH("/banner/{id}").WithPath("id", 2).
		WithJSON(models.BannerIdDeleteRequest{
			TagIds:    &[]int32{54, 85},
			FeatureId: &featureID,
		}).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK)
}

func TestPatch400_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "PATCH /banner/{id}, status 400",
	})
	var featureID int32 = 2001
	exp.PATCH("/banner/{id}").WithPath("id", 55000).
		WithJSON(models.BannerIdDeleteRequest{
			TagIds:    &[]int32{0, 85},
			FeatureId: &featureID,
		}).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusBadRequest)
}

func TestPatch401_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "PATCH /banner/{id}, status 401 (wrong_token)",
	})
	var featureID int32 = 2001
	exp.PATCH("/banner/{id}").WithPath("id", 55000).
		WithJSON(models.BannerIdDeleteRequest{
			TagIds:    &[]int32{54, 85},
			FeatureId: &featureID,
		}).
		WithHeader("token", "wrong_token").
		Expect().Status(http.StatusUnauthorized)
}

func TestPatch403_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "PATCH /banner/{id}, status 403 (user_token)",
	})
	var featureID int32 = 2001
	exp.PATCH("/banner/{id}").WithPath("id", 55000).
		WithJSON(models.BannerIdDeleteRequest{
			TagIds:    &[]int32{54, 85},
			FeatureId: &featureID,
		}).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusForbidden)
}

func TestPatch404_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "PATCH /banner/{id}, status 404 invalid id",
	})
	var featureID int32 = 2001
	exp.PATCH("/banner/{id}").WithPath("id", 11155000).
		WithJSON(models.BannerIdDeleteRequest{
			TagIds:    &[]int32{54, 85},
			FeatureId: &featureID,
		}).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusNotFound)
}
