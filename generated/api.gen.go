// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Errors  *[]ValidationError `json:"errors,omitempty"`
	Message string             `json:"message"`
}

// GetProfileResponse defines model for GetProfileResponse.
type GetProfileResponse struct {
	CountLogin int    `json:"count_login"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

// HelloResponse defines model for HelloResponse.
type HelloResponse struct {
	Message string `json:"message"`
}

// LoginRequest defines model for LoginRequest.
type LoginRequest struct {
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterRequest defines model for RegisterRequest.
type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

// RegisterResponse defines model for RegisterResponse.
type RegisterResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// UpdateProfileRequest defines model for UpdateProfileRequest.
type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

// ValidationError defines model for ValidationError.
type ValidationError struct {
	ErrRules  []string `json:"err_rules"`
	FieldName string   `json:"field_name"`
}

// HelloParams defines parameters for Hello.
type HelloParams struct {
	Id int `form:"id" json:"id"`
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LoginRequest

// RegisterJSONRequestBody defines body for Register for application/json ContentType.
type RegisterJSONRequestBody = RegisterRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// This is just a test endpoint to get you started. Please delete this endpoint.
	// (GET /hello)
	Hello(ctx echo.Context, params HelloParams) error
	// User login
	// (POST /login)
	Login(ctx echo.Context) error
	// Get user profile
	// (GET /users)
	GetProfile(ctx echo.Context) error
	// Register a new user
	// (POST /users)
	Register(ctx echo.Context) error
	// Update a user
	// (PUT /users)
	UpdateProfile(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// Hello converts echo context to params.
func (w *ServerInterfaceWrapper) Hello(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params HelloParams
	// ------------- Required query parameter "id" -------------

	err = runtime.BindQueryParameter("form", true, true, "id", ctx.QueryParams(), &params.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Hello(ctx, params)
	return err
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Login(ctx)
	return err
}

// GetProfile converts echo context to params.
func (w *ServerInterfaceWrapper) GetProfile(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProfile(ctx)
	return err
}

// Register converts echo context to params.
func (w *ServerInterfaceWrapper) Register(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Register(ctx)
	return err
}

// UpdateProfile converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateProfile(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateProfile(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/hello", wrapper.Hello)
	router.POST(baseURL+"/login", wrapper.Login)
	router.GET(baseURL+"/users", wrapper.GetProfile)
	router.POST(baseURL+"/users", wrapper.Register)
	router.PUT(baseURL+"/users", wrapper.UpdateProfile)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xX33PbNgz+V3DcHp3YXfsyPXbbdb396qXNXnq5HCfCMjOKZEDQnZbz/74jJcuWLTXb",
	"zXH6Jksg8OED8BF+EKWrvbNoOYjiQYRyhbXMjz8QObrC4J0NmF54ch6JNebPmD7nJ81Y54evCZeiEF/N",
	"dz7nncP579JoJVk7m/2KzUxw41EUQhLJJv2uMQRZ5VDdp8CkbSU2m5kgvI+aUIniY2940/twf9xhycnJ",
	"G+R35Jba4DT00kXLt8ZV2u4F05axwozMynoMxkz4lbP/AmA+v7WeDeKNYf4RjXHTcE/Cy88p+hXeRwx8",
	"HMLLED45Uv8n6W22va/PwJhKld2faB8P1ZqN+b/CSgdGmsx0urQnoOCg7p9lYod0igytnqY7tRKzIdQx",
	"fNdeScZ+lv4zndNYjkIdasOY1txSNDiUm6OYh4Ky1GjU7QTEA1L2bGd78Y6ZSQe1Xbrk0ugSu9K1UcQv",
	"bz9kJJpN+nkdkOA90lqXye8aKWhnRSFeXC4uF8nSebTSa1GIl/lV6hpe5fzmq6QL6anCTH2iJNP0Vomi",
	"VY1sT7JGxiTGHx9E0jRxH5GabZWLtuK7bJkizjqpH2uxzU2ybtsyI/lmsWh10zLaDEV6b3SZwczvgrO7",
	"u+Oxi2AodplOhaEk7bml5gMGBkKOZBNBrxavThZ7eKWNxP7VMSxdtCr3R4h1LalJmFY6gA5wFwODBE4Q",
	"0SrvtGVgBxUyNC5CYEmM6hLeGZQBQaFBRuB0fGt/mX3P++vHuzBS3SyUXdUw8GunmpPRMLgLNsNJSL2x",
	"OSr/i1PHni7B+1iWGMIymgaMqypUoLtGWJyvEV5LBTuG9lshj3RbvFzHGLBdg0aHdLeOiCecqZGl5xFm",
	"CZk0rlFBwg++w5hpfnk+mq+tjLxypP9GBTTK9xvkI4zjI7O9UZ9oag5XizMPztG+8EiFS0LJqEBm9r6g",
	"AdomAhIsfurR+ThS0sES8pQTNLrtjKT1209wAfs0x3xQPR/JXavDBXjZGCcVrPt1Cro/ac861hc9RHbt",
	"JBMGF6ns5Obb8+H6ztml0WUi63vJIA2hVA3gXzpwAG0hNIGxPlT8XOJ+kPJHpPV24YpkRCFWzL6Yz40r",
	"pVkledrcbP4JAAD//xmV8YVdDwAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}