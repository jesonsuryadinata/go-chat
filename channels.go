package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func AddUserIntoChannel(userId int, channelId int) (string, error) {
	filter := bson.M{"channelId": channelId}
	update := bson.M{"$push": bson.M{"members": userId}}

	result, err := channelsCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("error while updating")
		return "", err
	}
	if result.ModifiedCount == 0 {
		fmt.Println("error : channel not found")
		return "", err
	}

	return "Added user into Global Channel", nil
}
