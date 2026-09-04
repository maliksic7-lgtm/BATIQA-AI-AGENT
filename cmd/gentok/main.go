package main

import (
	"fmt"
	"os"
	"time"

	"batiqa-ai/internal/config"
	"batiqa-ai/internal/guesttoken"
)

func main() {
	config.Load()
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/gentok <room>")
		os.Exit(1)
	}
	tok, err := guesttoken.New(os.Args[1], 90*24*time.Hour, guesttoken.Secret())
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERR:", err)
		os.Exit(1)
	}
	fmt.Print(tok)
}
