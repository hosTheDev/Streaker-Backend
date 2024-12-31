package main

//std lib for input output etc.
import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"


	//"fiber"
	"github.com/gofiber/fiber/v2"
)

func main(){
	err := godotenv.Load()
	if(err != nil){
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

}

//start the field names with a capital letter. 
type User struct {
	Email string
	Password string
}
