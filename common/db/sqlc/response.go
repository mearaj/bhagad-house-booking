package sqlc

type BookingsResponse struct {
	Bookings []Booking `json:"bookings,omitempty"`
	Error    string    `json:"error,omitempty"`
}

type AuthErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type CreateBookingResponse struct {
	Booking Booking `json:"booking,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type UpdateBookingResponse struct {
	Booking Booking `json:"booking,omitempty"`
	Error   string  `json:"error,omitempty"`
}

type DeleteBookingResponse struct {
	ID    int64  `json:"ID,omitempty"`
	Error string `json:"error,omitempty"`
}

type NewUserResponse struct {
	User  User   `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type UserResponse struct {
	User  User   `json:"user,omitempty"`
	Error string `json:"error,omitempty"`
}

type LoginUserResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	User        User   `json:"user,omitempty"`
	Error       string `json:"error,omitempty"`
}

type UsersResponse struct {
	Users []User `json:"users,omitempty"`
	Error string `json:"error,omitempty"`
}
