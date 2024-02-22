package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type Person struct {
	Name  string `json:"name" db:"name"`
	Age   int    `json:"age" db:"age"`
	Email string `json:"email" db:"email"`
}

var (
	stringConn = `postgres://joaof:123456@localhost/joaof?sslmode=disable`
	dockerConn = `postgres://joaof:123456@host.docker.internal/joaof?sslmode=disable`
)

type Conn struct {
	DB *sqlx.DB
}

func Connect() *Conn {
	env := os.Getenv("ENVI")
	fmt.Println("env:")
	fmt.Println(env)

	if env == "docker" {
		stringConn = dockerConn
	}

	fmt.Println(stringConn)

	db, err := sqlx.Open("postgres", stringConn)
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	conn := new(Conn)
	conn.DB = db

	return conn
}

func (c *Conn) InsertPerson(p Person) {
	q := `
		INSERT INTO public.person (name, age, email)
		VALUES (:name, :age, :email);
	`

	_, err := c.DB.NamedExec(q, p)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(p)
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/person", PostPerson)
	e.Logger.Fatal(e.Start(":3000"))
}

func PostPerson(e echo.Context) error {
	p := new(Person)

	if err := e.Bind(p); err != nil {
		return e.String(http.StatusBadRequest, "bad request")
	}

	conn := Connect()
	defer conn.DB.Close()

	conn.InsertPerson(*p)

	return e.JSON(http.StatusOK, p)
}
