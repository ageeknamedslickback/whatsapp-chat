package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ageeknamedslickback/whatsapp-chat/domain"
	database "github.com/ageeknamedslickback/whatsapp-chat/infrastructure/db"
	"github.com/ageeknamedslickback/whatsapp-chat/presentation/graph"
	"github.com/ageeknamedslickback/whatsapp-chat/presentation/graph/generated"
	"github.com/ageeknamedslickback/whatsapp-chat/presentation/rest"
	"github.com/ageeknamedslickback/whatsapp-chat/usecases"
	"github.com/gorilla/websocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(domain.Message{})

	d := database.NewMessageRepository(db)
	s := usecases.NewMessageService(*d)
	r := graph.NewResolver(*s)
	h := rest.NewRestHandlers(*s, *r)

	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: r},
		),
	)

	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		KeepAlivePingInterval: 10 * time.Second,
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	http.HandleFunc("/message", h.IncomingMessage)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
