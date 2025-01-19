package internal

type Credential struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Metadata string `json:"metadata"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
