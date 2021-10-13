package controllers

import (
	"github.com/gofiber/fiber/v2"
	"jwt/models"
	"golang.org/x/crypto/bcrypt"
	"jwt/database"
)

func Register(c *fiber.Ctx) error {
	// similar to array as string as key and string as value
	var data map[string]string

	// Passing as reference &
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// since won't accept string []byte converts to slice of bytes - essentially a dynamic array no prealoc size
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: password,
	}

	// Again passing as reference what this mean
	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	// the .First(&user) assigns to above user variable. why using & - getting it's memory address.
	// * goes in front of var that holds a mem adr and resolves it. It goes and gets the thing that the pointer is pointing to.
	// When * is put in front of a type, e.g. *string, it becomes part of the type declaration, so you can say "this variable holds a pointer to a string". For example: var str_pointer *string
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		// map is string interface?
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}
	// Why no & ?
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	return c.JSON(user)
}