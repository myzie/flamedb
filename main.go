package main

import (
	"os"
	"strings"

	"github.com/myzie/flamedb/database"

	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/myzie/flamedb/restapi"
	"github.com/myzie/flamedb/restapi/operations"
	"github.com/myzie/flamedb/service"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

func main() {

	var opts struct {
		Addr    string `short:"a" long:"addr" default:"localhost" description:"Listen address"`
		Port    int    `short:"p" long:"port" default:"8000" description:"Listen port"`
		Origins string `long:"origins" description:"List of origins to allow"`
	}

	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewFlamedbAPI(swaggerSpec)
	api.Logger = log.Infof
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Host = opts.Addr
	server.Port = opts.Port

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   strings.Split(opts.Origins, ","),
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
	})

	server.SetHandler(corsMiddleware.Handler(api.Serve(nil)))

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
