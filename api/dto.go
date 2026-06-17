package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TOTPLoginRequest struct {
	Code string `json:"code"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Tier     string `json:"tier,omitempty"`
}

type CreateMovieRequest struct {
	Title        string   `json:"title"`
	Year         int      `json:"year"`
	Genre        string   `json:"genre"`
	Format       string   `json:"format"`
	Director     string   `json:"director"`
	Cast         []string `json:"cast"`
	Synopsis     string   `json:"synopsis"`
	CopiesTotal  int      `json:"copies_total"`
	IsNewRelease bool     `json:"is_new_release"`
}

type UpdateMovieRequest struct {
	Title        *string   `json:"title,omitempty"`
	Year         *int      `json:"year,omitempty"`
	Genre        *string   `json:"genre,omitempty"`
	Format       *string   `json:"format,omitempty"`
	Director     *string   `json:"director,omitempty"`
	Cast         *[]string `json:"cast,omitempty"`
	Synopsis     *string   `json:"synopsis,omitempty"`
	CopiesTotal  *int      `json:"copies_total,omitempty"`
	IsNewRelease *bool     `json:"is_new_release,omitempty"`
}

type RentRequest struct {
	MovieID   string `json:"movie_id"`
	UseTicket bool   `json:"use_ticket"`
}

type ReturnRequest struct {
	RentalID string `json:"rental_id"`
}

type UpdateUserRequest struct {
	Tier   *string `json:"tier,omitempty"`
	Banned *bool   `json:"banned,omitempty"`
}

type TOTPRequest struct {
	Enabled bool `json:"enabled"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    int64       `json:"expires_at"`
	User         interface{} `json:"user"`
}

type MovieListResponse struct {
	Movies   []interface{} `json:"movies"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type StaffPickResponse struct {
	StaffPick bool `json:"staff_pick"`
}

type TOTPSetupResponse struct {
	Secret string `json:"secret"`
	URL    string `json:"url"`
}
