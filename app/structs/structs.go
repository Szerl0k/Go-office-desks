package structs

type Desk struct {
	ID       int    `json:"id"`
	Floor    int    `json:"floor"`
	Occupied bool   `json:"occupied"`
	Body     string `json:"body"`
}

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Password string `json:"password"`
}
