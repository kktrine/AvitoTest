package server_tests

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

//feature = [999] : ["is_active"] = true
//feature = [1000]: ["is_active"] = false

func TestGerUserBanner200_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 1, status 200",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 999).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusOK).JSON().Raw()
}

func TestGerUserBanner200_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 2, status 200 (from cache)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 999).
		WithQuery("use_last_revision", false).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusOK).JSON().Raw()
}

func TestGerUserBanner200_Test_3(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 3, status 200 (admin token, is_active = true)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 999).
		WithQuery("use_last_revision", true).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Raw()
}

func TestGerUserBanner200_Test_4(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 4, status 200 (admin token, is_active = false)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 1000).
		WithQuery("use_last_revision", true).
		WithHeader("token", "admin_token").
		Expect().Status(http.StatusOK).JSON().Raw()

}

func TestGerUserBanner400_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 4, status 400 (invalid feature and tag)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 0).
		WithQuery("feature_id", 0).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusBadRequest).JSON().IsEqual("Некорректные данные. Фича и тэг должны быть положительными числами")
}

func TestGerUserBanner400_Test_2(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 2, status 400 (invalid feature)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 0).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusBadRequest).JSON().IsEqual("Некорректные данные. Фича и тэг должны быть положительными числами")
}

func TestGerUserBanner400_Test_3(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 3, status 400 (invalid tag)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 0).
		WithQuery("feature_id", 1).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusBadRequest).JSON().IsEqual("Некорректные данные. Фича и тэг должны быть положительными числами")
}

func TestGerUserBanner401_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 1, status 401 (invalid token)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 1).
		WithQuery("use_last_revision", true).
		WithHeader("token", "wrong_token").
		Expect().Status(http.StatusUnauthorized)
}

func TestGerUserBanner403_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 1, status 403 (user_token to is_active = false banner)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 1000).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusForbidden).JSON().IsEqual("Пользователь не имеет доступа")
}

func TestGerUserBanner404_Test_1(t *testing.T) {
	exp := httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  "http://localhost:8080",
		Reporter: httpexpect.NewAssertReporter(t),
		TestName: "GET /user_banner 1, status 404 (banner not found)",
	})

	exp.GET("/user_banner").
		WithQuery("tag_id", 1).
		WithQuery("feature_id", 1999).
		WithQuery("use_last_revision", true).
		WithHeader("token", "user_token").
		Expect().Status(http.StatusNotFound).JSON().IsEqual("Баннер не найден")
}
