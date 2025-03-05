package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"

	"fmt"
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

func (dao *UserDAO) GetUserByID(id uint64) (*model.User, error) {
	query := `
		SELECT id, username, password, first_name, last_name, email
		FROM public.users
		WHERE id = $1
	`
	var user model.User
	err := dao.DB.QueryRow(query, id).Scan(
		&user.Id, &user.Username, &user.Password, &user.FirstName,
		&user.LastName, &user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAO) GetAllUsers() ([]*model.User, error) {
	query := `
		SELECT id, username, password, first_name, last_name, email
		FROM public.users
	`
	rows, err := dao.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.FirstName,
			&user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAO) UpdateUserByID(id uint64, user *model.User) error {
	query := `
		UPDATE public.users
		SET username = $1, password = $2, first_name = $3, last_name = $4, email = $5
		WHERE id = $6
	`
	result, err := dao.DB.Exec(query, user.Username, user.Password, user.FirstName, user.LastName, user.Email, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (dao *UserDAO) DeleteUserByID(id uint64) error {
	query := `
		DELETE FROM public.users
		WHERE id = $1
	`
	result, err := dao.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
