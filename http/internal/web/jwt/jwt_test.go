package jwtauth

import (
	"log"
	"testing"

	mockuser "github.com/Lafetz/loyalty_marketplace/internal/web/mockUser"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJwt(t *testing.T) {
	t.Run("Successfuly create jwt", func(t *testing.T) {
		user := mockuser.User{
			Id:       uuid.New(),
			Username: "helloworld",
			Email:    "hellow@world.com",
		}
		token, err := CreateJwt(user)
		if err != nil {
			log.Fatal(err)
		}
		assert.True(t, len(token) > 0)
		assert.Nil(t, err)

	})

	t.Run("Successfuly parse token", func(t *testing.T) {
		user := mockuser.User{
			Id:       uuid.New(),
			Username: "helloworld",
			Email:    "hellow@world.com",
		}
		token, err := CreateJwt(user)
		if err != nil {
			log.Fatal(err)
		}
		_, err = PareseJwt(token)
		assert.Nil(t, err)
	})

}
