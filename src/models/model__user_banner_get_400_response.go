package models

type UserBannerGet400Response struct {
	Error string `json:"error,omitempty"`
}

// AssertUserBannerGet400ResponseRequired checks if the required fields are not zero-ed
func AssertUserBannerGet400ResponseRequired(obj UserBannerGet400Response) error {
	return nil
}

// AssertUserBannerGet400ResponseConstraints checks if the values respects the defined constraints
func AssertUserBannerGet400ResponseConstraints(obj UserBannerGet400Response) error {
	return nil
}
