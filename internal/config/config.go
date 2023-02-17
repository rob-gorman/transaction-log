package config

// config data all hardcoded for convenience

import (
	"os"
	"time"
)

var (
	CassandraHost     string
	CassandraKeyspace string
	CassandraTable    string
	AppPort           string
	KeyTTL            time.Duration
	APIKeyTable       string
)

func InitConfig() {
	var ok bool
	if CassandraHost, ok = os.LookupEnv("CASSANDRA_HOST"); !ok {
		CassandraHost = "localhost"
	}

	if CassandraKeyspace, ok = os.LookupEnv("CASSANDRA_KEYSPACE"); !ok {
		CassandraKeyspace = "audit_log"
	}

	if CassandraTable, ok = os.LookupEnv("CASSANDRA_TABLE"); !ok {
		CassandraTable = "logs_by_time"
	}

	if AppPort, ok = os.LookupEnv("PORT"); !ok {
		AppPort = "3000"
	}

	if APIKeyTable, ok = os.LookupEnv("API_KEY_TABLE"); !ok {
		APIKeyTable = "api_keys"
	}

	KeyTTL = 7 * 24 * time.Hour
}
