package cmd

import (
	"fmt"
	"os"
)

// SaveTokenToFile - Функция для сохранения токена в файл
func SaveTokenToFile(token string) error {
	// Открытие файла для записи (если файла нет, он будет создан)
	file, err := os.Create("token")
	if err != nil {
		return fmt.Errorf("не удалось создать файл: %w", err)
	}
	defer file.Close()

	// Запись токена в файл
	_, err = file.WriteString(token)
	if err != nil {
		return fmt.Errorf("не удалось записать токен в файл: %w", err)
	}

	return nil
}

// ReadTokenFromFile - Функция для чтения токена из файла
func ReadTokenFromFile() (string, error) {
	// Открытие файла для чтения
	file, err := os.Open("token")
	if err != nil {
		return "", fmt.Errorf("не удалось открыть файл: %w", err)
	}
	defer file.Close()

	// Чтение токена из файла
	var token string
	_, err = fmt.Fscanf(file, "%s", &token)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать токен из файла: %w", err)
	}

	return token, nil
}
