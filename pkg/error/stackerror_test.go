package error_test

import (
	"errors"
	"testing"

	errorpkg "github.com/harungurubudi/mtsg/pkg/error"
	"github.com/stretchr/testify/suite"
)

type StackErrorTestSuite struct {
	suite.Suite
}

func (suite *StackErrorTestSuite) TestErrorMessage() {
	err := errorpkg.NewStackError("something went wrong", nil)
	suite.Equal("something went wrong", err.Error())
}

func (suite *StackErrorTestSuite) TestErrorWrapping() {
	cause := errors.New("root cause")
	err := errorpkg.NewStackError("outer error", cause)
	suite.Contains(err.Error(), "outer error")
	suite.Contains(err.Error(), "root cause")
	suite.Equal(cause, errors.Unwrap(err))
}

func (suite *StackErrorTestSuite) TestStackTrace() {
	err := errorpkg.NewStackError("trace error", nil)
	trace := err.StackTrace()
	suite.Contains(trace, "TestStackTrace")
	suite.Contains(trace, ".go:")
}

func (suite *StackErrorTestSuite) TestNewStackErrorf() {
	err := errorpkg.NewStackErrorf("error: %d", 42)
	suite.Equal("error: 42", err.Error())
}

func (suite *StackErrorTestSuite) TestCallerInfo() {
	err := errorpkg.NewStackError("caller info", nil)
	suite.NotEmpty(err.Function)
	suite.NotEmpty(err.File)
	suite.NotZero(err.Line)
}

func TestStackErrorTestSuite(t *testing.T) {
	suite.Run(t, new(StackErrorTestSuite))
}
