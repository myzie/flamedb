package main

import (
	"crypto/rsa"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/myzie/flamedb/database"
	"github.com/myzie/flamedb/restapi"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/service"
	"github.com/namsral/flag"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func getJwk(url string) *rsa.PublicKey {
	jwks, err := service.GetKeySet(url)
	if err != nil {
		log.Fatalln("Failed to retrieve JWKS:", err)
	}
	keys, err := jwks.GetRSAPublicKeys()
	if err != nil {
		log.Fatalln("Failed to decode JWKS:", err)
	}
	if len(keys) == 0 {
		log.Fatalln("JWKS is empty")
	}
	return keys[0]
}

func main() {

	var (
		addr    string
		port    int
		origins string
		keyURL  string
	)

	flag.StringVar(&addr, "addr", "localhost", "Listen address")
	flag.StringVar(&origins, "origins", "", "List of origins to allow")
	flag.IntVar(&port, "port", 8000, "Listen port")
	flag.StringVar(&keyURL, "jwks", "", "JWKS URL")
	flag.Parse()

	var jwk *rsa.PublicKey
	if keyURL != "" {
		jwk = getJwk(keyURL)
		log.Println("Loaded JWKS", keyURL)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewFlamedbAPI(swaggerSpec)
	api.Logger = log.Infof
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Host = addr
	server.Port = port

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(origins, ","),
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	})

	dbSettings := database.GetSettings()
	log.Infof("DB settings: %+v\n", dbSettings)

	gormDB, err := database.Connect(dbSettings)
	if err != nil {
		log.Fatalln(err)
	}
	if err := gormDB.AutoMigrate(&database.Record{}).Error; err != nil {
		log.Fatalln(err)
	}

	service.New(service.Opts{
		API:   api,
		Flame: database.NewFlame(gormDB),
		Key:   jwk,
	})

	server.SetHandler(corsMiddleware.Handler(api.Serve(nil)))

	log.Printf("Listening at %s:%d\n", addr, port)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
