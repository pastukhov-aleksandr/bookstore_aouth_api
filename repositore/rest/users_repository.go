package rest

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/domain/users"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, rest_errors.RestErr)
}

type usersRepository struct{}

func NewRestUsersRepository() RestUsersRepository {
	return &usersRepository{}
}

type AuthSuccess struct {
	/* variables */
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}

	// Create a Resty Client
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(&AuthSuccess{}).
		Post("http://localhost:8081/api/users/login")

	if err != nil || response.Body() == nil {
		return nil, rest_errors.NewInternalServerError("invalid restclient response when trying to login user", errors.New("restclient error"))
	}

	if response.StatusCode() > 299 {
		apiErr, err := rest_errors.NewRestErrorFromBytes(response.Body())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("invalid error interface when trying to login user", err)
		}
		return nil, apiErr
	}

	var user users.User
	if err := json.Unmarshal(response.Body(), &user); err != nil {
		return nil, rest_errors.NewInternalServerError("error when trying to unmarshal users login response", errors.New("json parsing error"))
	}
	return &user, nil
}
