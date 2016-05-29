package models

import (
    "time"
    "gopkg.in/mgo.v2/bson"
)

// need capitalized to expoert the variables
// apply alias to use lower case names when delivering json
type ImageInfo struct {
        Id              bson.ObjectId   `json:"id" bson:"id"`
        UserID          string          `json:"user_id" bson:"user_id"`
        FileName        string          `json:"file_name" bson:"file_name"`
        ImageURL        string          `json:"image_url" bson:"image_url"`
        Timestamp       time.Time       `json:"timestamp" bson:"timestamp"`
        Histogram       Histogram       `json:"histogram" bson:"histogram"`
}

type Histogram [16]int
// append json/mgo struct tag
// to instruct how to store the user info
