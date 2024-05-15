package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "9e8dw651cwe6d5c31wc86ew351cw6e534as1",
				Name:   "Admin",
				Age:    20,
				Email:  "admin@admin.com",
				Role:   "admin",
				Phones: []string{"01234567891"},
				meta:   nil,
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1234567890",
				Name:   "User",
				Age:    1,
				Email:  "test.ru",
				Role:   "moderator",
				Phones: []string{"1234567890"},
				meta:   nil,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrValidate,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrValidate,
				},
			},
		},
		{
			in: Response{
				Code: 200,
				Body: "test string",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 100,
				Body: "test",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrValidate,
				},
			},
		},
		{
			in: App{Version: "0.00"},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   ErrValidate,
				},
			},
		},
		{
			in:          App{Version: "10.01"},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte{0, 1, 5},
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in:          "empty",
			expectedErr: ErrNotStructProvided,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			require.Equal(t, tt.expectedErr, Validate(tt.in))
		})
	}
}
