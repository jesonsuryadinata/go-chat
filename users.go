package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, password string) error {
	filter := bson.M{"username": username}
	count, err := usersCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Println("error while filtering db on register")
		return err
	}

	if count > 0 {
		// User already exists
		fmt.Println("User already exists")
		return fmt.Errorf("Username has already been used", username)
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("Error while hashing password on register")
		log.Fatal(err)
	}

	// If user doesn't exist, INSERT
	user := bson.M{
		"username": username,
		"password": hashedPassword,
	}

	_, err = usersCollection.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("Error while inserting user")
		return err
	}

	userId, err := GetUserIdByUsername(username)
	if err != nil {
		fmt.Errorf("Error on GetUserIdByUsername (users.RegisterUser)")
		return err
	}

	addToGlobal, err := AddUserIntoChannel(int64(userId), 1)
	if err != nil {
		fmt.Errorf("Error on AddUserIntoChannel(users.RegisterUser)")
	}
	fmt.Println(addToGlobal)

	addUserChannel, err := AddUserChannelRelation(int64(userId), 1, "member")
	if err != nil {
		fmt.Errorf("Error on AddUserChannelRelation(users.RegisterUser)")
	}
	fmt.Println(addUserChannel)

	fmt.Printf("Username [ %s ] registered.\n", username)
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func AuthenticateUser(username, password string) (string, error) {
	// Find by username
	hashedPassword, err := HashPassword(password)
	if err != nil {
		fmt.Println("Authenticate user : Error while hashing password input")
		return "", err
	}
	var user bson.M
	filter := bson.M{"username": username}
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("Authenticate user : Error user not found. Input: \nUsername : %s | Hashed password : %s | Filter : %s\n", username, hashedPassword, filter)
			return "", fmt.Errorf("User %s is not found", username)
		}
		fmt.Println("Authenticate user : Error while retrieving user from db")
		return "", err
	}

	storedPassword, ok := user["password"].(string)
	// // Compare hashed password with db-stored password
	// if !ok || storedPassword != hashedPassword {
	// 	fmt.Println("Authenticate user : invalid username or password")
	// 	return "", fmt.Errorf("Invalid username or password")
	// }
	if !ok {
		fmt.Println("Password field in the db is not a string")
		return "", fmt.Errorf("Invalid password format in the db")
	}

	// Compare the stored hashed password with user password input
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		// Passwords do not match or another error occured
		fmt.Println("Password comparison failed")
		return "", err
	}

	fmt.Println("Authenticate user success")
	//Authentication successful
	return username, nil
}

func GetUserIdByUsername(username string) (int64, error) {
	var user bson.M
	filter := bson.M{"username": username}
	err := usersCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("GetUserIdByUsername : Error user not found. Input: [Username : %s]\n", username)
			return 0, fmt.Errorf("User %s is not found", username)
		}
		fmt.Println("GetUserIdByUsername : Error while retrieving user from db")
		return 0, err
	}
	userId, ok := user["userId"].(int64)

	fmt.Printf("GetUserIdByUsername : %d | ok : %t\n", userId, ok)
	return userId, nil
}
