package models

type GameSession struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	GameID    string  `json:"game_id"`
	GameTitle string  `json:"game_title"`
	StartedAt int64   `json:"started_at"`
	EndedAt   int64   `json:"ended_at"`
	Duration  int64   `json:"duration"`
	Cost      float64 `json:"cost"`
	Status    string  `json:"status"`
}
