package structs

type Desk struct {
	ID       int    `json:"id"`
	Floor    int    `json:"floor"`
	Occupied bool   `json:"occupied"`
	Body     string `json:"body"`
}
