package pages

import (
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/tui/styles"
)

type BrowseModel struct {
	movies   []models.MovieResponse
	selected int
	page     int
	total    int
}

func NewBrowseModel() *BrowseModel {
	return &BrowseModel{
		page: 1,
	}
}

func (m *BrowseModel) SetMovies(movies []models.MovieResponse, total int) {
	m.movies = movies
	m.total = total
}

func (m *BrowseModel) View(width, height int) string {
	if len(m.movies) == 0 {
		return styles.TextStyle.Render("Loading catalog...")
	}
	return styles.TextStyle.Render("Browse catalog — Phase 6")
}
