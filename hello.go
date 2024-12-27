package main

//std lib for input output etc.
import (
	"fmt"
	"log"
	//"fiber"
	"github.com/gofiber/fiber/v2"
)

func main(){
	//print statments
	fmt.Println("hello from go")

	//to make a new instance of Fiber.
	app := fiber.New()

	app.Get("/get", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg":"Habibi"})
	})

	//to catch any errors
	log.Fatal(app.Listen(":4000"))
}