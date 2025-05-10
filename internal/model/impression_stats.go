package model

type ImpressionStats struct {
	Id              uint64 `json:"id"`
	ImpressionLabel string `json:"impression_label"`
	Impressions     uint64 `json:"impressions"`
	Clicks          uint64 `json:"clicks"`
	CreatedAt       string `json:"created_at"`
}

func (I *ImpressionStats) CTRPercent() float64 {
	if I.Impressions == 0 {
		return 0
	}
	return float64(I.Clicks) / float64(I.Impressions) * 100
}
