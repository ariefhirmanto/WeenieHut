package service

import "context"

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	return "token", nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (string, error) {
	return "token", nil
}
