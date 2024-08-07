package tests

import (
	"bytes"
	"context"
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

func (s *APITestSuite) TestHousesCreateSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.HouseCreateInput{
		Address:   "test address 1",
		Year:      2001,
		Developer: "good developer",
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeModerator.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/house/create", bytes.NewBuffer(b))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var respBody v1.DataResponse[domain.House]
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)

	fmt.Println(respBody)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)
	r.Equal(input.Address, respBody.Data.Address)
	r.Equal(input.Year, respBody.Data.Year)
	r.Equal(input.Developer, respBody.Data.Developer)
}

func (s *APITestSuite) TestHousesCreateUnauthorized() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.HouseCreateInput{
		Address:   "test address 1",
		Year:      2001,
		Developer: "good developer",
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/house/create", bytes.NewBuffer(b))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}

func (s *APITestSuite) TestHousesGetFlatsById() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/house/%d", 1), nil)

	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var respBody v1.DataResponse[[]domain.Flat]
	err = json.Unmarshal(resp.Body.Bytes(), &respBody)
	s.NoError(err)

	// Compare to flat that was created through TestFlatsCreateSuccess.
	r.Equal(http.StatusOK, resp.Result().StatusCode)
	r.Equal(1, len(respBody.Data))
	r.Equal(1, respBody.Data[0].ID)
	r.Equal(53, respBody.Data[0].FlatNumber)
	r.Equal(20000, respBody.Data[0].Price)
	r.Equal(4, respBody.Data[0].Rooms)
	r.Equal(domain.StatusApproved, respBody.Data[0].Status)
}

func (s *APITestSuite) TestHousesSubscribeSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	// User that was created in data.
	input := dtos.HouseSubscribeInput{
		Email: "initTester@mail.ru",
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", fmt.Sprintf("/api/house/%d/subscribe", 1), bytes.NewBuffer(b))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	query, args, err := squirrel.
		Select("user_email").
		From("house_subscriptions").
		Where(squirrel.Eq{"house_id": 1}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	s.NoError(err)

	var email string
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&email)
	s.NoError(err)

	r.Equal(http.StatusOK, resp.Result().StatusCode)
	r.Equal(input.Email, email)
}
