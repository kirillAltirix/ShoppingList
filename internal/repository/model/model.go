package repository

type User struct {
	UserID   int
	ChatID   string
	Username string
}

type List struct {
	ListID int
	Key    string
	Status string
}

type Item struct {
	ItemID int
	Name   string
	Status string
}
