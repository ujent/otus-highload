package settings

import (
	"fmt"
	"os"
)

const serverPortEnv = "APP_SERVER_PORT"
const serverPortDefault = "4000"
const dBConnStrEnv = "DB_CONN_STRING"
const saltEnv = "USER_SALT"

type Settings struct {
	Server    *ServerSettings
	DBConnStr string
}

type ServerSettings struct {
	Port string
	Salt string
}

func Load() (*Settings, error) {
	dbConnStr := os.Getenv(dBConnStrEnv)
	if dbConnStr == "" {
		return nil, fmt.Errorf("%s not set", dBConnStrEnv)
	}

	salt := os.Getenv(saltEnv)
	if salt == "" {
		return nil, fmt.Errorf("%s not set", saltEnv)
	}

	serverPort := os.Getenv(serverPortEnv)
	if serverPort == "" {
		serverPort = serverPortDefault
	}

	return &Settings{DBConnStr: dbConnStr, Server: &ServerSettings{Port: serverPort, Salt: salt}}, nil
}
