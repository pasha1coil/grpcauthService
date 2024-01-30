package postgres

import (
	"authgrpc/internal/config"
	"authgrpc/internal/domain/models"
	"authgrpc/internal/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type Storage struct {
	conn *pgx.Conn
}

func New(dbConfig config.DBConfig) (*Storage, error) {
	const op = "storage.postgres.New"

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{conn: conn}, nil
}

func (s *Storage) Stop() error {
	return s.conn.Close(context.Background())
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, email, phone string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64
	err := s.conn.QueryRow(ctx,
		"INSERT INTO users(email, pass_hash, phone) VALUES($1, $2, $3) RETURNING id",
		email, passHash, phone).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User
	err := s.conn.QueryRow(ctx,
		"SELECT id, email,phone, pass_hash FROM users WHERE email = $1",
		email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	var app models.App
	err := s.conn.QueryRow(ctx,
		"SELECT id, name, secret FROM apps WHERE id = $1",
		id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool
	err := s.conn.QueryRow(ctx,
		"SELECT is_admin FROM users WHERE id = $1",
		userID).Scan(&isAdmin)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
