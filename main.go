package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/develharsh/golang-fiber-postgresql-rest-api/models"
	"github.com/develharsh/golang-fiber-postgresql-rest-api/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

func (r *Repository) CreateBook(context *fiber.Ctx) error {
	book := Book{}
	err := context.BodyParser(&book)
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "Request Failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}
	context.Status(http.StatusCreated).JSON(&fiber.Map{"message": "Book has been added"})
	return nil
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"error": "could not get the books"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Books fetched successfully",
		"data": bookModels})
	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	ID, err := strconv.Atoi(context.Params("id"))
	if err != nil {
		fmt.Println(err)
	}
	bookModel := &models.Books{ID: uint(ID)}
	err = r.DB.First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"error": "could not get the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book fetched successfully",
		"data": bookModel})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	ID, err := strconv.Atoi(context.Params("id"))
	if err != nil {
		fmt.Println(err)
	}
	bookModel := &models.Books{ID: uint(ID)}
	err = r.DB.Delete(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"error": "could not delete the book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message": "Book deleted successfully"})
	return nil
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	fmt.Print(config)
	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal(err)
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	app.Use(logger.New())
	r.SetupRoutes(app)
	app.Listen("127.0.0.1:8080")
}
