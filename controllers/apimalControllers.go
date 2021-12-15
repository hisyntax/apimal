package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hisyntax/apimal/database"
	"github.com/hisyntax/apimal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var animalCollection *mongo.Collection = database.OpenCollection(database.Client, "animal")
var validate = validator.New()

func CreateAnimalHandler(c *gin.Context) {
	//open up a dataase connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//close that connection on the insertion is done
	defer cancel()

	//declear a variable which represents the animal model struct
	var animal models.Animal

	//bind the animal model together and check if an error occured while binding
	if err := c.BindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//user the validator package to validate the user inputes to make sure some
	//some required fields are not left blank and throw an error
	validateErr := validate.Struct(&animal)
	if validateErr != nil {
		// log.Panic(validateErr)
		msg := "Some fields are left blank"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": msg,
		})
		return
	}

	//count throw al the database in the database to to check for an exception
	//if the name of the animal privided is already availavle in the database or not
	count, err := animalCollection.CountDocuments(ctx, bson.M{"name": animal.Name})
	if err != nil {
		log.Panic(err)
		msg := "Sorry, an error occured"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": msg,
		})
		return
	}

	//if the privided animal name exists, throw an error
	//notifying the user that the animal name exists already
	if count > 0 {
		msg := "The name already exists"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": msg,
		})
		return
	}

	//add a created at time to the animal data to be created
	animal.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// animal.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	//Init a mongodb default id for the Id fields in the models
	animal.ID = primitive.NewObjectID()
	//assign that ID to the animal_id so that the are the same
	animal.Animal_id = animal.ID.Hex()

	//tru to save the prvided animal data and if an error occured
	//return an error message
	insert, err := animalCollection.InsertOne(ctx, animal)
	if err != nil {
		msg := "An error occured while trying the save the animal information"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": msg,
		})
		return
	}

	//if there has been no errors this far
	//then save the animal data
	c.JSON(http.StatusOK, gin.H{
		"name":      animal.Name,
		"desc":      animal.Description,
		"image":     animal.Image,
		"animal_id": insert,
	})
}

func GetAnimalsHandler(c *gin.Context) {
	//opens a database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//close that connection after user
	defer cancel()

	//declear an aniamls variable ro hold an array of the animal model
	var animals []models.Animal

	//declear a variable to find all the datas in the database
	//check if there is an error and handle that error appropriately
	cusor, err := animalCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
		msg := "Sorry, something went wrong. Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": msg,
		})
		return
	}

	//user the initially create varivale to finc all the aninmal data and handle the
	// error if any
	if err = cusor.All(ctx, &animals); err != nil {
		log.Panic(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	//close the connection
	defer cusor.Close(ctx)

	//check if there is an error in retrieving the animal datas
	if err := cusor.Err(); err != nil {
		log.Panic(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//list out the animal datas found
	c.JSON(http.StatusOK, gin.H{
		"animal": animals,
	})
}
