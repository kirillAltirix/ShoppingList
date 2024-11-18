package storage

type User struct {
	userID     int
	chatID     int
	usernameID string
}

type List struct {
	listID int
	key    string
	status string
}

type Iten struct {
	itemID int
	name   string
	status string
}
