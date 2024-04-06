package openapi

import (
	"banner/internal/cashe"
	"banner/internal/config"
	"banner/internal/postgresql"
	"banner/models"
	"context"
	//"database/sql"
	"errors"
	"net/http"
	"time"
)

// DefaultAPIService is a service that implements the logic for the DefaultAPIServicer
// This service should implement the business logic for every endpoint for the DefaultAPI API.
// Include any external packages or services that will be required by this service.
type DefaultAPIService struct {
	c  *cashe.Cache
	db *postgresql.Postgres
}

// NewDefaultAPIService creates a default api service
func NewDefaultAPIService() DefaultAPIServicer {
	cfg := config.MustLoad()
	database := postgresql.NewPostgresRepository(cfg.DbConfig)
	return &DefaultAPIService{
		c:  cashe.New(5*time.Minute, 5*time.Minute+30*time.Second),
		db: database,
	}
}

// BannerGet - Получение всех баннеров c фильтрацией по фиче и/или тегу
func (s *DefaultAPIService) BannerGet(ctx context.Context, token string, featureId int32, tagId int32, limit int32, offset int32) (ImplResponse, error) {
	// TODO - update BannerGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []BannerGet200ResponseInner{}) or use other options such as http.Ok ...
	// return Response(200, []BannerGet200ResponseInner{}), nil

	// TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	// return Response(401, nil),nil

	// TODO: Uncomment the next line to return response Response(403, {}) or use other options such as http.Ok ...
	// return Response(403, nil),nil

	// TODO: Uncomment the next line to return response Response(500, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(500, UserBannerGet400Response{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("BannerGet method not implemented")
}

// BannerIdDelete - Удаление баннера по идентификатору
func (s *DefaultAPIService) BannerIdDelete(ctx context.Context, id int32, token string) (ImplResponse, error) {
	// TODO - update BannerIdDelete with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(400, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(400, UserBannerGet400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	// return Response(401, nil),nil

	// TODO: Uncomment the next line to return response Response(403, {}) or use other options such as http.Ok ...
	// return Response(403, nil),nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(500, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(500, UserBannerGet400Response{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("BannerIdDelete method not implemented")
}

// BannerIdPatch - Обновление содержимого баннера
func (s *DefaultAPIService) BannerIdPatch(ctx context.Context, id int32, bannerIdDeleteRequest models.BannerIdDeleteRequest, token string) (ImplResponse, error) {
	// TODO - update BannerIdPatch with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	// return Response(200, nil),nil

	// TODO: Uncomment the next line to return response Response(400, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(400, UserBannerGet400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	// return Response(401, nil),nil

	// TODO: Uncomment the next line to return response Response(403, {}) or use other options such as http.Ok ...
	// return Response(403, nil),nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(500, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(500, UserBannerGet400Response{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("BannerIdPatch method not implemented")
}

// BannerPost - Создание нового баннера
func (s *DefaultAPIService) BannerPost(ctx context.Context, bannerGetRequest models.BannerGetRequest, token string) (ImplResponse, error) {
	// TODO - update BannerPost with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(201, BannerGet201Response{}) or use other options such as http.Ok ...
	// return Response(201, BannerGet201Response{}), nil

	// TODO: Uncomment the next line to return response Response(400, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(400, UserBannerGet400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	// return Response(401, nil),nil

	// TODO: Uncomment the next line to return response Response(403, {}) or use other options such as http.Ok ...
	// return Response(403, nil),nil

	// TODO: Uncomment the next line to return response Response(500, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(500, UserBannerGet400Response{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("BannerPost method not implemented")
}

// UserBannerGet - Получение баннера для пользователя
func (s *DefaultAPIService) UserBannerGet(ctx context.Context, tagId int32, featureId int32, useLastRevision bool, token string) (ImplResponse, error) {
	// TODO - update UserBannerGet with the required logic for this service method.
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, map[string]interface{}{}) or use other options such as http.Ok ...
	// return Response(200, map[string]interface{}{}), nil

	// TODO: Uncomment the next line to return response Response(400, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(400, UserBannerGet400Response{}), nil

	// TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	// return Response(401, nil),nil

	// TODO: Uncomment the next line to return response Response(403, {}) or use other options such as http.Ok ...
	// return Response(403, nil),nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(500, UserBannerGet400Response{}) or use other options such as http.Ok ...
	// return Response(500, UserBannerGet400Response{}), nil

	return Response(http.StatusNotImplemented, nil), errors.New("UserBannerGet method not implemented")
}
