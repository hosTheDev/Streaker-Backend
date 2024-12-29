package main

//std lib for input output etc.
import (
	"fmt"
	"log"

	//"fiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	//print statments
	fmt.Println("hello from go")

	//to make a new instance of Fiber.
	app := fiber.New()

	app.Get("/api/get", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"msg": "Habibi"})
	})

	todos := []Todo{}

	app.Post("/api/post", func(c *fiber.Ctx) error {
		todo := &Todo{}

		//not working always returning body as "" .. THIS IS NOT THE REASON .. CHECK THE COMMENT OVER Todo struct.
		err := c.BodyParser(todo)
		
		if err != nil {
			log.Fatal(err)
		}

		// if err := c.BodyParser(todo); err != nil{
		// 	return err
		// }

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body is required."})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(200).JSON(todo)
	})

	//to catch any errors
	log.Fatal(app.Listen(":4000"))
}

//start the field names with a capital letter. 
type Todo struct {
	Id   int
	Task string `json:"task"`
	Done bool   `json:"done"`
	Body string `json:"body"`
}
