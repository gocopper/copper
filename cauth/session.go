package cauth

type session struct {
	User         *user  `json:"user"`
	SessionToken string `json:"session_token"`
}
