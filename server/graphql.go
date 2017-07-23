package server

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/udryan10/grpc-rest-example/generated"
)

var GraphQLMapsType *graphql.Object
var GraphQLMarkerType *graphql.Object

func init() {
	GraphQLMapsType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Query",
		Description: "",
		Fields: (graphql.FieldsThunk)(func() graphql.Fields {
			return graphql.Fields{
				"markers": &graphql.Field{
					Type:        graphql.NewList(GraphQLMarkerType),
					Description: "",
					Args: graphql.FieldConfigArgument{
						"point": &graphql.ArgumentConfig{
							Description: "target a specific marker",
							Type:        graphql.NewList(graphql.Int),
						},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						maps, err := loadMapFromDisk()
						if err != nil {
							return &generated.Markers{}, fmt.Errorf("field markers not resolved")
						}
						// return all
						if p.Args["point"] == nil {
							return maps.Markers, nil
						}

						selectedMarkers := []*generated.Marker{}
						for _, marker := range maps.Markers {

							for _, arg := range p.Args["point"].([]interface{}) {
								if int32(arg.(int)) == marker.Point {
									selectedMarkers = append(selectedMarkers, marker)
								}
							}
						}

						return selectedMarkers, nil
					},
				},
			}
		}),
	})

	GraphQLMarkerType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Marker",
		Description: "",
		Fields: (graphql.FieldsThunk)(func() graphql.Fields {
			return graphql.Fields{
				"point": &graphql.Field{
					Type:        graphql.Int,
					Description: "",

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.Point, nil
						}
						return nil, fmt.Errorf("field point not resolved")
					},
				},
				"homeTeam": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.HomeTeam, nil
						}
						return nil, fmt.Errorf("field homeTeam not resolved")
					},
				},
				"awayTeam": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.AwayTeam, nil
						}
						return nil, fmt.Errorf("field awayTeam not resolved")
					},
				},
				"markerImage": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.MarkerImage, nil
						}
						return nil, fmt.Errorf("field markerImage not resolved")
					},
				},
				"information": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.Information, nil
						}
						return nil, fmt.Errorf("field information not resolved")
					},
				},
				"fixture": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.Fixture, nil
						}

						return nil, fmt.Errorf("field fixture not resolved")
					},
				},
				"capacity": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.Capacity, nil
						}
						return nil, fmt.Errorf("field capacity not resolved")
					},
				},
				"previousScore": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.PreviousScore, nil
						}
						return nil, fmt.Errorf("field previousScore not resolved")
					},
				},
				"tv": &graphql.Field{
					Type:        graphql.String,
					Description: "",
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						obj, ok := p.Source.(*generated.Marker)
						if ok {
							return obj.Tv, nil
						}
						return nil, fmt.Errorf("field tv not resolved")
					},
				},
			}
		}),
	})
}
