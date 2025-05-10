package dao

import (
	"context"
	"database/sql"
	"log"
	"time"
)

var puertoRicoLoc *time.Location

func init() {
	var err error
	puertoRicoLoc, err = time.LoadLocation("America/Puerto_Rico")
	if err != nil {
		log.Println("Error loading Puerto Rico location:", err)
	}
}

type StreaksDAO struct {
	DB *sql.DB
}

func NewStreaksDAO(db *sql.DB) *StreaksDAO {
	return &StreaksDAO{DB: db}
}

func (dao *StreaksDAO) GetStreaksByUserID(userID uint64) (int, error) {
	query := `
		SELECT streak_count
		  FROM streaks
		 WHERE user_id = $1
	`

	var streakCount int
	err := dao.DB.QueryRow(query, userID).Scan(&streakCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return streakCount, nil
}

func (dao *StreaksDAO) RecordSongRank(ctx context.Context, userID uint64) error {
	// Start a transaction
	tx, err := dao.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Println(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	nowAST := time.Now().In(puertoRicoLoc)
	y, m, d := nowAST.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, puertoRicoLoc)

	// SELECT existing streak row FOR UPDATE
	var (
		dailyCount    int
		streakCount   int
		lastCountDate sql.NullTime
	)
	errorScan := tx.QueryRowContext(ctx, `
		SELECT daily_count, streak_count, last_count_date
		  FROM streaks
		 WHERE user_id = $1
		 FOR UPDATE
	`, userID).Scan(&dailyCount, &streakCount, &lastCountDate)

	if errorScan == sql.ErrNoRows {
		// Insert initial row
		_, err = tx.ExecContext(ctx, `
			INSERT INTO streaks(user_id, daily_count, streak_count, last_count_date, last_streak_date, updated_at)
			VALUES($1, 0, 0, NULL, NULL, $2)
		`, userID, nowAST)
		if err != nil {
			return err
		}
		dailyCount, streakCount = 0, 0
		lastCountDate.Valid = false
	} else if errorScan != nil {
		return errorScan
	}

	var isSameDay bool
	if lastCountDate.Valid {
		// Normalize to Puerto Rico time and compare Y/M/D
		ay, am, ad := lastCountDate.Time.Date()
		isSameDay = (ay == y && am == m && ad == d)
	}

	newDaily := 1
	if lastCountDate.Valid && isSameDay {
		newDaily = dailyCount + 1
	}

	// If this newDaily hits the threshold, increment streakCount
	if newDaily == 10 {
		streakCount++
		_, err = tx.ExecContext(ctx, `
			UPDATE streaks
			   SET streak_count = $2,
			       daily_count = $3,
			       last_count_date = $4,
			       updated_at = $5
			 WHERE user_id = $1
		`, userID, streakCount, newDaily, today, nowAST)
		return err
	}

	// Otherwise just update daily_count and last_count_date
	_, err = tx.ExecContext(ctx, `
		UPDATE streaks
		   SET daily_count = $2,
		       last_count_date = $3,
		       updated_at = $4
		 WHERE user_id = $1
	`, userID, newDaily, today, nowAST)
	return err
}

// ResetStreaksDaily runs in a cron job at AST midnight to clear streaks for users
// who did not reach 10 rankings the previous day, and resets daily counts for all users.
func (dao *StreaksDAO) ResetStreaksDaily(ctx context.Context) error {
	nowAST := time.Now().In(puertoRicoLoc)
	y, m, d := nowAST.Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, puertoRicoLoc)
	yesterday := today.AddDate(0, 0, -1)

	// Reset streak_count to zero for users whose last_count_date == yesterday and daily_count < 10
	_, err := dao.DB.ExecContext(ctx, `
		UPDATE streaks
		   SET streak_count = 0
		 WHERE last_count_date = $1
		   AND daily_count < 10
	`, yesterday)
	if err != nil {
		return err
	}

	// Clear daily_count and set last_count_date to today for all users
	_, err = dao.DB.ExecContext(ctx, `
		UPDATE streaks
		   SET daily_count = 0,
		       last_count_date = $1,
		       updated_at = NOW()
	`, today)
	return err
}
