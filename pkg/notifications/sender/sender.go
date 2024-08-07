package sender

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Sender interface {
	SendEmail(ctx context.Context, recipient string, message string) error
}

type sender struct{}

func New() *sender {
	return &sender{}
}

func (s *sender) SendEmail(ctx context.Context, recipient string, message string) error {
	// Имитация отправки сообщения
	duration := time.Duration(rand.Int63n(3000)) * time.Millisecond
	time.Sleep(duration)

	// Имитация неуспешной отправки сообщения
	errorProbability := 0.1
	if rand.Float64() < errorProbability {
		return errors.New("internal error")
	}

	fmt.Printf("send message '%s' to '%s'\n", message, recipient)

	return nil
}
