package auth0

type Config struct {
	AcceptableTimeSkewSeconds int    `json:"acceptable_time_skew"`
	Audience                  string `json:"audience"`
	Domain                    string `json:"domain"`
	ClientID                  string `json:"client_id"`
	ClientSecret              string `json:"client_secret"`
}
