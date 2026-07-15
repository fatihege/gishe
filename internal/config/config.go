package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	HTTPAddress string

	DatabaseURL string

	JWTSecret   string
	JWTIssuer   string
	JWTAudience string
}

func Load() (Config, error) {
	databaseURL, err := buildDatabaseUrl()
	if err != nil {
		return Config{}, err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is not set")
	}

	httpAddress := os.Getenv("HTTP_ADDRESS")
	if httpAddress == "" {
		httpAddress = ":8080"
	}

	return Config{
		HTTPAddress: httpAddress,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		JWTIssuer:   "gishe-api",
		JWTAudience: "gishe-client",
	}, nil
}

func buildDatabaseUrl() (string, error) {
	required := []string{
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
	}

	for _, name := range required {
		if os.Getenv(name) == "" {
			return "", fmt.Errorf("%s is not set", name)
		}
	}

	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	u := &url.URL{
		Scheme: "postgres",
		Host:   fmt.Sprintf("%s:%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		Path:   os.Getenv("POSTGRES_DB"),
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		u.User = url.User(os.Getenv("POSTGRES_USER"))
	} else {
		u.User = url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"))
	}

	query := u.Query()
	query.Set("sslmode", sslMode)
	u.RawQuery = query.Encode()

	return u.String(), nil
}
