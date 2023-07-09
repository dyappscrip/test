package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type InvitationData struct {
	GuestName     string
	EventDate     string
	EventTime     string
	EventLocation string
	RSVPlink      string
	ContactName   string
	ContactEmail  string
	ContactPhone  string
}
type User struct {
	Name  string `json:"name" bson:"name"`
	Index int    `json:"index" bson:"index"`
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb+srv://root:dhanraj123@cluster0.c67dfhe.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("welcome")
	collection := db.Collection("users")

	router.GET("/:index", func(c *gin.Context) {
		indexStr := c.Param("index")
		num, _ := strconv.Atoi(indexStr)
		filter := bson.M{"index": int(num)}
		var user User
		collection.FindOne(context.Background(), filter).Decode(&user)

		data := InvitationData{
			GuestName:     user.Name,
			EventDate:     "July 15, 2023",
			EventTime:     "7:00 PM",
			EventLocation: "123 Main Street, City",
			RSVPlink:      "http://example.com/rsvp",
			ContactName:   "Jane Smith",
			ContactEmail:  "jane@example.com",
			ContactPhone:  "123-456-7890",
		}

		c.HTML(http.StatusOK, "invitation.html", gin.H{
			"data": data,
		})
	})

	router.GET("/add", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add.html", gin.H{})
	})

	router.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"index": user.Index}
		update := bson.M{"$set": bson.M{"name": user.Name}}

		resp, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
		fmt.Printf("resp: %v\n", resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User name added successfully"})
	})

	router.Run(":8080")
}
