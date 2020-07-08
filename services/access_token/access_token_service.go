package access_token

import (
	"log"
	"strings"

	"github.com/pastukhov-aleksandr/bookstore_aouth_api/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/db"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/repositore/rest"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/secret_code"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, rest_errors.RestErr)
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
	DeleteRefreshToken(access_token.AccessTokenRequest) rest_errors.RestErr
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

	at := access_token.GetNewAccessToken(request.UserID, request.ClientID)

	access_sicret := secret_code.Get_ACCESS_SECRET()
	refresh_sicret := secret_code.Get_REFRESH_SECRET()

	if access_sicret == "" && refresh_sicret == "" {
		log.Println(".env file not found")
		return nil, rest_errors.NewBadRequestError("invalid access token")
	}

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

func (s *service) DeleteRefreshToken(request access_token.AccessTokenRequest) rest_errors.RestErr {

	if err := s.dbRepo.DeleteRefreshToken(request.UserID, request.ClientID); err != nil {
		return err
	}

	return nil
}
