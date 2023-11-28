package routes

import (
	"faseflow-api/model"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type user struct{
	Username string `json:"username"`	
	Password string `json:"password"`
}

func NewMiddleware() fiber.Handler {
 return AuthMiddleware
}

func AuthMiddleware(c *fiber.Ctx) error {

	sess, err := store.Get(c)

	if strings.Split(c.Path(), "/")[1] == "auth" {
		return c.Next()
	}
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message":"not authorized",
		})
	}

	if sess.Get(AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message":"not authorized",
		})
	}
	return c.Next()
}

func Register(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var data user

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong: " + err.Error(),
				})
	}
	password, bcErr := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if bcErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong: " + err.Error(),
				})
	}
	user := model.User {
		Username: data.Username,		
		Password: string(password),
	}
	err = model.CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong: " + err.Error(),
				})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "registered",
	})
}

func Login(c *fiber.Ctx) error {
	var data user
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong: " + err.Error(),
				})
	}
	var user model.User
	if !model.CheckUser(data.Username, &user) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized. username not in db",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized. wrong password",
		})
	}
	sess, sessErr := store.Get(c)
	if sessErr != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong " + err.Error(),
		})
	}
	sess.Set(AUTH_KEY, true)
	sess.Set(USER_ID, user.ID)
	sessErr = sess.Save()
	if sessErr != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map {
		"message" : "logged in",
	})
}

func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map {
			"message" : "logged out (no session)",
		})
	}
	err = sess.Destroy()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map {
			"message" : "something went wrong " + err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map {
		"message" : "logged out",
	})
}

func HealthCheck(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
	auth := sess.Get(AUTH_KEY)

	if auth != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map {
			"message": "authenticated",
		})
	}else{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
}

func GetUser(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
	if sess.Get(AUTH_KEY) == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
	userId := sess.Get(USER_ID)
	if userId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
	var user model.User
	
	user, err = model.GetUser(fmt.Sprint(userId))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map {
			"message" : "not authorized",
		})
	}
	user.Password = ""
	return c.Status(fiber.StatusOK).JSON(user)
}