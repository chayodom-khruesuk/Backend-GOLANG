package controllers

import (
	"go-fiber-test/database"
	"go-fiber-test/models"

	"github.com/gofiber/fiber/v2"
)

func GetCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var companies []models.Company
	db.Find(&companies)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    companies,
		"message": "Companies retrieved successfully",
	})
}

func GetCompanyById(c *fiber.Ctx) error {
	db := database.DBConn
	var company models.Company
	id := c.Params("id")
	result := db.First(&company, id)
	if result.Error != nil {
		return c.Status(404).SendString("Company not found")
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    company,
		"message": "Company retrieved successfully",
	})
}

func CreateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db.Create(&company)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    company,
		"message": "Company created successfully",
	})
}

func UpdateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company models.Company
	id := c.Params("id")
	result := db.First(&company, id)
	if result.Error != nil {
		return c.Status(404).SendString("Company not found")
	}
	if err := c.BodyParser(&company); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	db.Save(&company)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    company,
		"message": "Company updated successfully",
	})
}

func DeleteCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company models.Company
	id := c.Params("id")
	result := db.First(&company, id)
	if result.Error != nil {
		return c.Status(404).SendString("Company not found")
	}
	db.Delete(&company)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    company,
		"message": "Company deleted successfully",
	})
}
