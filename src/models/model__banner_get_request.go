/*
 * Сервис баннеров
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type BannerGetRequest struct {

	// Идентификаторы тэгов
	TagIds []int32 `json:"tag_ids,omitempty"`

	// Идентификатор фичи
	FeatureId int32 `json:"feature_id,omitempty"`

	// Содержимое баннера
	Content map[string]interface{} `json:"content,omitempty"`

	// Флаг активности баннера
	IsActive bool `json:"is_active,omitempty"`
}

// AssertBannerGetRequestRequired checks if the required fields are not zero-ed
func AssertBannerGetRequestRequired(obj BannerGetRequest) error {
	return nil
}

// AssertBannerGetRequestConstraints checks if the values respects the defined constraints
func AssertBannerGetRequestConstraints(obj BannerGetRequest) error {
	return nil
}