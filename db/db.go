package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	if DB != nil {
		return nil
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := getenv("DB_HOST", "localhost")
		port := getenv("DB_PORT", "5432")
		user := getenv("DB_USER", "postgres")
		pass := getenv("DB_PASSWORD", "120311")
		name := getenv("DB_NAME", "forumdb")
		ssl := getenv("DB_SSLMODE", "disable")
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, name, ssl)
	}
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	if err := runMigrations(DB); err != nil {
		return err
	}
	fmt.Println("DB connected and migrations applied")
	return nil
}

func GetDB() *sql.DB { return DB }

func CloseDB() {
	if DB != nil {
		_ = DB.Close()
	}
}

func runMigrations(db *sql.DB) error {
	// Try likely migration paths depending on working directory
	candidates := []string{
		filepath.Join("migrations", "001_init.sql"),             // run from repo root
		filepath.Join("..", "migrations", "001_init.sql"),       // run from cmd/forum
		filepath.Join("..", "..", "migrations", "001_init.sql"), // run from deeper dirs
	}
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			bytes, err := ioutil.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read migration %s: %w", path, err)
			}
			if _, err := db.Exec(string(bytes)); err != nil {
				return fmt.Errorf("apply migration %s: %w", path, err)
			}
			return nil
		}
	}
	// fallback minimal ensures (valid syntax for PostgreSQL)
	if _, err := db.Exec(`ALTER TABLE posts ADD COLUMN IF NOT EXISTS image_data BYTEA`); err != nil {
		return fmt.Errorf("fallback migration failed: %w", err)
	}
	return nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
