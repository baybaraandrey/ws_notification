package postgresql

import (
	"database/sql"
	"fmt"
	nativeLog "log"

	_ "github.com/lib/pq"

	config "github.com/baybaraandrey/ws_notification/internal/config"
)

// OpenConnection open new connection to db
func OpenConnection() *sql.DB {
	cfg, err := config.Load()
	if err != nil {
		nativeLog.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s", cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPass, cfg.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		nativeLog.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		nativeLog.Fatal(err)
	}

	return db
}
