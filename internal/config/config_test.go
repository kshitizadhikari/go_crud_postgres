package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	keys := []string{
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "APP_PORT",
		"MINIO_ENDPOINT", "MINIO_ACCESS_KEY", "MINIO_SECRET_KEY", "MINIO_BUCKET",
	}
	for _, k := range keys {
		os.Unsetenv(k)
	}

	cfg := LoadConfig()

	if cfg.DBHost != "postgres" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "postgres")
	}
	if cfg.DBPort != "5432" {
		t.Errorf("DBPort = %q, want %q", cfg.DBPort, "5432")
	}
	if cfg.DBUser != "postgres" {
		t.Errorf("DBUser = %q, want %q", cfg.DBUser, "postgres")
	}
	if cfg.DBPassword != "password" {
		t.Errorf("DBPassword = %q, want %q", cfg.DBPassword, "password")
	}
	if cfg.DBName != "go_crud" {
		t.Errorf("DBName = %q, want %q", cfg.DBName, "go_crud")
	}
	if cfg.DBSSLMode != "disable" {
		t.Errorf("DBSSLMode = %q, want %q", cfg.DBSSLMode, "disable")
	}
	if cfg.AppPort != "8080" {
		t.Errorf("AppPort = %q, want %q", cfg.AppPort, "8080")
	}
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	os.Setenv("DB_HOST", "customhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("APP_PORT", "9090")
	os.Setenv("MINIO_ENDPOINT", "minio.custom.com:9000")
	os.Setenv("MINIO_BUCKET", "my-bucket")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("MINIO_ENDPOINT")
		os.Unsetenv("MINIO_BUCKET")
	}()

	cfg := LoadConfig()

	if cfg.DBHost != "customhost" {
		t.Errorf("DBHost = %q, want %q", cfg.DBHost, "customhost")
	}
	if cfg.DBPort != "5433" {
		t.Errorf("DBPort = %q, want %q", cfg.DBPort, "5433")
	}
	if cfg.AppPort != "9090" {
		t.Errorf("AppPort = %q, want %q", cfg.AppPort, "9090")
	}
	if cfg.Minio.Endpoint != "minio.custom.com:9000" {
		t.Errorf("Minio.Endpoint = %q, want %q", cfg.Minio.Endpoint, "minio.custom.com:9000")
	}
	if cfg.Minio.Bucket != "my-bucket" {
		t.Errorf("Minio.Bucket = %q, want %q", cfg.Minio.Bucket, "my-bucket")
	}
}

func TestGetDSN(t *testing.T) {
	cfg := &Config{
		DBHost:     "db.example.com",
		DBUser:     "admin",
		DBPassword: "secret",
		DBName:     "mydb",
		DBPort:     "5432",
		DBSSLMode:  "require",
	}

	dsn := cfg.GetDSN()

	if !strings.Contains(dsn, "host=db.example.com") {
		t.Errorf("DSN missing host, got: %s", dsn)
	}
	if !strings.Contains(dsn, "user=admin") {
		t.Errorf("DSN missing user, got: %s", dsn)
	}
	if !strings.Contains(dsn, "password=secret") {
		t.Errorf("DSN missing password, got: %s", dsn)
	}
	if !strings.Contains(dsn, "dbname=mydb") {
		t.Errorf("DSN missing dbname, got: %s", dsn)
	}
	if !strings.Contains(dsn, "port=5432") {
		t.Errorf("DSN missing port, got: %s", dsn)
	}
	if !strings.Contains(dsn, "sslmode=require") {
		t.Errorf("DSN missing sslmode, got: %s", dsn)
	}
	if !strings.Contains(dsn, "TimeZone=Asia/Kathmandu") {
		t.Errorf("DSN missing timezone, got: %s", dsn)
	}
}

func TestMinioConfig_Defaults(t *testing.T) {
	os.Unsetenv("MINIO_ENDPOINT")
	os.Unsetenv("MINIO_ACCESS_KEY")
	os.Unsetenv("MINIO_SECRET_KEY")
	os.Unsetenv("MINIO_BUCKET")

	cfg := LoadConfig()

	if cfg.Minio.Endpoint != "localhost:9000" {
		t.Errorf("Minio.Endpoint = %q, want %q", cfg.Minio.Endpoint, "localhost:9000")
	}
	if cfg.Minio.AccessKey != "admin" {
		t.Errorf("Minio.AccessKey = %q, want %q", cfg.Minio.AccessKey, "admin")
	}
	if cfg.Minio.SecretKey != "admin123" {
		t.Errorf("Minio.SecretKey = %q, want %q", cfg.Minio.SecretKey, "admin123")
	}
	if cfg.Minio.UseSSL != false {
		t.Errorf("Minio.UseSSL = %v, want false", cfg.Minio.UseSSL)
	}
	if cfg.Minio.Bucket != "uploads" {
		t.Errorf("Minio.Bucket = %q, want %q", cfg.Minio.Bucket, "uploads")
	}
}

func TestGetEnv_Helper(t *testing.T) {
	os.Setenv("TEST_ENV_VAR", "hello")
	defer os.Unsetenv("TEST_ENV_VAR")

	val := getEnv("TEST_ENV_VAR", "default")
	if val != "hello" {
		t.Errorf("getEnv() = %q, want %q", val, "hello")
	}

	val = getEnv("NONEXISTENT_VAR_12345", "default")
	if val != "default" {
		t.Errorf("getEnv() = %q, want %q", val, "default")
	}
}
