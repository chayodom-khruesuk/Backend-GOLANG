package controllers

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	m "go-fiber-test/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func HelloTest(c *fiber.Ctx) error {
	return c.SendString("Hello")
}

func WordConnect(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

func BodyParser(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	log.Println(p.Name)
	log.Println(p.Pass)
	str := p.Name + p.Pass
	return c.JSON(str)
}

func QueryTest(c *fiber.Ctx) error {
	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

func ValidatorTest(c *fiber.Ctx) error {
	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(user)
}

func Factorial(c *fiber.Ctx) error {
	value := c.Params("value")
	num, err := strconv.Atoi(value)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	factorial := 1
	for i := 1; i <= num; i++ {
		factorial *= i
	}
	return c.Status(200).SendString(strconv.Itoa(num) + "! = " + strconv.Itoa(factorial))

}

func Register(c *fiber.Ctx) error {
	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	validate := validator.New()
	if err := validate.Struct(user); err != nil {
		var messages []string
		for _, e := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(e.Field())
			switch fieldName {
			case "name":
				messages = append(messages, "Please enter your name")
			case "email":
				messages = append(messages, "Please enter your email")
			case "password":
				messages = append(messages, "Please enter your password")
			case "lineid":
				messages = append(messages, "Please enter your Line ID")
			case "phone":
				messages = append(messages, "Please enter a valid phone number with 10 digits")
			case "websitename":
				messages = append(messages, "Please enter your website name")
			case "typebusiness":
				messages = append(messages, "Please select your business type")
			}
		}
		return c.Status(400).JSON(fiber.Map{
			"messages": messages,
		})
	}

	phonePattern := "^[0-9]{10}$"
	matched, _ := regexp.MatchString(phonePattern, user.Phone)
	if !matched {
		return c.Status(400).JSON(fiber.Map{
			"message": "Please enter a valid phone number with 10 digits",
		})
	}

	websitePattern := "^(https?://)?([a-z0-9.-]+/.[a-z]{2,})(/[a-z0-9._-]+)*(/)?$"
	matched, _ = regexp.MatchString(websitePattern, user.WebsiteName)
	if !matched {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid pattern website",
		})
	}

	businessPattern := "^(" + strings.Join(getValidTypeNames(), "|") + ")$"
	matched, _ = regexp.MatchString(businessPattern, user.TypeBusiness)
	if !matched {
		return c.Status(400).JSON(fiber.Map{
			"message":         fmt.Sprintf("Invalid business type: %s", user.TypeBusiness),
			"available_types": m.BusinessTypes,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func getValidTypeNames() []string {
	var validTypes []string
	for _, t := range m.BusinessTypes {
		validTypes = append(validTypes, t.Name)
	}
	return validTypes
}

func AsciiConv(c *fiber.Ctx) error {
	tax_id := c.Query("tax_id")

	if tax_id == "" {
		return c.Status(400).SendString("tax_id is required")
	}

	var asciiValues []string
	for _, char := range tax_id {
		asciiValues = append(asciiValues, strconv.Itoa(int(char)))
	}

	return c.Status(200).SendString(strings.Join(asciiValues, " "))
}
