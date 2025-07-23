package token_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/harungurubudi/mtsg/pkg/token"
	"github.com/harungurubudi/mtsg/testmock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GeneratorTestSuite struct {
	suite.Suite
	mockRedis *testmock.MockRedisAdapter
	generator token.GeneratorRepository
	ctx       context.Context
	key       string
}

func (suite *GeneratorTestSuite) SetupTest() {
	suite.mockRedis = new(testmock.MockRedisAdapter)
	suite.key = "test-secret-key"
	suite.generator = token.NewGenerator(suite.mockRedis, suite.key)
	suite.ctx = context.Background()
}

func (suite *GeneratorTestSuite) TestGenerate_Success() {
	claims := token.Claims{
		Subject:    "access",
		Identifier: "user-123",
		EXP:        time.Now().Add(1 * time.Hour).Unix(),
	}
	suite.mockRedis.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
	tok, err := suite.generator.Generate(suite.ctx, claims)
	suite.NoError(err)
	suite.NotEmpty(tok)
	suite.mockRedis.AssertCalled(suite.T(), "Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything)
}

func (suite *GeneratorTestSuite) TestGenerate_ExpiredClaims() {
	claims := token.Claims{
		Subject:    "access",
		Identifier: "user-123",
		EXP:        time.Now().Add(-1 * time.Hour).Unix(),
	}
	tok, err := suite.generator.Generate(suite.ctx, claims)
	suite.Error(err)
	suite.Empty(tok)
}

func (suite *GeneratorTestSuite) TestGenerate_RedisSetError() {
	claims := token.Claims{
		Subject:    "access",
		Identifier: "user-123",
		EXP:        time.Now().Add(1 * time.Hour).Unix(),
	}
	suite.mockRedis.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(errors.New("redis error"))
	tok, err := suite.generator.Generate(suite.ctx, claims)
	suite.Error(err)
	suite.Empty(tok)
}

func (suite *GeneratorTestSuite) TestValidate_Success() {
	claims := token.Claims{
		Subject:    "access",
		Identifier: "user-123",
		EXP:        time.Now().Add(1 * time.Hour).Unix(),
	}
	tok := token.Token("test-token")
	suite.mockRedis.On("GetByKey", suite.ctx, string(tok), mock.AnythingOfType("*token.Claims")).Run(func(args mock.Arguments) {
		ptr := args.Get(2).(*token.Claims)
		*ptr = claims
	}).Return(nil)
	result, err := suite.generator.Validate(suite.ctx, tok)
	suite.NoError(err)
	suite.Equal(claims.Subject, result.Subject)
	suite.Equal(claims.Identifier, result.Identifier)
}

func (suite *GeneratorTestSuite) TestValidate_NotFound() {
	tok := token.Token("notfound-token")
	suite.mockRedis.On("GetByKey", suite.ctx, string(tok), mock.AnythingOfType("*token.Claims")).Return(errors.New("not found"))
	result, err := suite.generator.Validate(suite.ctx, tok)
	suite.ErrorIs(err, token.ErrTokenNotFound)
	suite.Nil(result)
}

func (suite *GeneratorTestSuite) TestValidate_NBF() {
	claims := token.Claims{
		Subject:    "access",
		Identifier: "user-123",
		NBF:        time.Now().Add(1 * time.Hour).Unix(),
		EXP:        time.Now().Add(2 * time.Hour).Unix(),
	}
	tok := token.Token("future-token")
	suite.mockRedis.On("GetByKey", suite.ctx, string(tok), mock.AnythingOfType("*token.Claims")).Run(func(args mock.Arguments) {
		ptr := args.Get(2).(*token.Claims)
		*ptr = claims
	}).Return(nil)
	result, err := suite.generator.Validate(suite.ctx, tok)
	suite.ErrorIs(err, token.ErrTokenNotYetValid)
	suite.Nil(result)
}

func (suite *GeneratorTestSuite) TestRevoke_Success() {
	tok := token.Token("test-token")
	suite.mockRedis.On("DeleteByKeys", suite.ctx, []string{string(tok)}).Return(nil)
	err := suite.generator.Revoke(suite.ctx, tok)
	suite.NoError(err)
}

func (suite *GeneratorTestSuite) TestRevoke_Error() {
	tok := token.Token("test-token")
	suite.mockRedis.On("DeleteByKeys", suite.ctx, []string{string(tok)}).Return(errors.New("redis error"))
	err := suite.generator.Revoke(suite.ctx, tok)
	suite.Error(err)
}

func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(GeneratorTestSuite))
}
