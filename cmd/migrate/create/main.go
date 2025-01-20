package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run create.go <migration_name>")
		return
	}

	migrationName := os.Args[1]
	timestamp := time.Now().Format("20060102150405")
	basePath := "cmd/migrate/migrations"
	upFile := filepath.Join(basePath, fmt.Sprintf("%s_%s.up.sql", timestamp, migrationName))
	downFile := filepath.Join(basePath, fmt.Sprintf("%s_%s.down.sql", timestamp, migrationName))

	// Ensure the migrations directory exists
	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		fmt.Printf("Error creating migrations directory: %v\n", err)
		return
	}

	// Create the .up.sql file
	if err := os.WriteFile(upFile, []byte("-- Add your SQL migration up script here\n"), os.ModePerm); err != nil {
		fmt.Printf("Error creating file %s: %v\n", upFile, err)
		return
	}

	// Create the .down.sql file
	if err := os.WriteFile(downFile, []byte("-- Add your SQL migration down script here\n"), os.ModePerm); err != nil {
		fmt.Printf("Error creating file %s: %v\n", downFile, err)
		return
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
}
