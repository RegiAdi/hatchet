package controllers

import (
	"context"
	"time"

	"github.com/RegiAdi/pos-mobile-backend/bootstrap"
	"github.com/RegiAdi/pos-mobile-backend/helpers"
	"github.com/RegiAdi/pos-mobile-backend/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProduct(c *fiber.Ctx) error {
	productCollection := bootstrap.MongoDB.Database.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product models.Product

	objID, _ := primitive.ObjectIDFromHex(c.Params("id"))
	err := productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Product not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Product retrieved successfully",
		"data":    product,
	})
}

func CreateProduct(c *fiber.Ctx) error {
	productCollection := bootstrap.MongoDB.Database.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	product := new(models.Product)

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Failed to parse body",
			"error":   err,
		})
	}

	product.CreatedAt = primitive.NewDateTimeFromTime(helpers.GetCurrentTime())
	product.UpdatedAt = primitive.NewDateTimeFromTime(helpers.GetCurrentTime())

	result, err := productCollection.InsertOne(ctx, product)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Create new product failed",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    result,
		"success": true,
		"message": "Product created successfully",
	})
}
