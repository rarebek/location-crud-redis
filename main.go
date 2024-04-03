package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()

	location := Location{
		Latitude:  41.299496,
		Longitude: 69.240074,
		Name:      "Tashkent",
	}

	err := SaveLocation(ctx, rdb, "tash", location)
	if err != nil {
		fmt.Println("Error saving location:", err)
	}

	savedLocation, err := GetLocation(ctx, rdb, "tash")
	if err != nil {
		fmt.Println("Error retrieving location:", err)
	} else {
		fmt.Println("Retrieved location:", savedLocation)
	}

	updatedLocation := Location{
		Latitude:  40.993599,
		Longitude: 71.677452,
		Name:      "Namangan",
	}
	err = UpdateLocation(ctx, rdb, "tash", updatedLocation)
	if err != nil {
		fmt.Println("Error updating location:", err)
	}

	updatedSavedLocation, err := GetLocation(ctx, rdb, "tash")
	if err != nil {
		fmt.Println("Error retrieving updated location:", err)
	} else {
		fmt.Println("Updated location:", updatedSavedLocation)
	}

	err = DeleteLocation(ctx, rdb, "tash")
	if err != nil {
		fmt.Println("Error deleting location:", err)
	}

	_, err = GetLocation(ctx, rdb, "tash")
	if err != nil {
		fmt.Println("Location does not exist after deletion")
	}
}

func SaveLocation(ctx context.Context, rdb *redis.Client, key string, location Location) error {
	return rdb.HSet(ctx, key, "latitude", location.Latitude, "longitude", location.Longitude, "name", location.Name).Err()
}

func GetLocation(ctx context.Context, rdb *redis.Client, key string) (Location, error) {
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return Location{}, err
	}

	latitude := result["latitude"]
	longitude := result["longitude"]
	name:= result["name"]

	return Location{
		Latitude:  ParseFloat64(latitude),
		Longitude: ParseFloat64(longitude),
		Name:      name,
	}, nil
}

func UpdateLocation(ctx context.Context, rdb *redis.Client, key string, location Location) error {
	return rdb.HSet(ctx, key, "latitude", location.Latitude, "longitude", location.Longitude, "name", location.Name).Err()
}

func DeleteLocation(ctx context.Context, rdb *redis.Client, key string) error {
	return rdb.Del(ctx, key).Err()
}

func ParseFloat64(s string) float64 {
	f := 0.0
	if _, err := fmt.Sscanf(s, "%f", &f); err != nil {
		return 0.0
	}
	return f
}
