package user

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DatatypesTestSuite struct {
	suite.Suite
}

func TestDatatypesTestSuite(t *testing.T) {
	suite.Run(t, new(DatatypesTestSuite))
}

func (suite *DatatypesTestSuite) TestPasswordValidate() {
	tests := []struct {
		name     string
		password Password
		want     bool
	}{
		{"valid password", "Abcdef1!", true},
		{"too short", "Ab1!", false},
		{"no uppercase", "abcdef1!", false},
		{"no lowercase", "ABCDEF1!", false},
		{"no number", "Abcdefg!", false},
		{"no special", "Abcdefg1", false},
		{"all requirements", "A1b2c3d4!", true},
		{"only special", "!!!!!!!!", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			suite.Equal(tt.want, tt.password.Validate())
		})
	}
}

func (suite *DatatypesTestSuite) TestNewCipherTextWithValidPasswordReturnsCipherText() {
	tests := []struct {
		name     string
		password Password
		wantErr  bool
	}{
		{"valid password", "MySecurePassword123!", false},
		{"empty password", "", false},
		{"complex password with special chars", "P@ssw0rd!@#$%^&*()", false},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ciphertext, err := NewCipherText(tt.password)
			if tt.wantErr {
				suite.Error(err)
				suite.Empty(ciphertext)
			} else {
				suite.NoError(err)
				suite.NotEmpty(ciphertext)
				suite.IsType(CipherText(""), ciphertext)
			}
		})
	}
}

func (suite *DatatypesTestSuite) TestCipherTextMatchesWithCorrectPasswordReturnsTrue() {
	tests := []struct {
		name     string
		password Password
	}{
		{"simple password", "password123"},
		{"complex password", "MySecurePassword123!"},
		{"empty password", ""},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ciphertext, err := NewCipherText(tt.password)
			suite.NoError(err)
			result := ciphertext.Matches(tt.password)
			suite.True(result)
		})
	}
}

func (suite *DatatypesTestSuite) TestCipherTextMatchesWithIncorrectPasswordReturnsFalse() {
	tests := []struct {
		name          string
		originalPass  Password
		incorrectPass Password
	}{
		{"different password", "password123", "password456"},
		{"case sensitive", "Password123", "password123"},
		{"empty vs non-empty", "", "somepassword"},
		{"non-empty vs empty", "somepassword", ""},
	}
	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ciphertext, err := NewCipherText(tt.originalPass)
			suite.NoError(err)
			result := ciphertext.Matches(tt.incorrectPass)
			suite.False(result)
		})
	}
}

func (suite *DatatypesTestSuite) TestCipherTextMatchesWithSamePasswordMultipleTimesReturnsConsistentResults() {
	password := Password("MySecurePassword123!")
	ciphertext, err := NewCipherText(password)
	suite.NoError(err)
	for i := 0; i < 5; i++ {
		result := ciphertext.Matches(password)
		suite.True(result, "Password should match consistently")
	}
}
