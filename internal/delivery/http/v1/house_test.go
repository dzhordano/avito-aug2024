package v1

import (
	"errors"
	"github.com/dzhordano/avito-bootcamp2024/internal/domain"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	mocks_service "github.com/dzhordano/avito-bootcamp2024/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetHouseById(t *testing.T) {
	type mockBehaviour func(s *mocks_service.MockHouses, id int)

	tests := []struct {
		name               string
		id                 int
		mockBehaviour      mockBehaviour
		expectedStatusCode int
		expectedReqBody    string
	}{
		{
			name: "OK",
			id:   1,
			mockBehaviour: func(s *mocks_service.MockHouses, id int) {
				s.
					EXPECT().
					GetById(gomock.Any(), id).
					Return([]domain.Flat{
						{
							ID:         1,
							FlatNumber: 256,
							Price:      1000,
							Rooms:      3,
							Status:     domain.StatusCreated,
						},
					}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedReqBody:    `{"data":[{"ID":1,"FlatNumber":256,"Price":1000,"Rooms":3,"Status":"created"}]}`,
		},
		{
			name: "Not found",
			id:   1,
			mockBehaviour: func(s *mocks_service.MockHouses, id int) {
				s.
					EXPECT().
					GetById(gomock.Any(), id).
					Return(nil, domain.ErrHouseNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedReqBody:    `{"message":"house not found"}`,
		},
		{
			name: "Internal server error",
			mockBehaviour: func(s *mocks_service.MockHouses, id int) {
				s.
					EXPECT().
					GetById(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("internal server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedReqBody:    `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			houses := mocks_service.NewMockHouses(c)
			tt.mockBehaviour(houses, tt.id)

			services := &service.Services{
				Houses: houses,
			}

			handler := NewHandler(services, nil)

			r := gin.New()
			r.GET("/api/house/:id", handler.getHouseById)

			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/api/house/1", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedReqBody, w.Body.String())
		})
	}
}
