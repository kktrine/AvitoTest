package openapi

import (
	"banner/models"
	"context"
	"net/http"
)

// DefaultAPIRouter defines the required methods for binding the api requests to a responses for the DefaultAPI
// The DefaultAPIRouter implementation should parse necessary information from the http request,
// pass the data to a DefaultAPIServicer to perform the required actions, then write the service results to the http response.
type DefaultAPIRouter interface {
	BannerGet(http.ResponseWriter, *http.Request)
	BannerIdDelete(http.ResponseWriter, *http.Request)
	BannerIdPatch(http.ResponseWriter, *http.Request)
	BannerPost(http.ResponseWriter, *http.Request)
	UserBannerGet(http.ResponseWriter, *http.Request)
}

// DefaultAPIServicer defines the api actions for the DefaultAPI service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can be ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type DefaultAPIServicer interface {
	BannerGet(context.Context, string, int32, int32, int32, int32) (ImplResponse, error)
	BannerIdDelete(context.Context, int32, string) (ImplResponse, error)
	BannerIdPatch(context.Context, int32, models.BannerIdDeleteRequest, string) (ImplResponse, error)
	BannerPost(context.Context, models.BannerGetRequest, string) (ImplResponse, error)
	UserBannerGet(context.Context, int32, int32, bool, string) (ImplResponse, error)
}
