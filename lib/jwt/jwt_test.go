package jwt_test

import (
	"fmt"
	"github.com/k6mil6/distributed-calculator/internal/model"
	"github.com/k6mil6/distributed-calculator/lib/jwt"
	"testing"
	"time"
)

type JWTTest struct {
	UserID        int64
	ErrorExpected bool
}

var tests = []JWTTest{
	{
		UserID:        1,
		ErrorExpected: false,
	},
	{
		UserID:        10002,
		ErrorExpected: false,
	},
}

func TestNewToken(t *testing.T) {

	for _, tt := range tests {
		got, err := jwt.NewToken(model.User{ID: tt.UserID}, 10*time.Minute, "secret")
		if err != nil {
			t.Errorf("NewToken() error = %v", err)
		}
		if got == "" {
			t.Errorf("NewToken() got = %v", got)
		}
		fmt.Println(got)
	}
}

func TestGetUserID(t *testing.T) {
	for _, tt := range tests {
		got, err := jwt.GetUserID("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTMyODkyNDgsImlkIjoxMDAwMiwibG9naW4iOiIifQ.CMbj52UvbwBMvqXujCPkAtOuoVngWRzJNUr4NXIh75k", "secret")
		if err != nil && !tt.ErrorExpected {
			t.Errorf("GetUserID() error = %v", err)
		}
		if err == nil && tt.ErrorExpected {
			t.Errorf("GetUserID() error = %v", err)
		}
		fmt.Print(got)
	}
}
