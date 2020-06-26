package access_token

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

const (
	ACCESS_SECRET  = "ACCESS_SECRET"
	REFRESH_SECRET = "REFRESH_SECRET"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type service struct {
	restUsersRepo rest.RestUsersRepository
	dbRepo        db.DbRepository
}

func NewService(usersRepo rest.RestUsersRepository, dbRepo db.DbRepository) Service {
	return &service{
		restUsersRepo: usersRepo,
		dbRepo:        dbRepo,
	}
}

func (s *service) GetById(accessTokenId string) (*access_token.AccessToken, rest_errors.RestErr) {
	accessTokenId = strings.TrimSpace(accessTokenId)
	if len(accessTokenId) == 0 {
		return nil, rest_errors.NewBadRequestError("invalid access token id")
	}
	accessToken, err := s.dbRepo.GetById(accessTokenId)
	if err != nil {
		return nil, err
	}
	return accessToken, nil
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	//TODO: Support both grant types: client_credentials and password

	// Если нужна авторизация:
	// user, err := s.restUsersRepo.LoginUser(request.Username, request.Password)
	// if err != nil {
	// 	return nil, err
	// }

	// Generate a new access token:
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
		return nil, rest_errors.NewBadRequestError("invalid access token")
	}

	at := access_token.GetNewAccessToken(request.UserID)
	at.ClientID = request.ClientID
	if len(request.UuID) > 0 {
		at.AccessUuID = request.UuID
	}

	var access_sicret = getEnv(ACCESS_SECRET, "")
	var refresh_sicret = getEnv(REFRESH_SECRET, "")
	if err := at.Generate(access_sicret, refresh_sicret); err != nil {
		return nil, err
	}

	// Save the new refresh access token in Cassandra:
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}
	return &at, nil
}

func (s *service) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := at.Validate(); err != nil {
		return err
	}
	return s.dbRepo.UpdateExpirationTime(at)
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
