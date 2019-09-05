package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"example.com/go-graphql-api/musicutil"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
)

var ginLambda *ginadapter.GinLambda

func createGinRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/pets", getPets)
	r.GET("/pets/:id", getPet)
	r.POST("/pets", createPet)
	r.Any("/discs", graphQlHandler())

	return r
}

// Handler is the main entry point for Lambda. Receives a proxy request and
// returns a proxy response
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if ginLambda == nil {
		// stdout and stderr are sent to AWS CloudWatch Logs
		log.Printf("Gin cold start")
		r := createGinRouter()
		ginLambda = ginadapter.New(r)
	}

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// If the "local" arg is provided, then run gin directly
	if len(os.Args) > 1 && os.Args[1] == "local" {
		r := createGinRouter()
		r.Run(":4000")
		return
	} else {
		lambda.Start(Handler)
	}
}

func getPets(c *gin.Context) {
	limit := 10
	if c.Query("limit") != "" {
		newLimit, err := strconv.Atoi(c.Query("limit"))
		if err != nil {
			limit = 10
		} else {
			limit = newLimit
		}
	}
	if limit > 50 {
		limit = 50
	}
	pets := make([]Pet, limit)

	for i := 0; i < limit; i++ {
		pets[i] = getRandomPet()
	}

	c.JSON(200, pets)
}

func getPet(c *gin.Context) {
	petID := c.Param("id")
	randomPet := getRandomPet()
	randomPet.ID = petID
	c.JSON(200, randomPet)
}

func createPet(c *gin.Context) {
	newPet := Pet{}
	err := c.BindJSON(&newPet)

	if err != nil {
		return
	}

	newPet.ID = getUUID()
	c.JSON(http.StatusAccepted, newPet)
}

func graphQlHandler() gin.HandlerFunc {
	// Creates a GraphQL-go HTTP handler with the defined schema
	h := handler.New(&handler.Config{
		Schema:   &musicutil.MusicSchema, // &schema.Schema,
		Pretty:   true,
		GraphiQL: true,
	})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
