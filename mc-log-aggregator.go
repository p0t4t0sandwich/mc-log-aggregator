package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// // urlExample := "postgres://username:password@localhost:5432/database_name"
// conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
// if err != nil {
// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
// 	os.Exit(1)
// }
// defer conn.Close(context.Background())

// var name string
// var weight int64
// err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
// if err != nil {
// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
// 	os.Exit(1)
// }

// fmt.Println(name, weight)

func main() {
	// Get IP from env
	ip := os.Getenv("IP_ADDRESS")
	if ip == "" {
		ip = "0.0.0.0"
	}

	// Get port from env
	port := os.Getenv("REST_PORT")
	if port == "" {
		port = "8080"
	}

	var router *gin.Engine = gin.Default()

	// Minecraft Server Status
	router.GET("/", getRoot)

	// Test
	router.GET("/test", getTest)
	router.POST("/upload", uploadTest)

	router.Run(ip + ":" + port)
}

// -------------- Structs --------------

// -------------- Enums --------------

// -------------- Functions --------------

// -------------- Handlers --------------

// Get root route
func getRoot(c *gin.Context) {
	// Read the html file
	html, err := os.ReadFile("templates/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Return the html
	c.Data(http.StatusOK, "text/html", html)
}

// Returns test HTMX html inject
func getTest(c *gin.Context) {
	// Return the html
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, "<div><h1>Whoop, Worked!</h1></div>")
}

// Upload test
func uploadTest(c *gin.Context) {
	// Get the file
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Save the file
	err = c.SaveUploadedFile(file, "uploads/"+file.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
