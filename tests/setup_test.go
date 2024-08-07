package tests

import (
	"context"
	"crypto/sha1"
	"fmt"
	v1 "github.com/dzhordano/avito-bootcamp2024/internal/delivery/http/v1"
	"github.com/dzhordano/avito-bootcamp2024/internal/repository"
	"github.com/dzhordano/avito-bootcamp2024/internal/service"
	"github.com/dzhordano/avito-bootcamp2024/pkg/auth"
	"github.com/dzhordano/avito-bootcamp2024/pkg/databases/postgres"
	"github.com/dzhordano/avito-bootcamp2024/pkg/emails/validation"
	"github.com/dzhordano/avito-bootcamp2024/pkg/logger"
	"github.com/dzhordano/avito-bootcamp2024/pkg/notifications/sender"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	dbDSN    string
	tokenTTL = 6 * time.Hour
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
	dbDSN = os.Getenv("TEST_DB_DSN")
}

type APITestSuite struct {
	suite.Suite

	db       *pgxpool.Pool
	handler  *v1.Handler
	services *service.Services
	repos    *repository.Repository

	emailValidations validation.EmailValidator
	tokensManager    auth.TokensManager
	notifications    sender.Sender
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func (s *APITestSuite) SetupSuite() {
	if client, err := postgres.NewClient(dbDSN); err != nil {
		s.FailNow("failed to connect to postgres: " + err.Error())
	} else {
		s.db = client
	}

	s.initDeps()

	if err := s.seedDB(); err != nil {
		s.FailNow("failed to seed db: " + err.Error())
	}
}

func (s *APITestSuite) initDeps() {
	repos := repository.New(s.db)
	tokensManager := auth.NewJWTManager("secret", tokenTTL)
	notifications := sender.New()
	longTasks := &sync.WaitGroup{}
	emailsValidator := validation.NewEmailValidator()
	inpLogger := logger.NewLogger("debug")

	services := service.New(service.Deps{
		Repos:         repos,
		TokensManager: tokensManager,
		Notifications: notifications,
		WaitGroup:     longTasks,
		Logger:        inpLogger,
	})

	m, err := migrate.New("file://../migrations", dbDSN)
	if err != nil {
		panic(err)
	}

	if err = m.Up(); err != nil {
		panic(err)
	}

	s.repos = repos
	s.emailValidations = emailsValidator
	s.tokensManager = tokensManager
	s.notifications = notifications
	s.services = services
	s.handler = v1.NewHandler(services, tokensManager)
}

func (s *APITestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *APITestSuite) seedDB() error {
	if _, err := s.repos.Houses.Create(context.Background(), house); err != nil {
		return err
	}

	if _, err := s.repos.Flats.Create(context.Background(), houseId, flatApproved); err != nil {
		return err
	}

	if _, err := s.repos.Flats.Create(context.Background(), houseId, flatCreated); err != nil {
		return err
	}

	passwordHash := sha1.Sum([]byte(userModerator.Password))
	userModerator.Password = fmt.Sprintf("%x", passwordHash)

	if err := s.repos.Users.Create(context.Background(), userModerator); err != nil {
		return err
	}

	return nil
}
