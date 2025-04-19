package main

import "api/internal/config"

func main() {
	config, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}

	_ = config // Use the config as needed
}
