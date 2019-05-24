package oauth

import (
	"github.com/genghongjie/go-oauth2-server/models"
)

// AccessTokenResponse ...
type AccessTokenResponse struct {
	UserID       string `json:"user_id,omitempty"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// IntrospectResponse ...
type IntrospectResponse struct {
	Active       bool   `json:"active"`
	Scope        string `json:"scope,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	Username     string `json:"username,omitempty"`
	RoleId       string `json:"role_id,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresAt    int    `json:"exp,omitempty"`
	SaicUserInfo *SaicUserInfo
}
type SaicUserInfo struct {
	Id       string `json:"id"`
	SaicId   string `json:"saicId"`
	UserName string `json:"userName"`
	PinYin   string `json:"pinYin"`
	UserType string `json:"userType"`
	Sex      string `json:"sex"`
	AdAccout string `json:"adAccout"`
}

// NewAccessTokenResponse ...
func NewAccessTokenResponse(accessToken *models.OauthAccessToken, refreshToken *models.OauthRefreshToken, lifetime int, theTokenType string) (*AccessTokenResponse, error) {
	response := &AccessTokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   lifetime,
		TokenType:   theTokenType,
		Scope:       accessToken.Scope,
	}
	if accessToken.UserID.Valid {
		response.UserID = accessToken.UserID.String
	}
	if refreshToken != nil {
		response.RefreshToken = refreshToken.Token
	}
	return response, nil
}
