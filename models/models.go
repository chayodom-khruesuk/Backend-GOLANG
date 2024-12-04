package models

import "gorm.io/gorm"

type Person struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

type BusinessType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var BusinessTypes = []BusinessType{
	{ID: 1, Name: "Technology"},
	{ID: 2, Name: "Manufacturing"},
	{ID: 3, Name: "Retail"},
	{ID: 4, Name: "Healthcare"},
	{ID: 5, Name: "Finance"},
	{ID: 6, Name: "E-commerce"},
	{ID: 7, Name: "Education"},
	{ID: 8, Name: "Consulting"},
	{ID: 9, Name: "Marketing"},
	{ID: 10, Name: "Food & Beverage"},
}

type User struct {
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email,omitempty" validate:"required,email"`
	Password     string `json:"password,omitempty" validate:"required,min=6,max=20"`
	LineId       string `json:"lineid,omitempty" validate:"required"`
	Phone        string `json:"phone,omitempty" validate:"required"`
	WebsiteName  string `json:"website_name,omitempty" validate:"required"`
	TypeBusiness string `json:"type_business,omitempty" validate:"required"`
}

type Dogs struct {
	gorm.Model
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
}

type ResultData struct {
	Data       []DogsRes `json:"data"`
	Name       string    `json:"name"`
	Count      int       `json:"count"`
	SumRed     int       `json:"sumRed"`
	SumGreen   int       `json:"sumGreen"`
	SumPink    int       `json:"sumPink"`
	SumNoColor int       `json:"sumNoColor"`
}

type DogsRes struct {
	Name  string `json:"name"`
	DogID int    `json:"dog_id"`
	Type  string `json:"type"`
}

type Company struct {
	CompanyName     string
	Address         string
	TaxID           string
	Type            string
	Phone           string
	Email           string
	Website         string
	EstablishedYear int
	EmployeeCount   int
	Revenue         float64
	Industry        string
}

type ProfileRes struct {
	Name       string `json:"name"`
	EmployeeId int    `json:"employee_id"`
	Age        int    `json:"age"`
	Gen        string `json:"gen"`
}

type ResultGen struct {
	Data          []ProfileRes `json:"data"`
	Name          string       `json:"name"`
	Count         int          `json:"count"`
	SumGenZ       int          `json:"gen_z"`
	SumGenY       int          `json:"gen_y"`
	SumGenX       int          `json:"gen_x"`
	SumBaByBoomer int          `json:"baby_boomer"`
	SumGI         int          `json:"gen_g"`
}

type Profile struct {
	gorm.Model
	EmployeeId int    `json:"employee_id,omitempty" validate:"required" gorm:"unique"`
	Name       string `json:"name,omitempty" validate:"required,min=2,max=30"`
	LastName   string `json:"lastname,omitempty" validate:"required,min=2,max=30"`
	BirthDay   string `json:"birthday,omitempty" validate:"required" gorm:"type:date"`
	Age        int    `json:"age,omitempty" validate:"required"`
	Email      string `json:"email,omitempty" validate:"required"  gorm:"unique"`
	Tel        string `json:"tel,omitempty" validate:"required"  gorm:"unique"`
}
