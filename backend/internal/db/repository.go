package db

import (
	"database/sql"
)

type UserRepository struct {
	DB *sql.DB
}

var Repo *UserRepository

func InitRepo(database *sql.DB) {
	Repo = &UserRepository{DB: database}
}

// Проверка существования пользователя
func (r *UserRepository) Exists(telegramID int64) (bool, error) {
	var exists bool
	err := r.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE telegram_id = $1)",
		telegramID,
	).Scan(&exists)
	return exists, err
}

// Сохранение нового пользователя
func (r *UserRepository) CreateUser(u User) error {
	_, err := r.DB.Exec(`
        INSERT INTO users (telegram_id, username, role, first_name, last_name, created_at)
        VALUES ($1, $2, $3, $4, $5, NOW())`,
		u.ID, u.Username, u.Role, u.FirstName, u.LastName,
	)
	return err
}

// Обновление роли пользователя по Telegram ID
func (r *UserRepository) UpdateUserRole(telegramID int64, newRole string) error {
	_, err := r.DB.Exec(
		"UPDATE users SET role = $1 WHERE telegram_id = $2",
		newRole, telegramID,
	)
	return err
}

// Удаление пользователя по Telegram ID
func (r *UserRepository) DeleteUser(telegramID int64) (bool, error) {
	res, err := r.DB.Exec(
		"DELETE FROM users WHERE telegram_id = $1",
		telegramID,
	)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rows > 0, nil
}

func (r *UserRepository) GetAll() ([]User, error) {
	rows, err := r.DB.Query(`
    SELECT 
        id,
        telegram_id,
        role,
        username,
        first_name,
        COALESCE(last_name, ''),
        COALESCE(email, ''),
        COALESCE(phone, ''),
        created_at
    FROM users
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.TelegramID,
			&u.Role,
			&u.Username,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}
