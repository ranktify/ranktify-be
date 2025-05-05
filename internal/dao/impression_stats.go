package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"
)

type ImpressionDAO struct {
	DB *sql.DB
}

func NewImpressionDAO(db *sql.DB) *ImpressionDAO {
	return &ImpressionDAO{DB: db}
}

func (dao *ImpressionDAO) GetImpressionStatsByLabel(label string) (*model.ImpressionStats, error) {
	query := `
		SELECT
			impression_label,
			impressions,
			clicks,
			created_at
		FROM 
			impression_stats
		WHERE 
			impression_label = $1
		LIMIT 1 -- in case somethings bad happens
	`
	var impression model.ImpressionStats
	err := dao.DB.QueryRow(query, label).Scan(
		&impression.ImpressionLabel,
		&impression.Impressions,
		&impression.Clicks,
		&impression.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &impression, nil
}

func (dao *ImpressionDAO) UpdateImpressionStats(impression *model.ImpressionStats) error {
	query := `
		INSERT INTO impression_stats (impression_label, impressions, clicks)
		VALUES ($1, $2, $3)
		ON CONFLICT (impression_label)
		DO UPDATE
			SET impressions = impression_stats.impression + $2,
				clicks      = impression_stats.clicks + $3,
				created_at  = NOW();
	`
	_, err := dao.DB.Exec(query, impression.ImpressionLabel, impression.Impressions, impression.Clicks)
	return err
}
