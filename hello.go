package main

//std lib for input output etc.
import (
	//"fmt"
	"context"
	"fmt"
	//"errors"
	"log"
	"os"

	//for easy access of .env file
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	//fiber
	"github.com/gofiber/fiber/v2"

	//for mongo
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	//for jwt
	jtoken "github.com/golang-jwt/jwt/v4"
)

type User struct {
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

//related to middleware (figure out what is this.)
type Claims struct {
	Email string `json:"email"`
	jtoken.StandardClaims
}

var (
	mongoClient *mongo.Client
	jwtKey = []byte("") //todo: add the secret key.
)

func main(){
	err := godotenv.Load()
	if(err != nil){
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	//connect to MongoDB Atlas
	clientOptions := options.Client().ApplyURI(os.Getenv("CONNECTION_STRING"))
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		fmt.Println("error in connecting to mongo db")
		panic(err)
	}

	fmt.Println("no error in connecting to mongo db")

	//the following will get executed after the current function returns.
	defer client.Disconnect(context.Background())
	
	mongoClient = client

	//creating index for email uniqueness
	collection := mongoClient.Database("auth").Collection("users")
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "email",Value:1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil{
		fmt.Println("error in creating a collection")
		panic(err)
	}

	fmt.Println("starting new fiber app.")

	app := fiber.New()

	//routes
	app.Post("/auth/signup", SignUp)
	app.Post("/auth/signin", SignIn)

	protected := app.Group("/protected")
	protected.Use(AuthMiddleware())
	protected.Get("/", ProtectedEndpoint)

	app.Listen(port)
}

func SignUp(c *fiber.Ctx) error {
	var user User
	if err:= c.BodyParser(&user); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":"Cannot parse JSON",
		})
	}

	//Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg" : "couldn't hash password",
		})
	}

	//store user in mongodb
	collection := mongoClient.Database("auth").Collection("users")
	user.Password = string(hashedPassword)
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"msg":"Email already exsists",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"msg":"User created successfully",
	})
}

func SignIn(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"msg":"successfull"})

}

func AuthMiddleware() fiber.Handler{
	return func(c *fiber.Ctx) error{
		cookie := c.Cookies("token")
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg":"Unathorized",
			})
		}
		
		claims := &Claims{}
		token,err := jtoken.ParseWithClaims(cookie, claims, func(token *jtoken.Token) (interface{}, error ){
			return jwtKey, nil
		})

		if err != nil || !token.Valid{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg":"Unauthorized",
			})
		}

		c.Locals("email",claims.Email) // related to protected.
		return c.Next()
	}
}

func ProtectedEndpoint(c *fiber.Ctx) error {
	email := c.Locals("email").(string)
	return c.JSON(fiber.Map{
		"msg":"Welcome to protected resource",
		"email":email,
	})
}

