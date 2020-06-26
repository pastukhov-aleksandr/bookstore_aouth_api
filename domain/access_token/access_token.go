package access_token

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pastukhov-aleksandr/bookstore_utils-go/rest_errors"
	uuid "github.com/satori/go.uuid"
)

const (
	expirationTime             = 15
	refreshExpirationTime      = 60 * 24 * 7
	grantTypePassword          = "password"
	grandTypeClientCredentials = "client_credentials"
)

type AccessTokenRequest struct {
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
	UuID     string `json:"uuid"`
}

func (at *AccessTokenRequest) Validate() rest_errors.RestErr {
	// switch at.GrantType {
	// case grantTypePassword:
	// 	break

	// case grandTypeClientCredentials:
	// 	break

	// default:
	// 	return rest_errors.NewBadRequestError("invalid grant_type parameter")
	// }

	//TODO: Validate parameters for each grant_type
	return nil
}

type AccessToken struct {
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
	UserID         int64  `json:"user_id"`
	ClientID       int64  `json:"client_id,omitempty"`
	Expires        int64  `json:"expires"`
	ExpiresRefresh int64  `json:"expires_refresh"`
	AccessUuID     string `json:"access_uuid"`
	Permission     string `json:"permission"`
}

func (at *AccessToken) Validate() rest_errors.RestErr {
	at.AccessToken = strings.TrimSpace(at.AccessToken)
	if at.AccessToken == "" {
		return rest_errors.NewBadRequestError("invalid access token id")
	}
	if at.UserID <= 0 {
		return rest_errors.NewBadRequestError("invalid user id")
	}
	if at.ClientID <= 0 {
		return rest_errors.NewBadRequestError("invalid client id")
	}
	if at.Expires <= 0 {
		return rest_errors.NewBadRequestError("invalid expiration time")
	}
	if at.ExpiresRefresh <= 0 {
		return rest_errors.NewBadRequestError("invalid refresh expiration time")
	}
	if len(at.AccessUuID) <= 0 {
		return rest_errors.NewBadRequestError("invalid access uuid")
	}
	return nil
}

func GetNewAccessToken(userId int64) AccessToken {
	return AccessToken{
		UserID:         userId,
		Expires:        time.Now().UTC().Add(expirationTime * time.Minute).Unix(),
		ExpiresRefresh: time.Now().UTC().Add(refreshExpirationTime * time.Minute).Unix(),
		AccessUuID:     uuid.NewV4().String(),
		Permission:     "aouth",
	}
}

func (at AccessToken) IsExpired() bool {
	return time.Unix(at.Expires, 0).Before(time.Now().UTC())
}

func (at *AccessToken) Generate(access_sicret string, refresh_sicret string) rest_errors.RestErr {
	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = at.UserID
	atClaims["client_id"] = at.ClientID
	atClaims["access_uuid"] = at.AccessUuID
	atClaims["exp"] = at.Expires

	atJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	at.AccessToken, err = atJwt.SignedString([]byte(access_sicret))
	if err != nil {
		return rest_errors.NewBadRequestError("invalid JWT")
	}

	atClaims["exp"] = at.ExpiresRefresh
	atJwt = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	at.RefreshToken, err = atJwt.SignedString([]byte(refresh_sicret))
	if err != nil {
		return rest_errors.NewBadRequestError("invalid JWT")
	}

	return nil
}
