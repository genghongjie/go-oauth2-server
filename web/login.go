package web

import (
	"net/http"

	"github.com/genghongjie/go-oauth2-server/session"
)

func (s *Service) loginForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "login.html", map[string]interface{}{
		"error":       errMsg,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) login(w http.ResponseWriter, r *http.Request) {

	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the client from the request context
	client, err := getClient(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate the user
	user, err := s.oauthService.AuthUser(
		r.Form.Get("email"),    // username
		r.Form.Get("password"), // password
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Get the scope string
	scope, err := s.oauthService.GetScope(r.Form.Get("scope"))
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}
	// Log in the user
	accessToken, refreshToken, err := s.oauthService.Login(
		client,
		user,
		scope,
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}
	// Log in the user and store the user session in a cookie
	userSession := &session.UserSession{
		ClientID:     client.Key,
		Username:     user.Username,
		AccessToken:  accessToken.Token,
		RefreshToken: refreshToken.Token,
	}
	if err := sessionService.SetUserSession(userSession); err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}
	// Redirect to the authorize page by default but allow redirection to other
	// pages by specifying a path with login_redirect_uri query string param
	loginRedirectURI := r.URL.Query().Get("login_redirect_uri")
	if loginRedirectURI == "" {
		loginRedirectURI = "/web/admin"
	}
	// loginRedirectURI = "http://www.baidu.com"
	//redirectWithQueryString(loginRedirectURI, r.URL.Query(), w, r)

	_, client, user, responseType, redirectURI, err := s.authorizeCommon(r)
	query := redirectURI.Query()
	authorizationCode, err := s.oauthService.GrantAuthorizationCode(
		client,                       // client
		user,                         // user
		s.cnf.Oauth.AuthCodeLifetime, // expires in
		redirectURI.String(),         // redirect URI
		scope,                        // scope
	)
	if err != nil {
		errorRedirect(w, r, redirectURI, "server_error", "", responseType)
		return
	}

	// Set query string params for the redirection URL
	query.Set("code", authorizationCode.Code)

	// And we're done here, redirect
	redirectWithQueryString(redirectURI.String(), query, w, r)
}
