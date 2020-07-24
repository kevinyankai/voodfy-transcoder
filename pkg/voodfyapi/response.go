package voodfyapi

// User struct used to bind user from api
type User struct {
	Token  string `json:"token"`
	Device string `json:"device"`
}

// Powergate struct used to bind powergate instance
type Powergate struct {
	InstanceID string `json:"instanceID"`
	Token      string `json:"token"`
	Address    string `json:"address"`
}

// Response struct used to bind response from api
type Response struct {
	Result struct {
		User      User      `json:"user"`
		Powergate Powergate `json:"powergate"`
	} `json:"result"`
}
