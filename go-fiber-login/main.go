package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type User struct {
	Realname string `json:"realname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func sha1Hash(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	app := fiber.New()

	app.Post("/login", func(c *fiber.Ctx) error {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}

		val, err := rdb.Get(ctx, "login_"+body.Username).Result()
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "User not found"})
		}

		var user User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Data error"})
		}

		if user.Password != sha1Hash(body.Password) {
			return c.Status(401).JSON(fiber.Map{"error": "Wrong password"})
		}

		return c.JSON(fiber.Map{
			"message":  "Login success",
			"realname": user.Realname,
			"email":    user.Email,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
