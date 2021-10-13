package models

// struct similar to classes in other languages

type User struct {
	// the `json:` part allows to change how it will show to the client
	Id uint `json:"id"`
	Name string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	// change from string to []byte in order to use properly
	Password []byte `json:"-"`
}