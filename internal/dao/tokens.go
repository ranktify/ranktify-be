package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"
)

type TokensDAO struct {
	DB *sql.DB
}

func NewTokensDAO(db *sql.DB) *TokensDAO {
	return &TokensDAO{DB: db}
}

func (dao *TokensDAO) SaveJWTRefreshToken(jwtTokenStruct *model.JWTRefreshToken) error {
	query := `
		INSERT INTO 
			public.jwt_refresh_tokens (user_id, jti, refresh_token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	_, err := dao.DB.Exec(query,
		jwtTokenStruct.UserID,
		jwtTokenStruct.JTI,
		jwtTokenStruct.RefreshToken,
		jwtTokenStruct.ExpiresAt,
	)
	return err
}

func (dao *TokensDAO) GetJWTRefreshTokenByJTI(jti string) (*model.JWTRefreshToken, error) {
	query := `
		SELECT
			user_id, jti, refresh_token, expires_at, created_at
		FROM
			public.jwt_refresh_tokens
		WHERE
			jti = $1
	`
	row := dao.DB.QueryRow(query, jti)
	var jwtTokenStruct model.JWTRefreshToken
	err := row.Scan(
		&jwtTokenStruct.UserID,
		&jwtTokenStruct.JTI,
		&jwtTokenStruct.RefreshToken,
		&jwtTokenStruct.ExpiresAt,
		&jwtTokenStruct.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &jwtTokenStruct, nil
}

func (dao *TokensDAO) UpdateJWTRefreshTokenByJTI(oldJti string, rt *model.JWTRefreshToken) error {
	query := `
		UPDATE
			public.jwt_refresh_tokens
		SET
			user_id = $1, jti = $2, refresh_token = $3, expires_at = $4, created_at = NOW()
		WHERE
			jti = $5
		RETURNING
			user_id, jti, refresh_token, expires_at, created_at
	`

	err := dao.DB.QueryRow(query,
		rt.UserID,
		rt.JTI,
		rt.RefreshToken,
		rt.ExpiresAt,
		oldJti,
	).Scan(
		&rt.UserID,
		&rt.JTI,
		&rt.RefreshToken,
		&rt.ExpiresAt,
		&rt.CreatedAt,
	)

	return err
}

func (dao *TokensDAO) DeleteJWTRefreshTokenByJTI(jti string) error {
	query := `
		DELETE FROM
			public.jwt_refresh_tokens
		WHERE
			jti = $1
	`
	_, err := dao.DB.Exec(query, jti)
	return err
}
