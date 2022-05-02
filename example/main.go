package main

import (
	"context"
	"encoding/json"
	"fmt"

	fluffy_paginate "github.com/fluffy-bunny/mongo-go-pagination-generics"
	paginate "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product struct
type Product struct {
	Id       string  `json:"_id,omitempty" bson:"_id"`
	Name     string  `json:"name,omitempty" bson:"name"`
	Quantity float64 `json:"quantity,omitempty" bson:"quantity"`
	Price    float64 `json:"price,omitempty" bson:"price"`
}

func insertExamples(db *mongo.Database) (insertedIds []interface{}, err error) {
	var data []interface{}
	for i := 0; i < 30; i++ {
		data = append(data, bson.M{
			"name":     fmt.Sprintf("product-%d", i),
			"quantity": float64(i),
			"price":    float64(i*10 + 5),
		})
	}
	result, err := db.Collection("products").InsertMany(
		context.Background(), data)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

var dbConnection *mongo.Database

func main() {
	// Establishing mongo db connection
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://127.0.0.1:27017/"))
	if err != nil {
		panic(err)
	}
	dbConnection = client.Database("myaggregate")
	_, insertErr := insertExamples(client.Database("myaggregate"))
	if insertErr != nil {
		panic(insertErr)
	}
	var normal = func(page int64, limit int64) string {
		filter := bson.M{}

		collection := dbConnection.Collection("products")
		projection := bson.D{
			{"name", 1},
			{"quantity", 1},
		}
		// Querying paginated data
		// If you want to do some complex sort like sort by score(weight) for full text search fields you can do it easily
		// sortValue := bson.M{
		//		"$meta" : "textScore",
		//	}
		// paginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("score", sortValue)...
		products, paginatedData, err := fluffy_paginate.New[Product](collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Sort("quantity", -1).Select(projection).Filter(filter).Find()
		if err != nil {
			panic(err)
		}

		payload := struct {
			Data       []Product               `json:"data"`
			Pagination paginate.PaginationData `json:"pagination"`
		}{
			Pagination: paginatedData.Pagination,
			Data:       products,
		}
		json, _ := json.MarshalIndent(payload, "", "    ")
		return string(json)
	}
	fmt.Println(normal(1, 10))

	var agg = func(page int64, limit int64) string {

		collection := dbConnection.Collection("products")

		//Example for Aggregation

		//match query
		match := bson.M{"$match": bson.M{"quantity": bson.M{"$gt": 0}}}
		//
		//group query
		projectQuery := bson.M{"$project": bson.M{"_id": 1, "name": 1, "quantity": 1}}

		// you can easily chain function and pass multiple query like here we are passing match
		// query and projection query as params in Aggregate function you cannot use filter with Aggregate
		// because you can pass filters directly through Aggregate param
		aggPaginatedData, err := paginate.New(collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Aggregate(match, projectQuery)
		if err != nil {
			panic(err)
		}

		var aggProductList []Product
		for _, raw := range aggPaginatedData.Data {
			var product *Product
			if marshallErr := bson.Unmarshal(raw, &product); marshallErr == nil {
				aggProductList = append(aggProductList, *product)
			}

		}

		s, aggPaginatedData, err := fluffy_paginate.New[Product](collection).Context(ctx).Limit(limit).Page(page).Sort("price", -1).Aggregate(match, projectQuery)
		fmt.Println(s)
		payload := struct {
			Data       []Product               `json:"data"`
			Pagination paginate.PaginationData `json:"pagination"`
		}{
			Pagination: aggPaginatedData.Pagination,
			Data:       aggProductList,
		}
		json, _ := json.MarshalIndent(payload, "", "    ")
		return string(json)
	}
	fmt.Println(agg(1, 10))

}
