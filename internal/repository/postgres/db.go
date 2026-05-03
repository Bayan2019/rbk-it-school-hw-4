package postgres

import (
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-4/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.DSN())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// var dbName string
	// err = db.Get(&dbName, "SELECT current_database()")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Database: %s\n", dbName)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
