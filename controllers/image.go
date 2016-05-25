package controllers

import (
    // Standard library packages
    "fmt"
    "net/http"
    "encoding/json"
    "log"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    // Third party packages
    "github.com/julienschmidt/httprouter"
    // Use defined model in another dir
    "../models"
    "../utils"
)

// ImageController represents the controller for operating on the image resource
// i.e) controller has no property, only methods
type ImageController struct{
    // use reference to access mongodb
    session *mgo.Session
}

// Constructor
func NewImageController(s *mgo.Session) *ImageController {
    // instantiate with the session received as an arg
    return &ImageController{s}
}

func (ic ImageController) CreateImage(rw http.ResponseWriter, req *http.Request, param httprouter.Params) {
    imgInfo, err := createImage(ic, req, param)

    // Create response
    // Write content-type, statuscode, payload
    if err != nil {
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
    } else {
        // Marshal provided interface into JSON structure
        imgJson, _ := json.Marshal(imgInfo)
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(201)
        fmt.Fprintf(rw, "%s\n", imgJson)
    }
}

func createImage(ic ImageController, req *http.Request, param httprouter.Params) (models.ImageInfo, error) {
    userID := param.ByName("userId")

    file, header, err := req.FormFile("image")
	if err != nil {
		log.Println(err)
		return models.ImageInfo{}, err
	}

    fileName, imageURL, err := utils.UploadImage(file, header, userID)
    if err != nil {
        log.Println(err)
        return models.ImageInfo{}, err
    }

    histogram, err := utils.ConvertToGrayScale(file)
    if err != nil {
        log.Println(err)
        return models.ImageInfo{}, err
    }

    // Stub an image to be populated from the body
    imgInfo := models.ImageInfo{}
    // assign unique string ID
    imgInfo.Id = bson.NewObjectId()
    // Append all data
    imgInfo.UserID = userID
    imgInfo.FileName = fileName
    imgInfo.ImageURL = imageURL
    imgInfo.Timestamp = time.Now()
    imgInfo.Histogram = histogram

    // Persist the data to mongodb
    conn := ic.session.DB("image").C("image_info")
    err = conn.Insert(imgInfo)
    if err != nil {
        log.Println(err)
        return models.ImageInfo{}, err
    }

    return imgInfo, nil
}
