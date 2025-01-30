package structs

import (
	"regexp"
)

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
	Role     string `json:"role"`
}

func (u *User) ValidEmail() bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(u.Email)
}
