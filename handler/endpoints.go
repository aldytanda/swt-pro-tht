package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aldytanda/swt-pro-tht/generated"
	"github.com/aldytanda/swt-pro-tht/repository"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// TODO: delete
// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

// (POST /users)
func (s *Server) Register(ctx echo.Context) error {
	var payload generated.RegisterRequest

	err := ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	ctxx := ctx.Request().Context()
	err = validateRegister(ctxx, payload)
	if err != nil {
		if validationErr, ok := err.(ValidationErrorResp); ok {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "Error Validation",
				Errors:  validationErr.ToResponseErrors(),
			})
		}

		log.Println(err)

		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Error Validation",
		})
	}

	user, err := s.Repository.CreateUser(ctxx, repository.User{
		Name:     payload.Name,
		Phone:    payload.Phone,
		Password: payload.Password,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateColumn) {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "Data [phone number] already exists.",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp := generated.RegisterResponse{
		Id:    int(user.ID),
		Name:  user.Name,
		Phone: user.Phone,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

// (PUT /users)
func (s *Server) UpdateProfile(ctx echo.Context) error {
	var payload generated.UpdateProfileRequest

	jwtClaims, ok := s.verifyToken(ctx)
	if !ok {
		log.Println("err verify token")
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	idStr, ok := jwtClaims["id"].(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("err token: id type is invalid")
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	err = ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload: failed to parse",
		})
	}

	ctxx := ctx.Request().Context()
	err = validateUpdateUser(ctxx, payload)
	if err != nil {
		if validationErr, ok := err.(ValidationErrorResp); ok {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "Error Validation",
				Errors:  validationErr.ToResponseErrors(),
			})
		}

		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Error Validation",
		})
	}

	updateUser, err := s.Repository.GetUser(ctxx, id)
	if err != nil {
		log.Println("Get Profile err: %w", err)

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Internal Server Error",
		})
	}

	updateUser.ID = id
	if payload.Name != nil {
		updateUser.Name = *payload.Name
	}
	if payload.Phone != nil {
		updateUser.Phone = *payload.Phone
	}

	user, err := s.Repository.UpdateUser(ctxx, updateUser)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateColumn) {
			return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
				Message: "Data [phone number] already exists.",
			})
		} else if errors.Is(err, repository.ErrNotFound) {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "Bad Request",
			})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	resp := generated.RegisterResponse{
		Id:    int(user.ID),
		Name:  user.Name,
		Phone: user.Phone,
	}

	return ctx.JSON(http.StatusOK, resp)
}

// (GET /users)
func (s *Server) GetProfile(ctx echo.Context) error {
	jwtClaims, ok := s.verifyToken(ctx)
	if !ok {
		log.Println("err verify token")
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	idStr, ok := jwtClaims["id"].(string)
	if !ok {
		log.Println("err token: id is empty")
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Println("err token: id type is invalid")
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	user, err := s.Repository.GetUser(ctx.Request().Context(), id)
	if err != nil {
		log.Println("Get Profile err: %w", err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Internal Server Error",
		})
	}

	resp := generated.GetProfileResponse{
		Name:       user.Name,
		Phone:      user.Phone,
		CountLogin: int(user.CountLogin),
	}

	return ctx.JSON(http.StatusOK, resp)
}

// (POST /login)
func (s *Server) Login(ctx echo.Context) error {
	var payload generated.LoginRequest

	err := ctx.Bind(&payload)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid payload",
		})
	}

	ctxx := ctx.Request().Context()
	token, err := s.Repository.Login(ctxx, repository.Login{
		Phone:    payload.Phone,
		Password: payload.Password,
	})
	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Unauthorized: incorrect credentials",
		})
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Token: token,
	})
}

func (s *Server) verifyToken(ctx echo.Context) (map[string]interface{}, bool) {
	result := map[string]interface{}{}

	if ctx.Request().Header["Authorization"] == nil {
		log.Println("empty authorization token")
		return result, false
	}

	log.Println(ctx.Request().Header["Authorization"][0])
	token, err := jwt.Parse(ctx.Request().Header["Authorization"][0], func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Println("err signing method")
			return nil, errors.New("Unauthorized")
		}

		return []byte(s.JWTSecretKey), nil
	})
	if err != nil {
		log.Printf("err parsing token: %s\n", err)
		return result, false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		result["id"] = claims["sub"]
		result["user"] = claims["user"]

		return result, true
	}

	log.Println("token invalid")
	return result, false
}
