package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	secretLength = 64
)

func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist // Дошли до корня файловой системы, go.mod не найден
		}
		dir = parent
	}
}

func GetEnvFilePath() (string, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return "", fmt.Errorf("не удалось найти корень проекта: %v", err)
	}
	return filepath.Join(root, ".env"), nil
}

func CreateSecrets() {
	envFilePath, err := GetEnvFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка получения пути к .env: %v\n", err)
		os.Exit(1)
	}

	secret := make([]byte, secretLength)
	_, err = rand.Read(secret)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ошибка генерации секрета: %v\n", err)
		os.Exit(1)
	}

	hexSecret := hex.EncodeToString(secret)

	if err := writeToEnv(refreshKey, hexSecret, envFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка записи %s в .env: %v\n", refreshKey, err)
		os.Exit(1)
	}
	if err := writeToEnv(accessKey, hexSecret, envFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "ошибка записи %s в .env: %v\n", accessKey, err)
		os.Exit(1)
	}
}

func writeToEnv(key, value, envFilePath string) error {
	data, err := os.ReadFile(envFilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("ошибка чтения .env: %v", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	keyExists := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), key+"=") {
			keyExists = true
			fmt.Printf("Ключ %s уже существует в .env\n", key)
			return nil
		}
	}

	newLine := fmt.Sprintf("%s=%s", key, value)
	if keyExists {
		return nil
	}

	f, err := os.OpenFile(envFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("ошибка открытия .env: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(newLine + "\n"); err != nil {
		return fmt.Errorf("ошибка записи в .env: %v", err)
	}

	return nil
}
