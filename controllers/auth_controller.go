package controllers

import (
	"context"
	"log"
	"time"

	"github.com/RegiAdi/hatchet/helpers"
	"github.com/RegiAdi/hatchet/kernel"
	"github.com/RegiAdi/hatchet/models"
	"github.com/RegiAdi/hatchet/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := kernel.Mongo.DB.Collection("users")

	var request models.User
	var user models.User
	var userLoginResponse responses.UserLoginResponse

	if err := c.BodyParser(&request); err != nil {
		log.Println(err)

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	err := userCollection.FindOne(ctx, bson.D{{Key: "username", Value: request.Username}}).Decode(&user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
			"error":   err,
		})
	}

	if !helpers.CheckPasswordHash(request.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Password do not match",
			"error":   nil,
		})
	}

	APIToken, _ := helpers.GenerateAPIToken()
	APITokenExpirationDate := helpers.GenerateAPITokenExpiration()

	userID, _ := primitive.ObjectIDFromHex(user.ID)
	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "api_token", Value: APIToken},
			{Key: "device_name", Value: request.DeviceName},
			{Key: "token_expires_at", Value: APITokenExpirationDate},
			{Key: "updated_at", Value: helpers.GetCurrentTime()},
		},
		}}

	err = userCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&userLoginResponse)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to generate API Token",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User authenticated successfully",
		"data":    userLoginResponse,
	})
}

func Logout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := kernel.Mongo.DB.Collection("users")

	reqHeader := c.GetReqHeaders()

	filter := bson.D{{Key: "api_token", Value: reqHeader["Authorization"]}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "api_token", Value: nil},
			{Key: "token_expires_at", Value: time.Time{}},
			{Key: "updated_at", Value: helpers.GetCurrentTime()},
		},
		}}

	result, err := userCollection.UpdateOne(ctx, filter, update)

	if err != nil || result.ModifiedCount < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete API Token",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "User logged out successfully",
	})
}

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userCollection := kernel.Mongo.DB.Collection("users")

	user := new(models.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	usernameCount, _ := userCollection.CountDocuments(ctx, bson.D{{Key: "username", Value: user.Username}})
	if usernameCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Username already exist",
		})
	}

	emailCount, _ := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})
	if emailCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"success": false,
			"message": "Email already exist",
		})
	}

	password, err := helpers.HashPassword(user.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User registration failed",
			"error":   err,
		})
	}

	user.Password = password
	user.CreatedAt = helpers.GetCurrentTime()
	user.UpdatedAt = helpers.GetCurrentTime()

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "User registration failed",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "User created successfully",
	})
}
