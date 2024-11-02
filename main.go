package main

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/joho/godotenv"

	// "mime/multipart"
	"net/http"

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
func Forward(c *fiber.Ctx, rdb *redis.Client, ctx context.Context, key string,pathTodel string) error {
	// Perform a proxy forward for the DELETE request
	url := os.Getenv("DB_URL")
	err := proxy.Forward(url+key, &fasthttp.Client{
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
	})(c)
	if err != nil {
		return c.Status(505).SendString("Failed to forward request")
	}

	if pathTodel != ""{
		err = redisPkg.RemoveFromKey(rdb, ctx, pathTodel)
		if err != nil {
			return c.Status(505).SendString("Failed to remove in redis cache")
		}
	}
	return nil
}


func userApi(app *fiber.App,rdb *redis.Client,ctx context.Context) error{

	app.Get("/api/user/getAllUser",func(c *fiber.Ctx) error {
		return get(c, rdb, ctx,"user/getAllUser")
	})

	app.Post("/api/user/getUserByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"user/getUserByID","")
	})

	app.Post("/api/user/login",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"user/login","")
	})

	//not test
	app.Post("/api/user/post",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"user/post","user/getAllUser")
	})

	app.Put("/api/user/updateUserByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"user/updateUserByID","user/getAllUser")
	})

	//not test
	app.Delete("/api/user/deleteUserByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"user/deleteUserByID","user/getAllUser")
	})

	return nil
}

func petApi(app *fiber.App,rdb *redis.Client,ctx context.Context) error{

	app.Get("/api/pet/getAllPet",func(c *fiber.Ctx) error {
		return get(c, rdb, ctx,"pet/getAllPet")
	})

	app.Post("/api/pet/getPetByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"pet/getPetByID","")
	})

	app.Post("/api/pet/post",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"pet/post","pet/getAllPet")
	})

	app.Put("/api/pet/updatePetByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"pet/updatePetByID","pet/getAllPet")
	})

	app.Delete("/api/pet/deletePetByID",func(c *fiber.Ctx) error {
		return Forward(c, rdb, ctx,"pet/deletePetByID","pet/getAllPet")
	})

	return nil
}

func adoptionApi(app *fiber.App,rdb *redis.Client,ctx context.Context) error {

	app.Get("/api/adoption/getAllAdoption",func(c *fiber.Ctx) error {
		return get(c,rdb,ctx,"adoption/getAllAdoption")
	})

	app.Post("/api/adoption/getAdoptionByID",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/getAdoptionByID","")
	})

	app.Post("/api/adoption/post",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/post","adoption/getAllAdoption")
	})

	app.Put("/api/adoption/updateAdoptionByID",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/updateAdoptionByID","adoption/getAllAdoption")
	})

	app.Put("/api/adoption/adopt",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/adopt","adoption/getAllAdoption")
	})

	app.Delete("/api/adoption/delete",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/delete","adoption/getAllAdoption")
	})

	app.Get("/api/adoption/history",func(c *fiber.Ctx) error {
		return Forward(c,rdb,ctx,"adoption/history","")
	})

	return nil
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("err")
	}

	rdb, ctx := redisPkg.InitRedis()


	userApi(app,rdb,ctx)
	petApi(app,rdb,ctx)
	adoptionApi(app,rdb,ctx)

	log.Fatal(app.Listen(":4243"))
}
