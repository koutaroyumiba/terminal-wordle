package game

import (
	"encoding/json"
	"os"
)

const statsFile = "stats.json"

type Stats struct {
	GamesPlayed    int         `json:"games_played"`
	Wins           int         `json:"wins"`
	CurrentStreak  int         `json:"current_streak"`
	MaxStreak      int         `json:"max_streak"`
	GuessFrequency map[int]int `json:"guess_frequency"`
}

func (s Stats) WinRate() float64 {
	if s.GamesPlayed == 0 {
		return 0
	}

	return float64(s.Wins) / float64(s.GamesPlayed) * 100
}

func (s Stats) AverageGuesses() float64 {
	if s.Wins == 0 {
		return 0
	}

	total := 0
	for key, val := range s.GuessFrequency {
		total += (key * val)
	}

	return float64(total) / float64(s.Wins)
}

func loadStats() Stats {
	f, err := os.ReadFile(statsFile)
	if err != nil {
		return Stats{
			GuessFrequency: make(map[int]int),
		}
	}

	var s Stats
	if err := json.Unmarshal(f, &s); err != nil {
		return Stats{
			GuessFrequency: make(map[int]int),
		}
	}

	return s
}

func saveStats(s Stats) {
	data, _ := json.Marshal(s)
	_ = os.WriteFile(statsFile, data, 0644)
}
