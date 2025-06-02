// package models
//
//	type User struct {
//		ID       int    `json:"id"`
//		Username string `json:"username"`
//		Password string `json:"password"`
//		Role     string `json:"role"`
//	}
package models

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Password       string `json:"password,omitempty"`
	Role           string `json:"role"`
	Nickname       string `json:"nickname,omitempty"`
	Specialization string `json:"specialization,omitempty"`
	Experience     int    `json:"experience,omitempty"`
	IsOnline       bool   `json:"is_online"`
}
type TraineeWithProfile struct {
	User
	Profile   *Profile `json:"profile,omitempty"`
	CreatedAt string   `json:"created_at"`
	LastSeen  *string  `json:"last_seen,omitempty"`
}
