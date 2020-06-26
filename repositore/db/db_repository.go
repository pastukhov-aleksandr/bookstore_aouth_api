package db

import (
	"errors"
	"time"

	"github.com/gocql/gocql"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/clients/cassandra"
	"github.com/pastukhov-aleksandr/bookstore_aouth_api/domain/access_token"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessToken    = "SELECT refresh_tokens, user_id, client_id, expires FROM refresh_tokens WHERE refresh_tokens=?;"
	queryCreateAccessToken = "INSERT INTO refresh_tokens(uuid, user_id, client_id, now) VALUES (?, ?, ?, ?);"
	queryUpdateExpires     = "UPDATE refresh_tokens SET expires=? WHERE refresh_tokens=?;"
)

func NewRepository() DbRepository {
	return &dbRepository{}
}

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, rest_errors.RestErr)
	Create(access_token.AccessToken) rest_errors.RestErr
	UpdateExpirationTime(access_token.AccessToken) rest_errors.RestErr
}

type dbRepository struct {
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, rest_errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessToken, id).Scan(
		&result.AccessToken,
		&result.UserID,
		&result.ClientID,
		&result.Expires,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, rest_errors.NewNotFoundError("no access token found with given id")
		}
		return nil, rest_errors.NewInternalServerError("error when trying to get current id", errors.New("database error"))
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessUuID,
		at.UserID,
		at.ClientID,
		time.Now(),
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying to save refresh token in database", err)
	}
	return nil
}

func (r *dbRepository) UpdateExpirationTime(at access_token.AccessToken) rest_errors.RestErr {
	if err := cassandra.GetSession().Query(queryUpdateExpires,
		at.Expires,
		at.AccessToken,
	).Exec(); err != nil {
		return rest_errors.NewInternalServerError("error when trying to update current resource", errors.New("database error"))
	}
	return nil
}
