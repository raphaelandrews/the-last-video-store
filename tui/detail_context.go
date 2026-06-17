package tui

import "github.com/thelastvideostore/internal/models"

func (m *Model) setDetailContext() {
	if m.detail == nil || m.userResp == nil {
		return
	}
	tier := models.TierByName(m.userResp.Subscription)
	m.detail.SetUserContext(m.userResp.FreeRentals, tier.FreeRentals, m.userResp.Balance)
	if m.detail.Movie.SequelTo != "" {
		for _, mv := range m.browse.Movies {
			if mv.ID == m.detail.Movie.SequelTo {
				m.detail.SequelTitle = mv.Title
				break
			}
		}
	}
	var franchise []models.MovieResponse
	currentID := m.detail.Movie.ID
	seen := map[string]bool{currentID: true}
	id := m.detail.Movie.SequelTo
	for id != "" && !seen[id] {
		found := false
		for _, mv := range m.browse.Movies {
			if mv.ID == id {
				seen[id] = true
				franchise = append([]models.MovieResponse{mv}, franchise...)
				id = mv.SequelTo
				found = true
				break
			}
		}
		if !found {
			break
		}
	}
	franchise = append(franchise, *m.detail.Movie)
	queue := []string{currentID}
	for len(queue) > 0 {
		prequelID := queue[0]
		queue = queue[1:]
		for _, mv := range m.browse.Movies {
			if mv.SequelTo == prequelID && !seen[mv.ID] {
				seen[mv.ID] = true
				franchise = append(franchise, mv)
				queue = append(queue, mv.ID)
			}
		}
	}
	if len(franchise) > 1 {
		m.detail.Franchise = franchise
	} else {
		m.detail.Franchise = nil
	}

	var sameGenre []models.MovieResponse
	for _, mv := range m.browse.Movies {
		if !seen[mv.ID] && mv.Genre == m.detail.Movie.Genre {
			sameGenre = append(sameGenre, mv)
			if len(sameGenre) >= 5 {
				break
			}
		}
	}
	m.detail.Recommendations = sameGenre
}
