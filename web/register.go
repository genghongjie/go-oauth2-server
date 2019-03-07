package web

import (
	"net/http"

	"github.com/genghongjie/go-oauth2-server/util/response"

	"github.com/genghongjie/go-oauth2-server/oauth/roles"
)

func (s *Service) registerForm(w http.ResponseWriter, r *http.Request) {
	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template
	errMsg, _ := sessionService.GetFlashMessage()
	renderTemplate(w, "register.html", map[string]interface{}{
		"error":       errMsg,
		"queryString": getQueryString(r.URL.Query()),
	})
}

func (s *Service) register(w http.ResponseWriter, r *http.Request) {

	// Get the session service from the request context
	sessionService, err := getSessionService(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check that the submitted email hasn't been registered already
	if s.oauthService.UserExists(r.Form.Get("email")) {
		sessionService.SetFlashMessage("Email taken")
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Create a user
	_, err = s.oauthService.CreateUser(
		roles.User,             // role ID
		r.Form.Get("email"),    // username
		r.Form.Get("password"), // password
	)
	if err != nil {
		sessionService.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	// Redirect to the login page
	redirectWithQueryString("/web/login", r.URL.Query(), w, r)
}

func (s *Service) addUser(w http.ResponseWriter, r *http.Request) {

	// Client auth
	_, err := s.oauthService.BasicAuthClient(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	if r.FormValue("userName") == "" {
		response.Error(w, "userName not found", http.StatusBadRequest)
		return
	}
	if r.FormValue("password") == "" {
		response.Error(w, "password not found", http.StatusBadRequest)
		return
	}
	// Check that the submitted email hasn't been registered already
	if s.oauthService.UserExists(r.Form.Get("userName")) {
		response.Error(w, "userName taken", http.StatusBadRequest)
		return
	}

	// Create a user
	user, err := s.oauthService.CreateUser(
		roles.User,             // role ID
		r.Form.Get("userName"), // username
		r.Form.Get("password"), // password
	)
	if err != nil {

		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response to json
	response.WriteJSON(w, user, 200)
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {

	// Client auth
	_, err := s.oauthService.BasicAuthClient(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}
	userId := r.FormValue("userId")
	if userId == "" {
		response.Error(w, "userId not found", http.StatusBadRequest)
		return
	}

	// Create a user
	err = s.oauthService.DeleteUserById(
		userId, // userId
	)
	if err != nil {
		response.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write response to json
	response.WriteJSON(w, "success", 200)
}
