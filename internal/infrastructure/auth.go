package infrastructure

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yvv4git/task-voting/internal/domain/entity"
)

type AuthStub struct {
	users map[string]entity.User
}

func NewAuthStub() *AuthStub {
	return &AuthStub{
		users: map[string]entity.User{
			"user1": {
				ID:       uuid.New(),
				Name:     "user1",
				Password: "secret1",
			},
			"user2": {
				ID:       uuid.New(),
				Name:     "user2",
				Password: "secret2",
			},
			"user3": {
				ID:       uuid.New(),
				Name:     "user3",
				Password: "secret3",
			},
			"user4": {
				ID:       uuid.New(),
				Name:     "user4",
				Password: "secret4",
			},
			"user5": {
				ID:       uuid.New(),
				Name:     "user5",
				Password: "secret5",
			},
		},
	}
}

func (a *AuthStub) CheckLoginPassword(_ context.Context, username, password string) error {
	userEntity, ok := a.users[username]
	if !ok {
		return ErrAuthUserNotFound
	}

	if userEntity.Password != password {
		return ErrAuthInvalidCred
	}

	return nil
}

func (a *AuthStub) UserIDByLoginPassword(_ context.Context, username, password string) (uuid.UUID, error) {
	userEntity, ok := a.users[username]
	if !ok {
		return uuid.Nil, ErrAuthUserNotFound
	}

	return userEntity.ID, nil
}

func ExtractBasicAuthValid(c *gin.Context) (string, string, error) {
	authHeader := c.GetHeader("Authorization")

	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", ErrAuthInvalidCred
	}

	encodedCredentials := authHeader[6:]
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return "", "", ErrAuthInvalidCred
	}

	credentials := string(decodedBytes)
	colonIndex := strings.IndexByte(credentials, ':')
	if colonIndex == -1 {
		return "", "", ErrAuthInvalidCred
	}

	username := credentials[:colonIndex]
	password := credentials[colonIndex+1:]

	if len(username) > 0 && len(password) > 0 {
		// Success
		return username, password, nil
	}

	return "", "", ErrAuthInvalidCred
}
