package voodfyapi

// User struct used to bind user from api
type User struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}

// Response struct used to bind response from api
type Response struct {
	Result struct {
		User User `json:"user"`
	} `json:"result"`
}
