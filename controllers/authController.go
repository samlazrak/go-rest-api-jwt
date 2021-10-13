package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"jwt/database"
	"jwt/models"
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

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	// since cookie is already string good to go
	// why & on the function
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error){
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	// need the standardclaims version for the issuer
	// Needs to be a POINTER! why?
	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	// the pointer was the issue with the prev version
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	// removal of cookies 'removal' - sets the exp to the past
	cookie := fiber.Cookie{
		Name: "jwt",
		Value: "",
		Expires: time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}


//TODO FIX ME
func Restricted(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error){
		return []byte(SecretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	// the pointer was the issue with the prev version
	database.DB.Where("id = ?", claims.Issuer).First(&user)

	var name = &user.Name

	return c.SendString("Welcome " + *name)
}