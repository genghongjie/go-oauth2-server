package oauth

import (
	"github.com/genghongjie/go-oauth2-server/util/routes"
	"github.com/gorilla/mux"
)

const (
	tokensResource     = "tokens"
	tokensPath         = "/" + tokensResource
	introspectResource = "introspect"
	introspectPath     = "/" + introspectResource
	userResource       = "user"
	userPath           = "/" + userResource
)

// RegisterRoutes registers route handlers for the oauth service
func (s *Service) RegisterRoutes(router *mux.Router, prefix string) {
	subRouter := router.PathPrefix(prefix).Subrouter()
	routes.AddRoutes(s.GetRoutes(), subRouter)
}

// GetRoutes returns []routes.Route slice for the oauth service
func (s *Service) GetRoutes() []routes.Route {
	return []routes.Route{
		{
			Name:        "oauth_tokens",
			Method:      "GET",
			Pattern:     "/token",
			HandlerFunc: s.tokensHandler,
		},
		{
			Name:        "oauth_introspect",
			Method:      "GET",
			Pattern:     "/userinfo",
			HandlerFunc: s.introspectHandler,
		},
	}

}
