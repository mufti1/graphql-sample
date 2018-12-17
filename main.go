package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/graphql-go/graphql"
)

// SherlockCase define object of sherlock holmes cases
type SherlockCase struct {
	ID       int64  `json:"id"`
	Name     string `json:"case_name"`
	Time     string `json:"time"`
	Location string `json:"location"`
}

var sherlockCases []SherlockCase

// definition of graphql data
var caseType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Case",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"time": &graphql.Field{
				Type: graphql.String,
			},
			"location": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// get case by id
			"case": &graphql.Field{
				Type:        caseType,
				Description: "Get Case By ID",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						for _, sherlockCase := range sherlockCases {
							if int(sherlockCase.ID) == id {
								return sherlockCase, nil
							}
						}
					}
					return nil, nil
				},
			},
			//get all case
			"allcase": &graphql.Field{
				Type:        graphql.NewList(caseType),
				Description: "Get all case",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return sherlockCases, nil
				},
			},
		},
	},
)

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"add": &graphql.Field{
			Type:        caseType,
			Description: "add new case",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"time": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"location": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				sherlockCase := SherlockCase{
					ID:       int64(rand.Intn(100000)), // generate random ID
					Name:     p.Args["name"].(string),
					Time:     p.Args["time"].(string),
					Location: p.Args["location"].(string),
				}
				sherlockCases = append(sherlockCases, sherlockCase)
				return sherlockCase, nil
			},
		},

		"update": &graphql.Field{
			Type:        caseType,
			Description: "update case",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"time": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"location": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, _ := p.Args["id"].(int)
				name, nameOk := p.Args["name"].(string)
				time, timeOk := p.Args["time"].(string)
				location, locationOk := p.Args["location"].(string)
				sherlockCase := SherlockCase{}
				for k, v := range sherlockCases {
					if int64(id) == v.ID {
						if nameOk {
							sherlockCases[k].Name = name
						}
						if timeOk {
							sherlockCases[k].Time = time
						}
						if locationOk {
							sherlockCases[k].Location = location
						}
						sherlockCase = sherlockCases[k]
						break
					}
				}
				return sherlockCase, nil
			},
		},

		"delete": &graphql.Field{
			Type:        caseType,
			Description: "Delete case by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				sherlockCase := SherlockCase{}
				for i, p := range sherlockCases {
					if int64(id) == p.ID {
						sherlockCase = sherlockCases[i]
						// Remove from product list
						sherlockCases = append(sherlockCases[:i], sherlockCases[i+1:]...)
					}
				}
				return sherlockCase, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func initSherlockCase(p *[]SherlockCase) {
	case1 := SherlockCase{
		ID:       1,
		Name:     "A Study in Scarlet",
		Time:     "March 1881",
		Location: "london",
	}
	case2 := SherlockCase{
		ID:       2,
		Name:     "A Scandal in Bohemia",
		Time:     "20 March 1888",
		Location: "london",
	}
	*p = append(*p, case1, case2)
}

func main() {
	//init data
	initSherlockCase(&sherlockCases)

	http.HandleFunc("/sherlockcase", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("listening on server 8001")
	http.ListenAndServe(":8001", nil)
}
