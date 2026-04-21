package finam

import (
	"context"
	"sync"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"google.golang.org/grpc/metadata"
)

type tokenService struct {
	authService auth.AuthServiceClient
	secret      string
	token       string
	expiresAt   time.Time
	mx          sync.Mutex
}

func newTokenService(authService auth.AuthServiceClient, secret string) *tokenService {
	return &tokenService{
		authService: authService,
		secret:      secret,
	}
}

func (s *tokenService) Context(ctx context.Context) (context.Context, error) {
	token, err := s.getToken()
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(ctx, "Authorization", token), nil
}

func (s *tokenService) getToken() (string, error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if s.token == "" || time.Since(s.expiresAt)+time.Second > 0 {
		err := s.refreshToken()
		if err != nil {
			return "", err
		}
	}
	return s.token, nil
}

func (s *tokenService) refreshToken() error {
	resp, err := s.authService.Auth(context.Background(), &auth.AuthRequest{Secret: s.secret})
	if err != nil {
		return err
	}

	s.token = resp.GetToken()

	details, err := s.authService.TokenDetails(context.Background(), &auth.TokenDetailsRequest{
		Token: s.token,
	})
	if err != nil {
		return err
	}

	s.expiresAt = details.ExpiresAt.AsTime()

	return nil
}
