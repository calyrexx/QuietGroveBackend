package handlers

type House struct {
	ID            int      `json:"id"`
	Name          string   `json:"title"`
	Description   string   `json:"description"`
	Capacity      int      `json:"people"`
	BasePrice     int      `json:"cost"`
	Images        []string `json:"images"`
	CheckInFrom   string   `json:"timeFirst"`
	CheckOutUntil string   `json:"timeSecond"`
}

type Extra struct {
	ID          int      `json:"id"`
	Name        string   `json:"title"`
	Description string   `json:"description"`
	BasePrice   int      `json:"cost"`
	Images      []string `json:"images"`
}
