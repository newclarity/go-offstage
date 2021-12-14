package offstage

import (
	"context"
	"fmt"
	"github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	labstack "github.com/labstack/echo/v4/middleware"
	"github.com/mikeschinkel/go-only"
	"github.com/pkg/errors"
	"github.com/newclarity/rep-go-offstage/git"
	"github.com/newclarity/rep-go-offstage/util"
	"net/http"
	"reflect"
)

const defaultServerPort = 9999

var InTesting = false

// Server provides a struct to "wire up" the API's services
type Server struct {
	Context *git.Context
	*echo.Echo
	Port      int
	Endpoints EndpointsConnector
}

// NewServer instantiates a pointer to a new Server instance
func NewServer(c *git.Context, e *echo.Echo) *Server {
	return &Server{
		Context: c,
		Echo:    e,
		Port:    defaultServerPort,
	}
}

// SetContext allows setting Server Context
func (srv *Server) SetContext(ctx *git.Context) {
	srv.Context = ctx
}

// GetLogger returns instance of echo.Logger from Echo
func (srv *Server) GetLogger() echo.Logger {
	return srv.Logger
}

// ConnectEndpoints configures the Echo server and connects interfaces
// communication between offstage, generated and the /src/ code.
func (srv *Server) ConnectEndpoints(ec EndpointsConnector) {
	for range only.Once {

		srv.validate(ec)

		// Captures the connector into a static package variable
		// that can connect offstage package with packages in the
		// service
		srv.Endpoints = ec

		// Allow
		ec.SetServer(srv)

		e := srv.Echo

		//goland:noinspection ALL
		if InTesting {

			// Hides display of Echo banner during testing
			e.HideBanner = true

			// Hides display of port number during testing
			e.HidePort = true

		} else {

			// Logs all requests
			// Maybe we will enable this for DEBUG in testing, too?
			e.Use(labstack.Logger())
		}

		////===[ THIS IS JUST A HOW-TO EXAMPLE ]========
		//// This shows how to use a custom context within middleware
		//e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		//	return func(c echo.Context) error {
		//		cc := &srv.CustomContext{c}
		//		return next(cc)
		//	}
		//})

		////===[ THIS IS JUST A HOW-TO EXAMPLE ]========
		//// This is a group for the route that begins with `/secret`
		//g := e.Group("/secret")
		//// This runs the MySecretMiddlewareFunc() middleware for any
		//// route starting with `/secret`
		//g.Use(MySecretMiddlewareFunc())

		// BearerAuthMiddlewareFunc up Bearer token-based authentication
		// CURRENTLY it only
		e.Use(srv.BearerAuthMiddlewareFunc())

		// Use our validation middleware to check all requests against the
		// OpenAPI schema.
		e.Use(middleware.OapiRequestValidator(srv.MakeSwagger()))

		// We now register our implementation above as the handler for the endpoints
		ec.RegisterHandlers(e, ec)

	}
}

// validate ensures the ServiceArgs are set, or if not panics
func (srv *Server) validate(ec EndpointsConnector) {
	var what string

	for range only.Once {
		t := reflect.TypeOf(ec)

		if t.Kind() != reflect.Ptr {
			what = "`Endpoints` as a POINTER to a struct with methods implementing API endpoints"
			break
		}
		if t.Elem().Kind() != reflect.Struct {
			what = "`Endpoints` as a pointer to a STRUCT with methods implementing API endpoints"
			break
		}
		what = ""
	}
	if what != "" {
		util.Fail(fmt.Sprintf("Error: Must specify %s when using support.ServiceArgs.",
			what))
	}
}

// StartServer the Echo server on localhost with the configured Port
func (srv *Server) StartServer() {
	for range only.Once {
		if srv.Port == 0 {
			util.Fail("server.Port must be non-zero")
		}
		// And we serve HTTP until the world ends.
		// Or we are told to shutdown.
		err := srv.Start(fmt.Sprintf("0.0.0.0:%d", srv.Port))
		if errors.Is(err, http.ErrServerClosed) {
			// This is expected, and thus not an error
			break
		}
		if err != nil {
			srv.GetLogger().Fatalf("Failed to start Echo server: %s",
				err)
		}
	}
}

// StopServer terminates the Echo server on localhost
func (srv *Server) StopServer() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		srv.GetLogger().Fatalf("Failed to shutdown Echo server gracefully: %s",
			err)
	}
}

// MakeSwagger Instantiates swagger instance and set server field to `nil`
func (srv *Server) MakeSwagger() *openapi3.Swagger {
	var swagger *openapi3.Swagger
	for range only.Once {
		eps := srv.Endpoints
		if eps == nil {
			util.Fail("Endpoints not set. TODO explain this better.")
			break
		}
		var err error
		swagger, err = eps.GetSwagger()
		if err != nil {
			util.Fail(fmt.Sprintf("unable to load swagger spec\n: %s", err))
			break
		}
		// Clear out the servers array in the swagger spec, that skips validating
		// that server names match. We don't know how this thing will be run.
		swagger.Servers = nil
	}
	return swagger
}

// BearerAuthMiddlewareFunc returns Echo middleware that verifies presence of token in Authorization header
func (srv *Server) BearerAuthMiddlewareFunc() echo.MiddlewareFunc {
	return labstack.KeyAuthWithConfig(labstack.KeyAuthConfig{
		AuthScheme: BearerAuthScheme,
		Validator: func(s string, e echo.Context) (bool, error) {
			// TODO This should validate the Bearer token somehow
			if len(s) > 0 {
				srv.Context.AuthToken = s
				return true, nil
			}
			return false, errors.New("token value is missing from Authorization header")
		},
	})
}
