package models

const (
	FormatVHS    = "VHS"
	FormatDVD    = "DVD"
	FormatBluRay = "Blu-ray"

	VHSLateFeeRate = 2.00
	DVDLateFeeRate = 3.00
	RewindFeeCost  = 1.00
)

var GenreList = []string{"Action", "Comedy", "Horror", "SciFi", "Drama", "Thriller", "Romance", "Animation"}

var SeriesGenreList = []string{"Crime", "Animation", "Drama", "SciFi", "Comedy", "Fantasy", "Horror", "Thriller", "Documentary"}

var GameGenreList = []string{"Action", "RPG", "Racing", "Platformer", "FPS", "Strategy", "Fighting", "Puzzle", "Sports"}
