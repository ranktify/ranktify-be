package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"
)

type UserDAO struct {
	DB *sql.DB
}

func NewUserDAO(db *sql.DB) *UserDAO {
	return &UserDAO{DB: db}
}

func (dao *UserDAO) GetUser(email string, username string) (*model.User, error) {
	query := `
		SELECT id, username, password, first_name, last_name,
			email, role, created_at
		FROM public.users
		WHERE email = $1 OR username = $2
	`
	var user model.User
	err := dao.DB.QueryRow(query, email, username).Scan(
		&user.Id, &user.Username, &user.Password, &user.FirstName,
		&user.LastName, &user.Email, &user.Role, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	query := `
		INSERT INTO public.users (username, password, first_name, last_name, email, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id
	`

	err := dao.DB.QueryRow(query, user.Username, user.Password, user.FirstName,
		user.LastName, user.Email, user.Role,
	).Scan(&user.Id)
	if err != nil {
		return err
	}

	return nil
}
