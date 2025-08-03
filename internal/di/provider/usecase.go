package provider

import (
	"github.com/harungurubudi/mtsg/internal/repository"
	"github.com/harungurubudi/mtsg/internal/usecase"
	"github.com/harungurubudi/mtsg/pkg/token"
)

// Use Case Providers

// ProvideAuthUseCase injects dependencies into Authentication usecase
func ProvideAuthUseCase(
	userRepo repository.UserRepository,
	tokenGen token.GeneratorRepository,
) usecase.Authentication {
	return usecase.NewAuthentication(userRepo, tokenGen)
}
