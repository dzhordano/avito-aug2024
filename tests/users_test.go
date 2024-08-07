package tests

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	v1 "github.com/dzhordano/avito-bootcamp2024/internal/delivery/http/v1"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestUsersDummyLoginClient() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	userType := domain.UserTypeClient.String()

	req, _ := http.NewRequest("GET", "/api/auth/dummyLogin", nil)
	q := req.URL.Query()
	q.Add("userType", userType)
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var authToken map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &authToken)
	s.NoError(err)

	utTokenResp, err := s.tokensManager.Parse(authToken["auth_token"])
	s.NoError(err)

	r.Equal(userType, utTokenResp)

}

func (s *APITestSuite) TestUsersDummyLoginInvalidType() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// Bad user type
	userType := "test"

	req, _ := http.NewRequest("GET", "/api/auth/dummyLogin", nil)
	q := req.URL.Query()
	q.Add("userType", userType)
	req.URL.RawQuery = q.Encode()

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
}

func (s *APITestSuite) TestUsersRegisterSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.UserRegisterInput{
		Email:    "tester1@mail.ru",
		Password: "test",
		UserType: domain.UserTypeClient,
	}

	b, _ := json.Marshal(input)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(b))

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	var respBody v1.UserIdResponse
	err := json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)

	r.NotEmpty(respBody)

	query, args, err := squirrel.
		Select("user_id", "email", "password_hash", "user_type").
		From("users").
		Where(squirrel.Eq{"user_id": respBody.UserId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	s.NoError(err)

	var user domain.User
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.UserType)
	s.NoError(err)

	// Hash password to compare with DB password.
	passwordHash := sha1.Sum([]byte(input.Password))
	input.Password = fmt.Sprintf("%x", passwordHash)

	r.Equal(respBody.UserId, user.ID.String())
	r.Equal(input.Email, user.Email)
	r.Equal(input.Password, user.Password)
	r.Equal(input.UserType, user.UserType)
}

func (s *APITestSuite) TestUsersLoginSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.UserLoginInput{
		Email:    "initTester@mail.ru",
		Password: "qwerty",
	}

	b, _ := json.Marshal(input)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	var token map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &token)
	s.NoError(err)

	r.NotEmpty(token)
}

func (s *APITestSuite) TestUsersLoginNotFound() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.UserLoginInput{
		Email:    "unknowTester@mail.ru",
		Password: "test",
	}

	b, _ := json.Marshal(input)

	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(b))

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusNotFound, resp.Result().StatusCode)
}
