package data

// Email data class
type Email struct {
	Domain        string
	APIKeyPrivate string
	APIKeyPublic  string
	From          string
	Subject       string
	Body          string
	HTML          string
	To            string
}
