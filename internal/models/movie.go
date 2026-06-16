package models

type Movie struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Year            int      `json:"year"`
	Genre           string   `json:"genre"`
	Format          string   `json:"format"`
	Director        string   `json:"director"`
	Cast            []string `json:"cast"`
	Synopsis        string   `json:"synopsis"`
	Rating          float64  `json:"rating"`
	RatingCount     int      `json:"rating_count"`
	Available       bool     `json:"available"`
	CopiesTotal     int      `json:"copies_total"`
	CopiesAvailable int      `json:"copies_available"`
	IsNewRelease    bool     `json:"is_new_release"`
	RentalPrice     float64  `json:"rental_price"`
	SequelTo        string   `json:"sequel_to"`
	CoverArt        string   `json:"cover_art"`
	CreatedAt       int64    `json:"created_at"`
}

type MovieResponse struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Year            int      `json:"year"`
	Genre           string   `json:"genre"`
	Format          string   `json:"format"`
	Director        string   `json:"director"`
	Cast            []string `json:"cast"`
	Synopsis        string   `json:"synopsis"`
	Rating          float64  `json:"rating"`
	RatingCount     int      `json:"rating_count"`
	Available       bool     `json:"available"`
	CopiesTotal     int      `json:"copies_total"`
	CopiesAvailable int      `json:"copies_available"`
	IsNewRelease    bool     `json:"is_new_release"`
	RentalPrice     float64  `json:"rental_price"`
	SequelTo        string   `json:"sequel_to"`
	CoverArt        string   `json:"cover_art"`
	CreatedAt       int64    `json:"created_at"`
	IsStaffPick     bool     `json:"is_staff_pick"`
}

func (m *Movie) ToResponse() MovieResponse {
	return MovieResponse{
		ID:              m.ID,
		Title:           m.Title,
		Year:            m.Year,
		Genre:           m.Genre,
		Format:          m.Format,
		Director:        m.Director,
		Cast:            m.Cast,
		Synopsis:        m.Synopsis,
		Rating:          m.Rating,
		RatingCount:     m.RatingCount,
		Available:       m.Available,
		CopiesTotal:     m.CopiesTotal,
		CopiesAvailable: m.CopiesAvailable,
		IsNewRelease:    m.IsNewRelease,
		RentalPrice:     m.RentalPrice,
		SequelTo:        m.SequelTo,
		CoverArt:        m.CoverArt,
		CreatedAt:       m.CreatedAt,
	}
}

func (m *Movie) HasCopies() bool {
	return m.CopiesAvailable > 0
}

func (m *Movie) IsLastChance() bool {
	return m.CopiesAvailable == 1 && !m.IsNewRelease
}
