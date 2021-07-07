package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"

	"github.com/Jimeux/go-redis/app"
	"github.com/Jimeux/go-redis/cache"
)

func main() {
	dynamoRepo := app.NewRepository(dynamodb.NewFromConfig(app.AWSConfig()))
	cacheRepo := cache.NewRepository(redis.NewClient(&redis.Options{
		Addr:     "localhost:26379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}))
	svc := app.NewService(dynamoRepo, cacheRepo)
	handler := app.NewHandler(svc)

	r := chi.NewRouter()
	r.Route("/messages", func(r chi.Router) {
		r.Get("/", handler.FindMessages)
		r.Post("/", handler.SendMessage)
	})
	r.Route("/reactions", func(r chi.Router) {
		r.Post("/", handler.SendReactions)
	})
	r.Route("/threads", func(r chi.Router) {
		r.Post("/", handler.CreateThread)
	})

	log.Println("Listening on localhost:8888")
	if err := http.ListenAndServe(":8888", r); err != nil {
		log.Println(err)
	}
}
