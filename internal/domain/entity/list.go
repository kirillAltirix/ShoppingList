package entity

type Item struct {
	ID     int
	Name   string
	Status string
}

type List struct {
	ID     int
	Key    string
	Status string
	Items  []Item
}
