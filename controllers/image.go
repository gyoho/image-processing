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
    "errors"
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

func (ic ImageController) GetHistogram(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    usr, err := retriveAllHistogramsOfUser(ic, param.ByName("id"))
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        usrJson, _ := json.Marshal(usr)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", usrJson)
    }
}

// func (uc UserController) GetMedianHistogram(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
//     usr, err := retriveUserById(uc, param.ByName("id"))
//     if err != nil {
// 		log.Println(err)
//         // Write content-type, statuscode, payload
//         rw.Header().Set("Content-Type", "plain/text")
//         rw.WriteHeader(400)
//         fmt.Fprintf(rw, "%s\n", err)
// 	} else {
//         // Create response
//         // Marshal provided interface into JSON structure
//         usrJson, _ := json.Marshal(usr)
//         // Write content-type, statuscode, payload
//         rw.Header().Set("Content-Type", "application/json")
//         rw.WriteHeader(200)
//         fmt.Fprintf(rw, "%s\n", usrJson)
//     }
// }
//

func retriveAllHistogramsOfUser(ic ImageController, user_id string) ([]bson.M, error) {
    // make connection
    conn := ic.session.DB("image").C("image_info")
    // prepare query
    year, week := time.Now().ISOWeek()
    pipeline := []bson.M {
        bson.M {
            "$project": bson.M {
    			"_id": 0,
    	        "user_id": 1,
    	        "year": bson.M{ "$year": "$timestamp"},
    	        "week": bson.M{ "$week": "$timestamp"},
    			"histogram": 1}},
        bson.M {
            "$match": bson.M {
    			"user_id": user_id,
    	        "year": year,
    	        "week": week}},
        bson.M {
            "$project": bson.M {"histogram": 1}}}

    pipe := conn.Pipe(pipeline)
    res := []bson.M{}
    err := pipe.All(&res)
    if err != nil {
        return []bson.M{}, errors.New("No user found with this ID")
    }
    fmt.Println(res)

    return res, nil
}

// func retriveAllHistogramsOfUser(ic ImageController, user_id string) ([]models.Histogram, error) {
//     // make connection
//     conn := ic.session.DB("image").C("image_info")
//     // prepare query
//     now := time.Now()
//     year := now.Year()
//     day := now.YearDay()
//     pipeline := []bson.M {
//         bson.M {
//             "$project": bson.M {
//     			"_id": 0,
//     	        "user_id": 1,
//     	        "year": bson.M{ "$year": "$timestamp"},
//     	        "day": bson.M{ "$dayOfYear": "$timestamp"},
//     			"histogram": 1}},
//         bson.M {
//             "$match": bson.M {
//     			"user_id": user_id,
//     	        "year": year,
//     	        "day": day}}}
//
//     pipe := conn.Pipe(pipeline)
//     res := []bson.M{}
//     err := pipe.All(&res)
//     if err != nil {
//         return []models.Histogram{}, errors.New("No user found with this ID")
//     }
//     fmt.Println(res)
//
//     return []models.Histogram{}, nil
// }
