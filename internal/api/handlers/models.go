package handlers

type (
	House struct {
		ID            int      `json:"id"`
		Name          string   `json:"title"`
		Description   string   `json:"description"`
		Capacity      int      `json:"people"`
		BasePrice     int      `json:"cost"`
		Images        []string `json:"images"`
		CheckInFrom   string   `json:"timeFirst"`
		CheckOutUntil string   `json:"timeSecond"`
	}

	Extra struct {
		ID          int      `json:"id"`
		Name        string   `json:"title"`
		Text        string   `json:"text"`
		Description string   `json:"description"`
		BasePrice   int      `json:"cost"`
		Images      []string `json:"images"`
	}

	GetAvailableHouses struct {
		CheckIn     string `schema:"in"`
		CheckOut    string `schema:"out"`
		GuestsCount int    `schema:"guests"`
	}

	CreateReservation struct {
		HouseID     int                    `json:"houseId"`
		Guest       Guest                  `json:"guest"`
		CheckIn     string                 `json:"checkIn"`
		CheckOut    string                 `json:"checkOut"`
		GuestsCount int                    `json:"guestsCount"`
		Extras      []ExtraReservation     `json:"extras,omitempty"`
		Bathhouse   []BathhouseReservation `json:"bathhouses,omitempty"`
	}

	BathhouseReservation struct {
		TypeID       int    `json:"id"`
		Date         string `json:"date"`
		TimeFrom     string `json:"timeFrom"`
		TimeTo       string `json:"timeTo"`
		FillOptionID int    `json:"fillId,omitempty"`
	}

	Guest struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	ExtraReservation struct {
		ID       int `json:"id"`
		Quantity int `json:"quantity"`
		Amount   int `json:"amount"`
	}

	VerifyRequest struct {
		Email string `json:"email"`
		Phone string `json:"phone"`
		Name  string `json:"name"`
	}

	Bathhouse struct {
		ID          int                   `json:"id"`
		HouseID     int                   `json:"houseId"`
		Name        string                `json:"name"`
		Price       int                   `json:"price"`
		Description string                `json:"description"`
		Images      []string              `json:"images"`
		FillOptions []BathhouseFillOption `json:"fillOptions,omitempty"`
	}

	BathhouseFillOption struct {
		ID          int    `json:"id"`
		BathhouseID int    `json:"bathhouseId"`
		Name        string `json:"name"`
		Image       string `json:"image"`
		Description string `json:"description"`
		Price       int    `json:"price"`
	}

	EventsNewApplication struct {
		Name        string `json:"name"`
		Phone       string `json:"phone"`
		CheckIn     string `json:"checkIn"`
		GuestsCount int    `json:"guestsCount"`
	}
)
