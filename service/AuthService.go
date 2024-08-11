package service

import (
	"HighArch-dialogs/api/private"
	"encoding/json"
	"io"
	"net/http"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Authenticate(token string, xRequestId string) (userId *string, err error) {
	// TODO move logic of http request in separate class
	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/internal/checkAuth/"+token, nil)
	req.Header.Set("X-Request-Id", xRequestId)
	req.Header.Set("Authorization", "internalApiAuthToken")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var authApiModel private.CheckAuthSuccessApiModel
	json.Unmarshal(body, &authApiModel)
	return &authApiModel.UserId, nil
}

const internalApiAuthToken = "X3sF9iQvQb9Q2JLHjd55ovISTk7gWLzp"
