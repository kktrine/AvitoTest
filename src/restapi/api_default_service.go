package openapi

import (
	"banner/internal/simple_auth"
	"banner/internal/storage"
	"banner/models"
	"context"
	"strconv"

	//"database/sql"
	"errors"
	"net/http"
)

// DefaultAPIService is a service that implements the logic for the DefaultAPIServicer
// This service should implement the business logic for every endpoint for the DefaultAPI API.
// Include any external packages or services that will be required by this service.
type DefaultAPIService struct {
	Storage *storage.Storage
}

// NewDefaultAPIService creates a default api service
func NewDefaultAPIService() DefaultAPIServicer {
	storage := storage.NewStorage()
	storage.Fill()
	return &DefaultAPIService{
		Storage: storage,
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
	ok, err := simple_auth.CheckAdminToken(token)
	if err != nil {
		return Response(401, "Пользователь не авторизован"), nil
	}
	if id <= 0 {
		return Response(400, "Некорректные данные"), nil
	}
	if !ok {
		return Response(403, "Пользователь не имеет доступа"), nil
	}
	found, err := s.Storage.Delete(id)
	if err != nil {
		return Response(500, "Внутренняя ошибка сервера"), nil
	}
	if !found {
		return Response(404, "Баннер для тэга не найден"), nil
	}
	return Response(204, "Баннер успешно удален"), nil

}

// BannerIdPatch - Обновление содержимого баннера
func (s *DefaultAPIService) BannerIdPatch(ctx context.Context, id int32, bannerIdDeleteRequest models.BannerIdDeleteRequest, token string) (ImplResponse, error) {
	ok, err := simple_auth.CheckAdminToken(token)
	if err != nil {
		return Response(401, "Пользователь не авторизован"), nil
	}
	if !ok {
		return Response(403, "Пользователь не имеет доступа"), nil
	}
	if id <= 0 {
		return Response(400, "Некорректные данные. Id должен быть положительным числом"), nil
	}
	toUpdate := models.InsertData{}
	if bannerIdDeleteRequest.FeatureId != nil {
		toUpdate.Feature = *bannerIdDeleteRequest.FeatureId
	}
	if bannerIdDeleteRequest.TagIds != nil {
		toUpdate.TagIds = *bannerIdDeleteRequest.TagIds
	}
	if bannerIdDeleteRequest.Content != nil {
		toUpdate.Content = *bannerIdDeleteRequest.Content
	}
	if bannerIdDeleteRequest.IsActive != nil {
		toUpdate.IsActive = *bannerIdDeleteRequest.IsActive
	}
	found, err := s.Storage.Update(id, &toUpdate)
	if !found {
		return Response(404, "Баннер не найден"), nil
	}
	if err != nil {
		return Response(500, err.Error()), nil
	}
	return Response(200, nil), nil
}

// BannerPost - Создание нового баннера
func (s *DefaultAPIService) BannerPost(ctx context.Context, bannerGetRequest models.BannerGetRequest, token string) (ImplResponse, error) {
	ok, err := simple_auth.CheckAdminToken(token)
	if err != nil {
		return Response(401, "Пользователь не авторизован"), nil
	}
	if !ok {
		return Response(403, "Пользователь не имеет доступа"), nil
	}
	if bannerGetRequest.FeatureId <= 0 {
		return Response(400, "Некорректные данные. Фича и тэг должны быть положительными числами"), nil
	}
	for _, i := range bannerGetRequest.TagIds {
		if i <= 0 {
			return Response(400, "Некорректные данные. Фича и тэг должны быть положительными числами"), nil
		}
	}
	id, err := s.Storage.Insert(&models.InsertData{
		Feature:  bannerGetRequest.FeatureId,
		TagIds:   bannerGetRequest.TagIds,
		Content:  bannerGetRequest.Content,
		IsActive: bannerGetRequest.IsActive,
	})
	if err != nil {
		return Response(500, err.Error()), nil
	}
	return Response(201, "Created with id: "+strconv.Itoa(int(id))), nil
}

// UserBannerGet - Получение баннера для пользователя
func (s *DefaultAPIService) UserBannerGet(ctx context.Context, tagId int32, featureId int32, useLastRevision bool, token string) (ImplResponse, error) {
	// Add api_default_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	ok, err := simple_auth.CheckAdminToken(token)
	if err != nil {
		return Response(401, "Пользователь не авторизован"), nil
	}
	if tagId <= 0 || featureId <= 0 {
		return Response(400, "Некорректные данные. Фича и тэг должны быть положительными числами"), nil
	}
	res, userAccess, found, err := s.Storage.GetUserBanner(featureId, tagId)
	if err != nil && !found {
		return Response(500, "Внутренняя ошибка сервера: "+err.Error()), nil
	}
	if !found {
		return Response(404, "Баннер не найден"), nil
	}
	if ok {
		return Response(200, map[string]interface{}{}), nil
	}
	if !userAccess {
		return Response(403, "Пользователь не имеет доступа"), nil
	}
	return Response(200, res), nil
}
