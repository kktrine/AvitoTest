package simple_auth

import "errors"

func CheckAdminToken(token string) (bool, error) {
	if token == "admin_token" {
		return true, nil
	}
	if token == "user_token" {
		return false, nil
	}
	return false, errors.New("invalid token")
}
