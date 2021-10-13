package models

// struct similar to classes in other languages

type User struct {
	Id uint
	Name string
	Email string `gorm:"unique"`
	// change from string to []byte in order to use properly
	Password []byte
}