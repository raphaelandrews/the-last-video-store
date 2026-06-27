package models

const (
	RentalActive   = "active"
	RentalReturned = "returned"
	RentalOverdue  = "overdue"
)

type Rental struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	MovieID      string  `json:"movie_id"`
	MovieFormat  string  `json:"movie_format"`
	RentedAt     int64   `json:"rented_at"`
	DueDate      int64   `json:"due_date"`
	ReturnedAt   int64   `json:"returned_at"`
	LateFee      float64 `json:"late_fee"`
	RewindFee    float64 `json:"rewind_fee"`
	NeedsRewind  bool    `json:"needs_rewind"`
	Status       string  `json:"status"`
	IsFreeRental bool    `json:"is_free_rental"`
}

type RentalResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	MovieID      string  `json:"movie_id"`
	MovieTitle   string  `json:"movie_title"`
	MovieFormat  string  `json:"movie_format"`
	RentedAt     int64   `json:"rented_at"`
	DueDate      int64   `json:"due_date"`
	ReturnedAt   int64   `json:"returned_at"`
	LateFee      float64 `json:"late_fee"`
	RewindFee    float64 `json:"rewind_fee"`
	NeedsRewind  bool    `json:"needs_rewind"`
	Status       string  `json:"status"`
	IsFreeRental bool    `json:"is_free_rental"`
	// PointsEarned is the change in 🍿 Popcorn Points from this
	// return: +10 on time, -5 if late, +5 bonus if the user owns the
	// Popcorn Bucket. Set by the return handler; 0 for active rentals
	// or rentals never returned through the API.
	PointsEarned int `json:"points_earned"`
}

func (r *Rental) ToResponse(movieTitle string) RentalResponse {
	return RentalResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		MovieID:      r.MovieID,
		MovieTitle:   movieTitle,
		MovieFormat:  r.MovieFormat,
		RentedAt:     r.RentedAt,
		DueDate:      r.DueDate,
		ReturnedAt:   r.ReturnedAt,
		LateFee:      r.LateFee,
		RewindFee:    r.RewindFee,
		NeedsRewind:  r.NeedsRewind,
		Status:       r.Status,
		IsFreeRental: r.IsFreeRental,
	}
}

func (r *Rental) IsOverdue(now int64) bool {
	return r.Status != RentalReturned && now > r.DueDate
}

func (r *Rental) CalculateLateFee(now int64) float64 {
	if now <= r.DueDate || r.IsFreeRental {
		return 0
	}
	const minute = int64(60)
	minutesLate := (now - r.DueDate) / minute
	if minutesLate < 1 {
		minutesLate = 1
	}
	rate := DailyLateFeeRate(r.MovieFormat) / 10 // per-minute rate
	return float64(minutesLate) * rate
}

func (r *Rental) CalculateRewindFee() float64 {
	if r.NeedsRewind && r.MovieFormat == "VHS" {
		return 1.00
	}
	return 0
}

func (r *Rental) TotalFee() float64 {
	return r.LateFee + r.RewindFee
}

func DueDateForFormat(format string, rentedAt int64) int64 {
	const minute = int64(60)
	switch format {
	case "VHS":
		return rentedAt + 1*minute
	case "DVD":
		return rentedAt + 2*minute
	case "Blu-ray":
		return rentedAt + 3*minute
	default:
		return rentedAt + 4*minute
	}
}

func DailyLateFeeRate(format string) float64 {
	switch format {
	case "VHS":
		return 2.00
	case "DVD", "Blu-ray":
		return 3.00
	default:
		return 3.00
	}
}
