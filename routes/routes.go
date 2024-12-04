package routes

import (
	c "go-fiber-test/controllers"
	m "go-fiber-test/middleware"

	"github.com/gofiber/fiber/v2"
)

func InetRoutes(app *fiber.App) {
	auth := m.AuthMiddleware()
	authDemo := m.AuthDemoMiddleware()

	// Group
	ex := app.Group("api/ex", authDemo)
	v3 := app.Group("api/v3", authDemo)
	dog := app.Group("/dog", authDemo)
	company := app.Group("/company", authDemo)
	user := app.Group("api/v1/user")

	user.Get("", c.GetProfile)
	user.Get("/id", auth, c.GetProfileById)
	user.Get("/range", auth, c.GetRangeProfile)
	user.Post("/register", auth, c.CreateProfile)
	user.Put("/update/:id", auth, c.UpdateProfile)
	user.Delete("/delete/:id", auth, c.DeleteProfile)

	ex.Get("/fact/:value", c.Factorial)
	ex.Post("/registerDemo", c.Register)

	v3.Get("/film", c.AsciiConv)

	dog.Get("", c.GetDogs)
	dog.Get("/filter", c.GetDog)
	dog.Get("/json", c.GetDogsJson)
	dog.Get("/del", c.GetDogDelete)
	dog.Get("/range", c.GetDogRange)
	dog.Post("/", c.AddDog)
	dog.Put("/:id", c.UpdateDog)
	dog.Delete("/:id", c.RemoveDog)

	company.Get("", c.GetCompany)
	company.Post("/create", c.CreateCompany)
	company.Put("/update/:id", c.UpdateCompany)
	company.Delete("/delete/:id", c.DeleteCompany)

	app.Get("", c.HelloTest)
	app.Get("/user/:name", c.WordConnect)
	app.Post("/", c.BodyParser)
	app.Post("/inet", c.QueryTest)
	app.Post("/valid", c.ValidatorTest)
}
