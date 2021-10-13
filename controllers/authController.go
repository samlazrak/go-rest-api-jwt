package controllers

import (
	"github.com/gofiber/fiber/v2"
	"jwt/models"
	"golang.org/x/crypto/bcrypt"
	"jwt/database"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const SecretKey = "secret"

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

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		// needs to be converted to string ( first int bc uint can't be used )
		Issuer: strconv.Itoa(int(user.Id)),
		// Conver to unix for general usablility?
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return c.JSON(fiber.Map{
			"message": "Not logged in",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24),
		// Bc front end doesn't need to access?
		HTTPOnly: true,
	}

	// Before &: Cannot use 'cookie' (type Cookie) as the type *Cookie
	// We do not actually want to use the actual data point of cookie, as it cannot intake that form of data
	// Therefore we pass reference to the data - essentially imagine it's key name vs key value
	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}