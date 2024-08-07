package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/dtos"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
)

func (s *APITestSuite) TestFlatsCreateSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.FlatCreateInput{
		FlatNumber: 1,
		HouseId:    1,
		Price:      1000,
		Rooms:      2,
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/flat/create", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	query, args, err := squirrel.
		Select("id", "flat_number", "price", "rooms", "status").
		From("flats").
		Where(squirrel.Eq{"flat_number": input.FlatNumber}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	s.NoError(err)

	var res domain.Flat
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&res.ID, &res.FlatNumber, &res.Price, &res.Rooms, &res.Status)
	s.NoError(err)

	r.NotEmpty(res.ID)
	r.Equal(input.FlatNumber, res.FlatNumber)
	r.Equal(input.Price, res.Price)
	r.Equal(input.Rooms, res.Rooms)
	s.NoError(err)
}

func (s *APITestSuite) TestFlatsCreateFlatAlreadyExists() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.FlatCreateInput{
		FlatNumber: 53,
		HouseId:    1,
		Price:      1000,
		Rooms:      2,
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/flat/create", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusInternalServerError, resp.Result().StatusCode)
}

func (s *APITestSuite) TestFlatsUpdateSuccess() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.FlatUpdateInput{
		FlatId: 1,
		Status: domain.StatusApproved,
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeModerator.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/flat/update", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusOK, resp.Result().StatusCode)

	query, args, err := squirrel.
		Select("id", "status").
		From("flats").
		Where(squirrel.Eq{"id": input.FlatId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	s.NoError(err)

	var res domain.Flat
	err = s.db.QueryRow(context.Background(), query, args...).Scan(&res.ID, &res.Status)

	r.Equal(input.Status, res.Status)
	r.Equal(input.FlatId, res.ID)
	s.NoError(err)
}

func (s *APITestSuite) TestFlatsUpdateUnauthorized() {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	s.handler.Init(router.Group("/api"))
	r := s.Require()

	input := dtos.FlatUpdateInput{
		FlatId: 1,
		Status: domain.StatusApproved,
	}

	b, _ := json.Marshal(input)

	token, err := s.tokensManager.GenerateJWT(domain.UserTypeClient.String())
	s.NoError(err)

	req, _ := http.NewRequest("POST", "/api/flat/update", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)
}

//func (s *APITestSuite) TestFlatsUpdate() {
//	gin.SetMode(gin.TestMode)
//
//	router := gin.New()
//	s.handler.Init(router.Group("/api"))
//	r := s.Require()
//
//	testsTable := []struct {
//		testName   string
//		input      dtos.FlatUpdateInput
//		statusCode int
//	}{
//		{
//			"Success",
//			dtos.FlatUpdateInput{
//				FlatId: 1,
//				Status: domain.StatusApproved,
//			},
//			http.StatusOK,
//		},
//		{
//			"Empty flat number",
//			dtos.FlatUpdateInput{
//				FlatId: 1,
//				Status: domain.StatusApproved,
//			},
//			http.StatusBadRequest,
//		},
//	}
//
//	for _, tt := range testsTable {
//
//		b, _ := json.Marshal(tt.input)
//
//		req, _ := http.NewRequest("POST", "/api/flat/update", bytes.NewBuffer(b))
//
//		resp := httptest.NewRecorder()
//		router.ServeHTTP(resp, req)
//	}
//}
