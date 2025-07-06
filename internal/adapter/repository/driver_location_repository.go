package mongo

import (
	"context"
	"errors"
	"github.com/envercigal/golang/internal/core/domain"
	"github.com/envercigal/golang/internal/core/port"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type driverLocationRepo struct {
	coll *mongo.Collection
}

func NewDriverLocationRepo(db *mongo.Database) port.DriverLocationRepository {
	_, err := db.Collection("driver_locations").Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "location", Value: "2dsphere"}},
		},
	)
	if err != nil {
		return nil

	}

	return &driverLocationRepo{
		coll: db.Collection("driver_locations"),
	}
}

func (r *driverLocationRepo) Create(ctx context.Context, dl *domain.DriverLocation) (*domain.DriverLocation, error) {
	res, err := r.coll.InsertOne(ctx, dl)
	if err != nil {
		return nil, err
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert InsertedID to ObjectID")
	}
	dl.ID = oid
	return dl, nil
}

func (r *driverLocationRepo) BulkCreate(ctx context.Context, batch []*domain.DriverLocation) error {
	docs := make([]interface{}, len(batch))
	for i, dl := range batch {
		docs[i] = dl
	}

	// Synchronized, acknowledged insert
	_, err := r.coll.InsertMany(
		ctx,
		docs,
		options.InsertMany().
			SetOrdered(false).                 // hata olsa bile devam et
			SetBypassDocumentValidation(true), // validasyon maliyetini atla
	)
	return err
}

func (r *driverLocationRepo) FindNearest(ctx context.Context, lon, lat float64) (*domain.DriverLocation, error) {
	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{lon, lat},
				},
			},
		},
	}

	var dl domain.DriverLocation
	if err := r.coll.FindOne(ctx, filter).Decode(&dl); err != nil {
		return nil, err
	}
	return &dl, nil
}
