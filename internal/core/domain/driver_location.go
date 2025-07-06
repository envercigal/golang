package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type GeoJSONPoint struct {
	Type        string    `bson:"type" json:"type"`
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [lat, lon]
}

type DriverLocation struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DriverID int                `bson:"driver_id"         json:"driver_id"`
	Location GeoJSONPoint       `bson:"location"          json:"location"`
	Updated  time.Time          `bson:"updated_at"        json:"updated_at"`
}
