package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
    "fmt"
    "log"
    "gopkg.in/mgo.v2"
    "./controllers"
)

const (
     dbUser string = "yaopeng"
     dbPassword string = "admin"
     dbServer string = "ds051553.mlab.com"
     dbPort string = "51553"
     dbName string = "image"
)

func main() {
    // Instantiate a new router
    router := httprouter.New()

    // Get a UserController instance
    imageController := controllers.NewImageController(getMongoSession())

    // Endpoints
    router.POST("/images/:userId", imageController.CreateImage)
    router.GET("/images/weekly/:userId", imageController.GetWeeklyHistograms)
    router.GET("/images/day/median", imageController.GetMedianHistogram)
    router.GET("/images/similarity/:userId/:n", imageController.GetUserIDWithSimilarHistogram)

    // Fire up the server
    fmt.Println("Server listening on 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getMongoSession() *mgo.Session {
    // Test
    // session, err := mgo.Dial("mongodb://localhost")

    // Production
    url := "mongodb://" + dbUser + ":" + dbPassword + "@" + dbServer + ":" + dbPort + "/" + dbName
    session, err := mgo.Dial(url)

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }

    return session
}
