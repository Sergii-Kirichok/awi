package main

import (
	"awi/awp"
	"awi/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cfg := config.New()
	cfg.Load(rootDir)

	l := awp.NewAuth(cfg)
	if err := l.Login(); err != nil {
		log.Fatalf("Login Err: %s", err)
	}
	fmt.Printf("Data: %#v", l)
}
