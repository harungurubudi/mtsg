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

func (suite *DatatypesTestSuite) TestNewCiphertextWithValidPasswordReturnsCiphertext() {
    tests := []struct {
        name     string
        password Password
        wantErr  bool
    }{
        {
            name:     "valid password",
            password: "MySecurePassword123!",
            wantErr:  false,
        },
        {
            name:     "empty password",
            password: "",
            wantErr:  false,
        },
        {
            name:     "complex password with special chars",
            password: "P@ssw0rd!@#$%^&*()",
            wantErr:  false,
        },
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            ciphertext, err := NewCiphertext(tt.password)
            
            if tt.wantErr {
                suite.Error(err)
                suite.Empty(ciphertext)
            } else {
                suite.NoError(err)
                suite.NotEmpty(ciphertext)
                suite.IsType(Ciphertext(""), ciphertext)
            }
        })
    }
}

func (suite *DatatypesTestSuite) TestCiphertextMatchesWithCorrectPasswordReturnsTrue() {
    tests := []struct {
        name     string
		password Password
	}{
		{
			name:     "simple password",
			password: "password123",
		},
		{
			name:     "complex password",
			password: "MySecurePassword123!",
		},
		{
			name:     "empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ciphertext, err := NewCiphertext(tt.password)
			suite.NoError(err)
			
			result := ciphertext.Matches(tt.password)
			suite.True(result)
		})
	}
}

func (suite *DatatypesTestSuite) TestCiphertextMatchesWithIncorrectPasswordReturnsFalse() {
	tests := []struct {
		name           string
		originalPass   Password
		incorrectPass  Password
	}{
		{
			name:          "different password",
			originalPass:  "password123",
			incorrectPass: "password456",
		},
		{
			name:          "case sensitive",
			originalPass:  "Password123",
			incorrectPass: "password123",
		},
		{
			name:          "empty vs non-empty",
			originalPass:  "",
			incorrectPass: "somepassword",
		},
		{
			name:          "non-empty vs empty",
			originalPass:  "somepassword",
			incorrectPass: "",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			ciphertext, err := NewCiphertext(tt.originalPass)
			suite.NoError(err)
			
			result := ciphertext.Matches(tt.incorrectPass)
			suite.False(result)
		})
	}
}

func (suite *DatatypesTestSuite) TestCiphertextMatchesWithSamePasswordMultipleTimesReturnsConsistentResults() {
	password := Password("MySecurePassword123!")
	
	ciphertext, err := NewCiphertext(password)
	suite.NoError(err)
	
	// Test multiple times to ensure consistency
	for i := 0; i < 5; i++ {
		result := ciphertext.Matches(password)
		suite.True(result, "Password should match consistently")
	}
} 