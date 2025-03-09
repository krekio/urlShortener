package postgres

import (
	"database/sql"
	"github.com/krekio/urlShortener.git/internal/config"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
	"log"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg config.Config) (*Storage, error) {
	// Открываем соединение с базой данных

	db, err := sql.Open("postgres", cfg.StorageDsn)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение к базе данных
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создаем таблицу, если она не существует
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS url (
            id SERIAL PRIMARY KEY,
            alias TEXT UNIQUE NOT NULL,
            url TEXT NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	// Создаем индекс, если он не существует
	_, err = db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
    `)
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established and tables are ready")
	return &Storage{DB: db}, nil
}
