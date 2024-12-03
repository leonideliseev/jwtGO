package main

import (
	"fmt"

	"github.com/leonideliseev/jwtGO/config"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		panic(fmt.Sprintf("error init configs: %s", err.Error()))
	}

	err = config.LoadEnv()
	if err != nil {
		panic(fmt.Sprintf("error loading env: %s", err.Error()))
	}

}
