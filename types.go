package offstage

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/newclarity/rep-go-offstage/git"
)

type ContextPropertyGetter = git.ContextPropertyGetter

type JSONConverter interface {
	ToJSON(ContextPropertyGetter) (string, error)
}

type EndpointsConnector interface {
	GetSwagger() (*openapi3.Swagger, error)
	RegisterHandlers(interface{}, interface{})
	SetServer(s *Server)
}



