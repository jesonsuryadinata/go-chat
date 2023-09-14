package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func AddUserIntoChannel(userId int64, channelId int64) (string, error) {
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

func AddUserChannelRelation(userId int64, channelId int64, status string) (string, error) {
	filter := bson.M{"userId": userId}
	update := bson.M{"$addToSet": bson.M{"channelId": channelId}}

	result, err := userChannelCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return "", err
	}
	if result.ModifiedCount == 0 {
		// New user registered, create new document
		newUser := bson.M{
			"userId":    userId,
			"channelId": []int64{channelId},
			"status":    status,
		}

		_, err := userChannelCollection.InsertOne(context.TODO(), newUser)
		if err != nil {
			fmt.Println("Error on InsertOne : ", err)
			return "", err
		}
		// fmt.Println("Success add to user-channel (new)")
		return "Success add to user-channel (new)", nil
	}
	fmt.Println("Success add to user-channel (existing)")
	return "Success add to user-channel (existing)", nil
}
