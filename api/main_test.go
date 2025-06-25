package api

import (
	"os"
	db "simplebank/db/model"
	"simplebank/factory"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := factory.Config{
		TokenSymmetricKey: factory.RandomString(32),
		TokenDuration:     15 * time.Minute,
	}

	server, err := NewServer(config, store)
	if err != nil {
		t.Fatal("Error! No es pot inicialitzar el server: ", err)
	}

	return server
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
