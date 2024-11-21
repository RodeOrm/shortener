package api

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetConfigFromFile(t *testing.T) {
	// Создаем временный конфигурационный файл для теста
	tmpFile, err := os.CreateTemp("", "config.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Пример конфигурации
	config := Config{
		ServerConfig: ServerConfig{
			ServerAddress:   "localhost:8181",
			BaseURL:         "http://localhost",
			FileStoragePath: "/path/to/file.db",
			TrustedSubnet:   "trustedSubnet",
		},
		DatabaseConfig: DatabaseConfig{
			DatabaseDSN: "dsn",
		},
		TLSConfig: TLSConfig{
			EnableHTTPS:  true,
			IsGivenHTTPS: false,
		},
	}

	// Записываем конфигурацию в файл
	data, err := json.Marshal(config)
	assert.NoError(t, err)
	_, err = tmpFile.Write(data)
	assert.NoError(t, err)

	// Тестируем SetConfigFromFile
	builder := ServerBuilder{}
	builder = builder.SetConfigFromFile(tmpFile.Name())

	// Проверка значений
	assert.Equal(t, "localhost:8181", builder.server.Config.ServerAddress)
	assert.Equal(t, "http://localhost", builder.server.Config.BaseURL)
	assert.Equal(t, "/path/to/file.db", builder.server.Config.FileStoragePath)
	assert.Equal(t, "dsn", builder.server.Config.DatabaseDSN)
	assert.Equal(t, "trustedSubnet", builder.server.Config.TrustedSubnet)
	assert.True(t, builder.server.Config.EnableHTTPS)
}

func TestSetConfig(t *testing.T) {
	builder := ServerBuilder{}
	builder = builder.SetConfig("localhost:8080", "http://localhost", "/path/to/file.db", "dsn", "true", "trustedSubnet")

	// Проверка значений
	assert.Equal(t, "localhost:8080", builder.server.Config.ServerAddress)
	assert.Equal(t, "http://localhost", builder.server.Config.BaseURL)
	assert.Equal(t, "/path/to/file.db", builder.server.Config.FileStoragePath)
	assert.Equal(t, "dsn", builder.server.Config.DatabaseDSN)
	assert.True(t, builder.server.Config.EnableHTTPS)
}

func TestSetDeleter(t *testing.T) {
	builder := ServerBuilder{}
	builder = builder.SetDeleter(5, 10, 100)

	// Проверка значений
	assert.Equal(t, 5, builder.server.Deleter.WorkerCount)
	assert.Equal(t, 10, builder.server.Deleter.BatchSize)
	assert.NotNil(t, builder.server.Deleter.DeleteQueue) // Проверяем, что очередь инициализирована
}
