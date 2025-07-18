package token

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "github.com/harungurubudi/mtsg/testmock"
)

type GeneratorTestSuite struct {
    suite.Suite
    mockRedis *testmock.MockRedisAdapter
    generator *Generator
    ctx       context.Context
    key       string
}

func (suite *GeneratorTestSuite) SetupTest() {
    suite.mockRedis = new(testmock.MockRedisAdapter)
    suite.key = "test-secret-key"
    suite.generator = NewGenerator(suite.mockRedis, suite.key)
    suite.ctx = context.Background()
}

func (suite *GeneratorTestSuite) TestGenerate_Success() {
    claims := Claims{
        Subject:    "access",
        Identifier: "user-123",
        EXP:        time.Now().Add(1 * time.Hour).Unix(),
    }
    suite.mockRedis.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(nil)
    token, err := suite.generator.Generate(suite.ctx, claims)
    suite.NoError(err)
    suite.NotEmpty(token)
    suite.mockRedis.AssertCalled(suite.T(), "Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything)
}

func (suite *GeneratorTestSuite) TestGenerate_ExpiredClaims() {
    claims := Claims{
        Subject:    "access",
        Identifier: "user-123",
        EXP:        time.Now().Add(-1 * time.Hour).Unix(),
    }
    token, err := suite.generator.Generate(suite.ctx, claims)
    suite.Error(err)
    suite.Empty(token)
}

func (suite *GeneratorTestSuite) TestGenerate_RedisSetError() {
    claims := Claims{
        Subject:    "access",
        Identifier: "user-123",
        EXP:        time.Now().Add(1 * time.Hour).Unix(),
    }
    suite.mockRedis.On("Set", suite.ctx, mock.AnythingOfType("string"), mock.Anything, mock.Anything).Return(errors.New("redis error"))
    token, err := suite.generator.Generate(suite.ctx, claims)
    suite.Error(err)
    suite.Empty(token)
}

func (suite *GeneratorTestSuite) TestValidate_Success() {
    claims := Claims{
        Subject:    "access",
        Identifier: "user-123",
        EXP:        time.Now().Add(1 * time.Hour).Unix(),
    }
    token := Token("test-token")
    suite.mockRedis.On("GetByKey", suite.ctx, string(token), mock.AnythingOfType("*token.Claims")).Run(func(args mock.Arguments) {
        ptr := args.Get(2).(*Claims)
        *ptr = claims
    }).Return(nil)
    result, err := suite.generator.Validate(suite.ctx, token)
    suite.NoError(err)
    suite.Equal(claims.Subject, result.Subject)
    suite.Equal(claims.Identifier, result.Identifier)
}

func (suite *GeneratorTestSuite) TestValidate_NotFound() {
    token := Token("notfound-token")
    suite.mockRedis.On("GetByKey", suite.ctx, string(token), mock.AnythingOfType("*token.Claims")).Return(errors.New("not found"))
    result, err := suite.generator.Validate(suite.ctx, token)
    suite.ErrorIs(err, ErrTokenNotFound)
    suite.Nil(result)
}

func (suite *GeneratorTestSuite) TestValidate_NBF() {
    claims := Claims{
        Subject:    "access",
        Identifier: "user-123",
        NBF:        time.Now().Add(1 * time.Hour).Unix(),
        EXP:        time.Now().Add(2 * time.Hour).Unix(),
    }
    token := Token("future-token")
    suite.mockRedis.On("GetByKey", suite.ctx, string(token), mock.AnythingOfType("*token.Claims")).Run(func(args mock.Arguments) {
        ptr := args.Get(2).(*Claims)
        *ptr = claims
    }).Return(nil)
    result, err := suite.generator.Validate(suite.ctx, token)
    suite.ErrorIs(err, ErrTokenNotYetValid)
    suite.Nil(result)
}

func (suite *GeneratorTestSuite) TestRevoke_Success() {
    token := Token("test-token")
    suite.mockRedis.On("DeleteByKeys", suite.ctx, []string{string(token)}).Return(nil)
    err := suite.generator.Revoke(suite.ctx, token)
    suite.NoError(err)
}

func (suite *GeneratorTestSuite) TestRevoke_Error() {
    token := Token("test-token")
    suite.mockRedis.On("DeleteByKeys", suite.ctx, []string{string(token)}).Return(errors.New("redis error"))
    err := suite.generator.Revoke(suite.ctx, token)
    suite.Error(err)
}

func TestGeneratorTestSuite(t *testing.T) {
    suite.Run(t, new(GeneratorTestSuite))
} 