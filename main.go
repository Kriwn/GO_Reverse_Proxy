package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"

	// "mime/multipart"
	"net/http"
	"strings"

	"github.com/Kriwn/Go_Reverse_Proxy/RedisPkg"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

func get(c *fiber.Ctx, rdb *redis.Client, ctx context.Context, pathParam string) error {
	key := redisPkg.GetValueFromKey(rdb, ctx, pathParam)
	if key.Err() != nil {
		path := os.Getenv("DB_URL")
		url := path + pathParam
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return c.Status(505).SendString("Failed to create request")
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return c.Status(505).SendString("Failed to perform request")
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return c.Status(505).SendString("Failed to read response body")
		}

		// Cache the response in Redis
		err = redisPkg.SetNew(rdb, ctx, pathParam, string(body))
		if err != nil {
			return c.Status(505).SendString("Failed to cache data")
		}

		return c.SendString(string(body))
	}
	// Return cached value from Redis
	return c.SendString(key.Val())
}

// nedd to add user id in body to pass proxy to delete
func Forward(c *fiber.Ctx, rdb *redis.Client, ctx context.Context, key string) error {
	// Perform a proxy forward for the DELETE request
	url := os.Getenv("DB_URL")
	err := proxy.Forward(url+key, &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})(c)
	if err != nil {
		return c.Status(505).SendString("Failed to forward request")
	}

	//Remove the key from Redis cache
	newString := strings.Split(key, "/")
	err = redisPkg.RemoveFromKey(rdb, ctx, newString[0]+"/get")
	if err != nil {
		return c.Status(505).SendString("Failed to remove in redis cache")
	}

	return nil
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowCredentials: true,

	}))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err")
	}

	rdb, ctx := redisPkg.InitRedis()

	app.Get("/api/*", func(c *fiber.Ctx) error {
		pathParam := c.Params("*")
		if pathParam != "" {
			newString := strings.Split(pathParam, "/")
			check := newString[len(newString)-1]

			if check == "get" {
				return get(c, rdb, ctx, pathParam)
			}
			return c.Status(400).SendString("Invalid endpoint")
		}

		return c.Status(400).SendString("Path parameter missing")
	})

	app.Delete("/api/*", func(c *fiber.Ctx) error {
		pathParam := c.Params("*")
		if pathParam != "" {
			newString := strings.Split(pathParam, "/")
			check := newString[len(newString)-1]

			if check == "delete" {
				return Forward(c, rdb, ctx, pathParam)
			}
			return c.Status(400).SendString("Invalid endpoint")
		}

		return c.Status(400).SendString("Path parameter missing")
	})

	app.Post("/api/*", func(c *fiber.Ctx) error {
		pathParam := c.Params("*")

		if pathParam != "" {
			newString := strings.Split(pathParam, "/")
			check := newString[len(newString)-1]

			if check == "post" {
				return Forward(c, rdb, ctx, pathParam)
			}
			return c.Status(400).SendString("Invalid endpoint")
		}

		return c.Status(400).SendString("Path parameter missing")
	})

	app.Post("/api/*", func(c *fiber.Ctx) error {
		pathParam := c.Params("*")

		if pathParam != "" {
			newString := strings.Split(pathParam, "/")
			check := newString[len(newString)-1]

			if check == "post" {
				return Forward(c, rdb, ctx, pathParam)
			}
			return c.Status(400).SendString("Invalid endpoint")
		}

		return c.Status(400).SendString("Path parameter missing")
	})

	app.Put("/api/*", func(c *fiber.Ctx) error {
		pathParam := c.Params("*")

		if pathParam != "" {
			newString := strings.Split(pathParam, "/")
			check := newString[len(newString)-1]

			if check == "put" {
				return Forward(c, rdb, ctx, pathParam)
			}
			return c.Status(400).SendString("Invalid endpoint")
		}

		return c.Status(400).SendString("Path parameter missing")
	})

	log.Fatal(app.Listen(":4243"))
}
