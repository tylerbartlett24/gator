package main

import (
	"fmt"

	"github.com/tylerbartlett24/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("Tyler")
	cfg = config.Read()
	fmt.Printf("%+v\n", cfg)
}
