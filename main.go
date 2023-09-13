package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pusher/pusher-http-go/v5"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/primitive"
)

var client *mongo.Client
var usersCollection *mongo.Collection
var channelsCollection *mongo.Collection
var userChannelCollection *mongo.Collection
var messagesCollection *mongo.Collection

func main() {
	app := fiber.New()

	app.Use(cors.New())

	pusherClient := pusher.Client{
		AppID:   "1648273",
		Key:     "c6998e24f5996790c71d",
		Secret:  "77a1884f6afc83c9dc72",
		Cluster: "ap1",
		Secure:  true,
	}

	// Connect to mongo
	uri := "mongodb+srv://admin:pa55word@cluster0.nfhn3fr.mongodb.net/?retryWrites=true&w=majority"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("chatdb").Collection("chatlog")
	usersCollection = client.Database("chatdb").Collection("users")
	channelsCollection = client.Database("chatdb").Collection("channels")
	userChannelCollection = client.Database("chatdb").Collection("user-channel")
	messagesCollection = client.Database("chatdb").Collection("messages")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	var results []bson.D
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to your MongoDB database!")

	// Messaging Routes

	app.Post("/api/messages", func(c *fiber.Ctx) error {
		// var data map[string]string
		var messageData struct {
			Username string `json:"username"`
			Message  string `json:"message"`
		}

		if err := c.BodyParser(&messageData); err != nil {
			return err
		}

		fmt.Println("Received data:", messageData)

		// Generate a BSON UTC datetime
		timestamp := primitive.NewDateTimeFromTime(time.Now().UTC())

		data := map[string]interface{}{
			"chat_timestamp": timestamp,
			"username":       messageData.Username,
			"message":        messageData.Message,
		}

		_, err := collection.InsertOne(context.Background(), data)
		if err != nil {
			return fmt.Errorf("failed to insert data: %v", err)
		}
		fmt.Println("Data inserted successfully")

		pusherClient.Trigger("chat", "message", data)

		return c.JSON([]string{})
	})

	app.Get("/api/messages", func(c *fiber.Ctx) error {
		cursor, err := collection.Find(context.Background(), bson.D{})
		if err != nil {
			return err
		}

		var messages []bson.M
		if err := cursor.All(context.Background(), &messages); err != nil {
			return err
		}

		fmt.Println("Messages data loaded.")

		return c.JSON(messages)
	})

	// Users routes

	app.Post("/register", func(c *fiber.Ctx) error {
		var requestData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&requestData); err != nil {
			return err
		}
		if err := RegisterUser(requestData.Username, requestData.Password); err != nil {
			return err
		}
		return c.JSON(fiber.Map{"message": "User registered successfully"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var requestData struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&requestData); err != nil {
			fmt.Println("Login : Error on bodyparse")
			return err
		}
		authenticatedUser, err := AuthenticateUser(requestData.Username, requestData.Password)
		if err != nil {
			if err.Error() == "User "+requestData.Username+" is not found" {
				fmt.Println("Login : wrong user/password")
				// Return a custom JSON error response for "User %s is not found"
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Invalid username or password"})
			} else {
				fmt.Println("Login : Error on the else")
				return err
			}
		}
		fmt.Println("Login : success")
		return c.JSON(fiber.Map{"message": "Login successful", "username": authenticatedUser})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello world !")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
