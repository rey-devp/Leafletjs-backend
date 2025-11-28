package handler

import (
	"context"
	"net/http"

	"git-uts/config"
	"git-uts/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func GetLocations(c *fiber.Ctx) error {
	var locations []model.Location

	collection := config.DB.Collection("location")

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data lokasi",
		})
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var location model.Location
		if err := cursor.Decode(&location); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Gagal decode data",
			})
		}
		locations = append(locations, location)
	}

	return c.JSON(locations)
}


func CreateLocation(c *fiber.Ctx) error {
	var input model.Location

	// Parse JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format JSON tidak valid",
		})
	}

	// Validasi minimal sesuai model
	if input.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama lokasi wajib diisi",
		})
	}

	// Pastikan ID kosong → MongoDB yang generate
	input.ID = primitive.NilObjectID

	collection := config.DB.Collection("location")

	result, err := collection.InsertOne(context.TODO(), input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data lokasi",
		})
	}

	input.ID = result.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(input)
}


func UpdateLocation(c *fiber.Ctx) error {
	idParam := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	// REQUEST STRUCT (biar 0 boleh di-update)
	type UpdateLocationRequest struct {
		Name        *string  `json:"name"`
		Latitude    *float64 `json:"latitude"`
		Longitude   *float64 `json:"longitude"`
		Description *string  `json:"description"`
	}

	var input UpdateLocationRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Format JSON tidak valid",
		})
	}

	updateData := bson.M{}

	if input.Name != nil {
		updateData["name"] = *input.Name
	}
	if input.Latitude != nil {
		updateData["latitude"] = *input.Latitude
	}
	if input.Longitude != nil {
		updateData["longitude"] = *input.Longitude
	}
	if input.Description != nil {
		updateData["description"] = *input.Description
	}

	if len(updateData) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Tidak ada data yang diperbarui",
		})
	}

	collection := config.DB.Collection("location")

	result, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memperbarui data lokasi",
		})
	}

	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Lokasi tidak ditemukan",
		})
	}

	// Ambil data terbaru
	var updated model.Location
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&updated)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data terbaru",
		})
	}

	return c.JSON(updated)
}


func DeleteLocation(c *fiber.Ctx) error {
	idParam := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "ID tidak valid",
		})
	}

	collection := config.DB.Collection("location")

	result, err := collection.DeleteOne(
		context.TODO(),
		bson.M{"_id": objectID},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menghapus data lokasi",
		})
	}

	if result.DeletedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Lokasi tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Lokasi berhasil dihapus",
	})
}
