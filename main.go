package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

//Person struct
type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
}

//People for Slice of Struct Person
type People []Person

var schema graphql.Schema
var people People

var personType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Person",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.ID,
		},
		"firstName": &graphql.Field{
			Type: graphql.String,
		},
		"lastName": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"gender": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"people": &graphql.Field{
			Type: graphql.NewList(personType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return people, nil
			},
		},
	},
})

func init() {
	f, _ := os.Open("./data.csv")
	r := csv.NewReader(bufio.NewReader(f))

	records, _ := r.ReadAll()
	for i, record := range records {
		if i == 0 {
			continue
		}

		people = append(people, Person{
			ID:        record[0],
			FirstName: record[1],
			LastName:  record[2],
			Email:     record[3],
			Gender:    record[4],
		})
	}
}

func main() {
	schema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		Playground: true, // False or delete (close gui)
	})

	http.Handle("/graphql", h)

	log.Println("Listening on port :8080")
	http.ListenAndServe(":8080", nil)
}
