package main

import (
	"fmt"
	"iosipoff/url-shortener/cmd/url-shortener/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Printf("config: %+v", cfg)

	//TODO: init logger: slog

	//TODO: init storage: sqlite

	//TODO: init router: chi, chi render

	//TODO: run server
}
