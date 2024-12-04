package controllers

import (
	"fmt"
	"go-fiber-test/database"
	m "go-fiber-test/models"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile []m.Profile
	db.Find(&profile)
	return c.Status(200).JSON(profile)
}

func GetProfileById(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile
	id := c.Params("id")
	result := db.First(&profile, id)
	if result.Error != nil {
		return c.Status(404).SendString("Company not found")
	}
	return c.Status(200).JSON(profile)
}

func GetRangeProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile []m.Profile

	db.Find(&profile)
	sum_GenZ := 0
	sum_GenY := 0
	sum_GenX := 0
	sum_BaByBoomer := 0
	sum_GenGI := 0

	var dataResults []m.ProfileRes
	for _, v := range profile {
		GenStr := ""
		if v.Age < 24 {
			GenStr = "Gen Z"
			sum_GenZ += 1
		} else if v.Age >= 24 && v.Age <= 41 {
			GenStr = "Gen Y"
			sum_GenY += 1
		} else if v.Age >= 42 && v.Age <= 56 {
			GenStr = "Gen X"
			sum_GenX += 1
		} else if v.Age >= 57 && v.Age <= 75 {
			GenStr = "Baby Boomer"
			sum_BaByBoomer += 1
		} else {
			GenStr = "G.I. Generation"
			sum_GenGI += 1
		}

		d := m.ProfileRes{
			Name:       v.Name,
			EmployeeId: v.EmployeeId,
			Age:        v.Age,
			Gen:        GenStr,
		}
		dataResults = append(dataResults, d)
	}
	r := m.ResultGen{
		Data:          dataResults,
		Name:          "golang-test",
		Count:         len(profile),
		SumGenX:       sum_GenX,
		SumGenY:       sum_GenY,
		SumGenZ:       sum_GenZ,
		SumBaByBoomer: sum_BaByBoomer,
		SumGI:         sum_GenGI,
	}
	return c.Status(200).JSON(r)
}

func CreateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	profile := new(m.Profile)
	if err := c.BodyParser(&profile); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(profile); err != nil {
		var message []string
		for _, e := range err.(validator.ValidationErrors) {
			profileName := strings.ToLower(e.Field())

			switch profileName {
			case "employee_id":
				message = append(message, "Please enter your employee id")
			case "name":
				message = append(message, "Please enter your name")
			case "lastname":
				message = append(message, "Please enter your lastname")
			case "birthday":
				message = append(message, "Please enter your birthday")
			case "age":
				message = append(message, "Please enter your age")
			case "email":
				message = append(message, "Please enter your E-mail")
			case "tel":
				message = append(message, "Please enter your telephone")
			}
		}
		return c.Status(400).JSON(fiber.Map{
			"message": message,
		})
	}

	var existingProfile m.Profile
	if err := db.Where("employee_id = ?", profile.EmployeeId).First(&existingProfile).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Employee ID already exists",
			"error":   fmt.Sprintf("Employee ID %d already exists", profile.EmployeeId),
		})
	}

	userPattern := "^[a-zA-Z]+$"
	isValid := regexp.MustCompile(userPattern).MatchString
	if !isValid(profile.Name) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Username should only contain letters",
		})
	}

	if !isValid(profile.LastName) {
		return c.Status(400).JSON(fiber.Map{
			"message": "Lastname should only contain letters",
		})
	}

	// Pattern BirthDay
	formats := []string{"02/01/2006", "2006-01-02", "02-01-2006"}
	var parsedDate time.Time
	var err error
	for _, layout := range formats {
		parsedDate, err = time.Parse(layout, profile.BirthDay)
		if err == nil {
			break
		}
	}
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid birthday format. Supported formats: DD/MM/YYYY, YYYY-MM-DD, DD-MM-YYYY",
			"error":   err.Error(),
		})
	}
	profile.BirthDay = parsedDate.Format("2006-01-02") // Standardize output format

	telPattern := "^[0-9]{10}$"
	isValidTel, _ := regexp.MatchString(telPattern, profile.Tel)
	if !isValidTel {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please enter a valid phone number with 10 digits",
		})
	}

	if err := db.Create(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not save profile to database",
			"error":   err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Profile registered successfully",
		"profile": profile,
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile
	id := c.Params("id")
	result := db.First(&profile, id)
	if result.Error != nil {
		return c.Status(404).SendString("Profile not found")
	}

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&profile)
	db.Save(&profile)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    profile,
		"message": "Profile updated successfully",
	})
}

func DeleteProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile
	id := c.Params("id")
	result := db.First(&profile, id)
	if result.Error != nil {
		return c.Status(404).SendString("Profile not found")
	}
	db.Delete(&profile)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"data":    profile,
		"message": "Profile deleted successfully",
	})
}

func SearchProfile(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))

	if search == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุคำค้นหา",
		})
	}

	var profiles []m.Profile
	result := db.Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Where("employee_id LIKE ?", "%"+search+"%").
			Or("fname LIKE ?", "%"+search+"%").
			Or("last_name LIKE ?", "%"+search+"%")
	}).Find(&profiles)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "เกิดข้อผิดพลาดในการค้นหา",
			"error":   result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":      "success",
		"count":       len(profiles),
		"data":        profiles,
		"search_term": search,
	})
}
