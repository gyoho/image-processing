package controllers

import (
    // Standard library packages
    "fmt"
    "net/http"
    "encoding/json"
    "log"
    "time"
    "strconv"
    "errors"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"

    // Third party packages
    "github.com/julienschmidt/httprouter"
    // Use defined model in another dir
    "../models"
    "../utils"
    "../utils/structs"
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

func (ic ImageController) GetWeeklyHistograms(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    hist, err := retriveHistogramsOfWeek(ic, param.ByName("userId"))
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        histJson, _ := json.Marshal(hist)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", histJson)
    }
}

func (ic ImageController) GetMedianHistogram(rw http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
    midHistMap, err := getMedianHistogram(ic)
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        midHistJson, _ := json.Marshal(midHistMap)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", midHistJson)
    }
}

func getMedianHistogram(ic ImageController) (map[string]models.Histogram, error) {
    col, err := retriveHistogramsOfDay(ic)
    if err != nil {
        return nil, err
    }

    midHist, err := calculateMedianHistogram(col)
    if err != nil {
        return nil, err
    }

    midHistMap := map[string]models.Histogram{"median_histogram" : midHist}
    return midHistMap, nil
}

func (ic ImageController) GetUserIDWithSimilarHistogram(rw http.ResponseWriter, _ *http.Request, param httprouter.Params) {
    userIdsMap, err := getUserIDWithSimilarHistogram(ic, param)
    if err != nil {
		log.Println(err)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "plain/text")
        rw.WriteHeader(400)
        fmt.Fprintf(rw, "%s\n", err)
	} else {
        // Create response
        // Marshal provided interface into JSON structure
        userIdsJson, _ := json.Marshal(userIdsMap)
        // Write content-type, statuscode, payload
        rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(200)
        fmt.Fprintf(rw, "%s\n", userIdsJson)
    }
}

func getUserIDWithSimilarHistogram(ic ImageController, param httprouter.Params) (map[string][]string, error) {
    col, err := retriveHistogramsOfDay(ic)
    if err != nil {
        return nil, err
    }

    nSimilarElem, err := extractUserIDWithSimilarHistogram(param.ByName("userId"), param.ByName("n"), col)
    if err != nil {
        return nil, err
    }

    n, err := strconv.Atoi(param.ByName("n"))
    if err != nil {
        return nil, errors.New("Invalid parameter")
    }
    userIds := make([]string, n)
    for i, elem := range nSimilarElem {
        userIds[i] = elem.UserID
    }

    nUserIdsMap := map[string][]string{"userIds": userIds}
    return nUserIdsMap, nil
}

func retriveHistogramsOfWeek(ic ImageController, user_id string) ([]bson.M, error) {
    // make connection
    conn := ic.session.DB("image").C("image_info")
    // prepare query
    year, week := time.Now().UTC().ISOWeek()    // mongo stores timestamp in UTC
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
        return res, errors.New("No user found with this ID")
    }

    return res, nil
}

func retriveHistogramsOfDay(ic ImageController) ([]models.ImageInfo, error) {
    // make connection
    conn := ic.session.DB("image").C("image_info")
    // prepare query
    now := time.Now().UTC() // mongo stores timestamp in UTC
    year := now.Year()
    day := now.YearDay()
    pipeline := []bson.M {
        bson.M {
            "$project": bson.M {
    			"_id": 0,
                "user_id": 1,
    	        "year": bson.M{ "$year": "$timestamp"},
    	        "day": bson.M{ "$dayOfYear": "$timestamp"},
    			"histogram": 1}},
        bson.M {
            "$match": bson.M {
    	        "year": year,
    	        "day": day}},
        bson.M {
            "$project": bson.M {
                "user_id": 1,
                "histogram": 1}}}

    pipe := conn.Pipe(pipeline)
    res := []models.ImageInfo{}
    err := pipe.All(&res)
    if err != nil {
        return nil, errors.New("No document found")
    }

    return res, nil
}

func calculateMedianHistogram(res []models.ImageInfo) (models.Histogram, error) {
    mh := structs.NewMedianHistogram()
    for _, elem := range res {
        mh.AddHistogram(elem.Histogram)
    }

    midHist, err := mh.GetMedianHistogram()
    if err != nil {
        return models.Histogram{}, err
    }

    return midHist, nil
}

func extractUserIDWithSimilarHistogram(user_id string, num string, inf []models.ImageInfo) ([]models.ImageInfo, error) {
    var src models.ImageInfo
    for _, elem := range inf {
         if user_id == elem.UserID {
             src = elem
         }
    }

    sh := structs.NewSimilarHistHeap()
    n, err := strconv.Atoi(num)
    if err != nil {
        return nil, errors.New("Invalid parameter")
    }
    sh.BuildHeap(src, inf, n)
    return sh.GetNSimilarElem()
}
