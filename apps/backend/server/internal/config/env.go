package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(root string) error {
	mode := os.Getenv("APP_ENV")
	if mode == "" {
		mode = "development"
	}

	switch mode {
	case "development", "test", "production":
	default:
		return fmt.Errorf("invalid APP_ENV value: %s", mode)
	}

	if err := os.Setenv("APP_ENV", mode); err != nil {
		return fmt.Errorf("failed to set APP_ENV: %v", err)
	}

	files := []string{
		filepath.Join(root, ".env."+mode+".local"),
		filepath.Join(root, ".env."+mode),
	}

	if mode == "test" {
		files = append(files, filepath.Join(root, ".env.local"))
	}

	files = append(files, filepath.Join(root, ".env"))

	for _, file := range files {
		if err := godotenv.Load(file); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return fmt.Errorf("load environment file %s: %w", file, err)
		}
	}

	return nil
}
