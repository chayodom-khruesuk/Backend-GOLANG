package controllers

import (
	"strings"

	"go-fiber-test/database"
	m "go-fiber-test/models"

	"github.com/gofiber/fiber/v2"
)

func GetDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs)
	return c.Status(200).JSON(dogs)
}

func GetDogSearch(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var dog []m.Dogs

	result := db.Find(&dog, "dog_id = ?", search)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&dog)
}

func GetDogDelete(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs
	result := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&dogs)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No deleted records found",
		})
	}
	return c.Status(200).JSON(dogs)
}

func GetDogRange(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	result := db.Where("dog_id > ? AND dog_id < ?", 50, 100).Find(&dogs)

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "No dogs found in this range",
		})
	}
	return c.Status(200).JSON(dogs)
}

func AddDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dogs

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dogs
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&dog)
	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var dog m.Dogs

	result := db.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs)
	sum_red := 0
	sum_green := 0
	sum_pink := 0
	sum_nocolor := 0

	var dataResults []m.DogsRes
	for _, v := range dogs {
		typeStr := ""
		if v.DogID >= 10 && v.DogID <= 50 {
			typeStr = "red"
			sum_red += 1
		} else if v.DogID >= 100 && v.DogID <= 150 {
			typeStr = "green"
			sum_green += 1
		} else if v.DogID >= 200 && v.DogID <= 250 {
			typeStr = "pink"
			sum_pink += 1
		} else {
			typeStr = "no color"
			sum_nocolor += 1
		}

		d := m.DogsRes{
			Name:  v.Name,
			DogID: v.DogID,
			Type:  typeStr,
		}
		dataResults = append(dataResults, d)
	}

	r := m.ResultData{
		Data:       dataResults,
		Name:       "golang-test",
		Count:      len(dogs),
		SumRed:     sum_red,
		SumGreen:   sum_green,
		SumPink:    sum_pink,
		SumNoColor: sum_nocolor,
	}
	return c.Status(200).JSON(r)
}
