package whydoweneedtest

import (
	"github.com/Rayato159/hello-sekai-shop-tutorial/config"
)

func NewTestConfig() *config.Config {
	cfg := config.LoadConfig("../env/test/.env")
	return &cfg
}
