package main

import (
	"URLShortenePetPrpoject/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
