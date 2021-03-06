package web

import (
	"errors"
	"net/http"

	"github.com/genghongjie/go-oauth2-server/config"
	"github.com/genghongjie/go-oauth2-server/oauth"
	"github.com/genghongjie/go-oauth2-server/session"
)

var (
	// ErrInvalidGrantType ...
	ErrInvalidGrantType = errors.New("Invalid grant type")
	// ErrInvalidClientIDOrSecret ...
	ErrInvalidClientIDOrSecret = errors.New("Invalid client ID or secret")
)

// Service struct keeps variables for reuse
type Service struct {
	cnf            *config.Config
	oauthService   oauth.ServiceInterface
	sessionService session.ServiceInterface
}

// NewService returns a new Service instance
func NewService(cnf *config.Config, oauthService oauth.ServiceInterface, sessionService session.ServiceInterface) *Service {
	return &Service{
		cnf:            cnf,
		oauthService:   oauthService,
		sessionService: sessionService,
	}
}

// GetConfig returns config.Config instance
func (s *Service) GetConfig() *config.Config {
	return s.cnf
}

// GetOauthService returns oauth.Service instance
func (s *Service) GetOauthService() oauth.ServiceInterface {
	return s.oauthService
}

// GetSessionService returns session.Service instance
func (s *Service) GetSessionService() session.ServiceInterface {
	return s.sessionService
}

// Close stops any running services
func (s *Service) Close() {}

func (s *Service) setSessionService(r *http.Request, w http.ResponseWriter) {
	s.sessionService.SetSessionService(r, w)
}
