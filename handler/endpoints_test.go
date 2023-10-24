package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aldytanda/swt-pro-tht/generated"
	"github.com/aldytanda/swt-pro-tht/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// TODO: delete
func TestServer_Hello(t *testing.T) {
	type args struct {
		params generated.HelloParams
	}
	type want struct {
		wantErr        error
		wantHTTPStatus int
	}
	tests := []struct {
		name string
		args args
		want
	}{
		{
			name: "Test 1",
			args: args{
				params: generated.HelloParams{
					Id: 1,
				},
			},
			want: want{
				wantErr:        nil,
				wantHTTPStatus: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			// mock repo interface
			repo := repository.NewMockRepositoryInterface(ctrl)

			// create server instance
			s := NewServer(NewServerOptions{
				Repository: repo,
			})

			path := fmt.Sprintf("/hello/%d", tt.args.params.Id)

			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			err := s.Hello(c, tt.args.params)

			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantHTTPStatus, rec.Code)
		})
	}
}

func TestServer_Register(t *testing.T) {
	type fields struct {
		JWTSecretKey string
		Repository   repository.RepositoryInterface
	}
	type args struct {
		ctx  echo.Context
		body string
	}
	type want struct {
		wantErr    error
		wantStatus int
		wantResp   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     want
		mockRepo repository.User
		mockErr  error
	}{
		{
			name: "test register with invalid payload",
			fields: fields{
				JWTSecretKey: "TestingKey",
			},
			args: args{
				body: `{"name":"a","phone":"+62813766712","password":"Abc123!"}`,
			},
			mockErr: nil,
			want: want{
				wantErr:    nil,
				wantStatus: http.StatusBadRequest,
				wantResp: `{
					"errors": [
						{
							"err_rules": [
								"Minimum_Length_3"
							],
							"field_name": "name"
						}
					],
					"message": "Error Validation"
				}`,
			},
		},
		{
			name: "test register with valid payload",
			fields: fields{
				JWTSecretKey: "TestingKey",
			},
			args: args{
				body: `{"name":"aldy","phone":"+62813766711","password":"Abc123!"}`,
			},
			mockRepo: repository.User{
				ID:    1,
				Name:  "aldy",
				Phone: "+62813766711",
			},
			mockErr: nil,
			want: want{
				wantErr:    nil,
				wantStatus: http.StatusCreated,
				wantResp: `{
					"id": 1,
					"name": "aldy",
					"phone": "+62813766711"
				}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			// mock repo interface
			repo := repository.NewMockRepositoryInterface(ctrl)
			if tt.mockRepo.Name != "" {
				repo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tt.mockRepo, tt.mockErr)
			}

			s := &Server{
				JWTSecretKey: tt.fields.JWTSecretKey,
				Repository:   repo,
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(tt.args.body)))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			err := s.Register(c)

			assert.Equal(t, tt.want.wantErr, err)
			assert.Equal(t, tt.want.wantStatus, rec.Code)
			assert.JSONEq(t, tt.want.wantResp, rec.Body.String())
		})
	}
}

func TestServer_UpdateProfile(t *testing.T) {
	type fields struct {
		JWTSecretKey string
		Repository   repository.RepositoryInterface
	}
	type args struct {
		ctx   echo.Context
		token string
		body  string
	}
	type want struct {
		wantErr    error
		wantStatus int
		wantResp   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     want
		mockRepo repository.User
		mockErr  error
	}{
		{
			name: "test update profile with invalid token",
			fields: fields{
				JWTSecretKey: "55860da815619c55f1cb7439e57b60cb8922cb63",
			},
			args: args{
				token: "TestingToken",
				body:  `{"name":"a","phone":"+62813766712"}`,
			},
			mockErr: nil,
			want: want{
				wantErr:    nil,
				wantStatus: http.StatusUnauthorized,
				wantResp: `{
					"message": "Unauthorized"
				}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			// mock repo interface
			repo := repository.NewMockRepositoryInterface(ctrl)
			if tt.mockRepo.Name != "" {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(tt.mockRepo, tt.mockErr)
			}

			s := &Server{
				JWTSecretKey: tt.fields.JWTSecretKey,
				Repository:   repo,
			}

			req := httptest.NewRequest(http.MethodPut, "/users", bytes.NewReader([]byte(tt.args.body)))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", tt.args.token)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			err := s.UpdateProfile(c)

			assert.Equal(t, tt.want.wantErr, err)
			assert.Equal(t, tt.want.wantStatus, rec.Code)
			assert.JSONEq(t, tt.want.wantResp, rec.Body.String())
		})
	}
}

func TestServer_GetProfile(t *testing.T) {
	type fields struct {
		JWTSecretKey string
		Repository   repository.RepositoryInterface
	}
	type args struct {
		ctx   echo.Context
		token string
		body  string
	}
	type want struct {
		wantErr    error
		wantStatus int
		wantResp   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     want
		mockRepo repository.User
		mockErr  error
	}{
		{
			name: "Test with invalid token",
			fields: fields{
				JWTSecretKey: "55860da815619c55f1cb7439e57b60cb8922cb63",
			},
			args: args{
				ctx:   nil,
				token: "Testing Token",
			},
			want: want{
				wantErr:    nil,
				wantStatus: http.StatusUnauthorized,
				wantResp: `{
				"message": "Unauthorized"
			}`,
			},
			mockRepo: repository.User{},
			mockErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			// mock repo interface
			repo := repository.NewMockRepositoryInterface(ctrl)
			if tt.mockRepo.Name != "" {
				repo.EXPECT().GetUser(gomock.Any(), gomock.Any()).Return(tt.mockRepo, tt.mockErr)
			}

			s := &Server{
				JWTSecretKey: tt.fields.JWTSecretKey,
				Repository:   repo,
			}

			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			req.Header.Set("Authorization", tt.args.token)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			err := s.GetProfile(c)

			assert.Equal(t, tt.want.wantErr, err)
			assert.Equal(t, tt.want.wantStatus, rec.Code)
			assert.JSONEq(t, tt.want.wantResp, rec.Body.String())
		})
	}
}

func TestServer_Login(t *testing.T) {
	type fields struct {
		JWTSecretKey string
		Repository   repository.RepositoryInterface
	}
	type args struct {
		ctx   echo.Context
		token string
		body  string
	}
	type want struct {
		wantErr    error
		wantStatus int
		wantResp   string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     want
		mockRepo string
		mockErr  error
	}{
		{
			name: "Test with invalid credentials",
			fields: fields{
				JWTSecretKey: "55860da815619c55f1cb7439e57b60cb8922cb63",
			},
			args: args{
				ctx: nil,
				body: `{
					"phone": "+6281376671",
					"password": "incorrectpassword"
				}`,
			},
			want: want{
				wantErr:    nil,
				wantStatus: http.StatusUnauthorized,
				wantResp: `{
				"message": "Unauthorized: incorrect credentials"
			}`,
			},
			mockRepo: "",
			mockErr:  repository.ErrUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			// mock repo interface
			repo := repository.NewMockRepositoryInterface(ctrl)
			// if tt.mockRepo != "" {
			repo.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tt.mockRepo, tt.mockErr)
			// }

			s := &Server{
				JWTSecretKey: tt.fields.JWTSecretKey,
				Repository:   repo,
			}

			req := httptest.NewRequest(http.MethodPost, "/login", nil)
			req.Header.Set("Authorization", tt.args.token)
			rec := httptest.NewRecorder()

			c := echo.New().NewContext(req, rec)

			err := s.Login(c)

			assert.Equal(t, tt.want.wantErr, err)
			assert.Equal(t, tt.want.wantStatus, rec.Code)
			assert.JSONEq(t, tt.want.wantResp, rec.Body.String())
		})
	}
}
