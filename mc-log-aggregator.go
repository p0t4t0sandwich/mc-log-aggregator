package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// -------------- Global Variables --------------

var postgresConn *pgx.Conn

func main() {
	var err error
	postgresConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer postgresConn.Close(context.Background())

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

	// router.Run(ip + ":" + port)

	stringList := parseFile("./uploads/2022-04-20-2.log")

	uploadToDatabase("Hub", "Bukkit", "test0101", stringList)
}

// -------------- Structs --------------

type LogEntry struct {
	Server     string
	ServerType string
	Timestamp  string
	Level      string
	Source     string
	Message    string
}

// -------------- Enums --------------

// -------------- Functions --------------

// Function to parse a base64 encoded string blob and split it into an array of strings
func parseBlob(blob string) []string {
	text, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		fmt.Println("decode error:", err)
		return nil
	}

	stringList := strings.Split(string(text), "\n")

	// for _, str := range stringList {
	// 	fmt.Println(str)
	// }

	return stringList
}

// Function to parse a file and split it into an array of strings
func parseFile(file string) []string {
	// Read the file
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Convert the file to a string
	stringData := string(data)

	// Split the string into an array of strings
	stringList := strings.Split(stringData, "\n")

	// for _, str := range stringList {
	// 	fmt.Println(str)
	// }

	return stringList
}

// Function to parse an array of strings and upload them to the database
func uploadToDatabase(server string, serverType string, logDate string, stringList []string) {
	// Timestamp in the form of "[02:26:00]"
	timeStampRegex := regexp.MustCompile(`\[\d{2}:\d{2}:\d{2}\]`)
	// Level in the form of "[text]:"
	levelRegex := regexp.MustCompile(`\[([^\]]+)\]:`)
	// Source in the form of ": [text]"
	sourceRegex := regexp.MustCompile(`: \[([^\]]+)\]`)
	// Message in the form of "] The text is yes"
	// messageRegex := regexp.MustCompile(`\] .+`)

	var logEntries []LogEntry

	for _, str := range stringList {
		// Check for a timestamp
		if timeStampRegex.MatchString(str) {
			// Get the timestamp
			timestamp := timeStampRegex.FindString(str)
			// Get the level
			level := levelRegex.FindString(str)
			// Get the message
			message := str
			// Get the source
			source := sourceRegex.FindString(message)
			// Remove the : from the source
			source = strings.Trim(source, ":")

			// Remove everything from the message except the message
			message = strings.Replace(message, level+" ", "", 1)
			message = strings.Replace(message, source, "", 1)
			message = strings.Replace(message, timestamp+" ", "", 1)

			// Remove the brackets from the timestamp
			timestamp = strings.Trim(timestamp, "[]")
			// Remove the brackets from the level
			level = strings.Trim(level, "[]:")
			// Remove the brackets from the source
			source = strings.Trim(source, " []")

			// Create a log entry
			logEntry := LogEntry{
				Server:     server,
				ServerType: serverType,
				Timestamp:  logDate + " " + timestamp,
				Level:      level,
				Source:     source,
				Message:    message,
			}

			// Append the log entry to the log entries
			logEntries = append(logEntries, logEntry)
		} else {
			// Assume that the message is a continuation of the previous message
			// Append the message to the previous message
			logEntries[len(logEntries)-1].Message += str
		}

		for _, logEntry := range logEntries {
			fmt.Println("--------------------")
			fmt.Println(logEntry)
		}
	}
}

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
