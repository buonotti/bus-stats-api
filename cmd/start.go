package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/buonotti/bus-stats-api/config"
	"github.com/buonotti/bus-stats-api/controllers"
	"github.com/buonotti/bus-stats-api/logging"
	"github.com/buonotti/bus-stats-api/middleware"
	"github.com/buonotti/bus-stats-api/surreal"

	apiV1 "github.com/buonotti/bus-stats-api/controllers/v1"
	"github.com/buonotti/bus-stats-api/docs"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// variable holding the value of the --dev flags
var isDev bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API and the database",
	Long:  `This command starts the api and the database. The api runs on http://localhost:8080 by default.`,
	Run: func(cmd *cobra.Command, args []string) {

		// set the environment
		if isDev {
			config.Env = config.Development
			log.SetLevel(log.DebugLevel)
		} else {
			config.Env = config.Production
			log.SetLevel(log.InfoLevel)
		}

		// init config library and load config values
		config.Setup()

		// config logger formats
		logging.Setup()

		// run the database
		surreal.Exec()

		// start the web api
		startApi(cmd, args)
	},
}

func init() {
	startCmd.Flags().BoolVar(&isDev, "dev", false, "Run the api in development mode")
	rootCmd.AddCommand(startCmd)
}

/*
Entry point for the api
*/
func startApi(cmd *cobra.Command, args []string) {
	docs.SwaggerInfo.BasePath = viper.GetString(config.Get("api.base_path"))
	gin.SetMode(viper.GetString(config.Get("gin.{env}.mode")))
	trustedProxies := viper.GetStringSlice(config.Get("gin.{env}.trusted_proxies"))

	store := persist.NewMemoryStore(2 * time.Minute)

	router := gin.New()
	router.Use(logging.LogrusLogger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())
	if config.Env == config.Development {
		router.SetTrustedProxies(trustedProxies)
	}
	router.GET("/health", controllers.HealthEndpoint)

	v1 := router.Group("/api/v1")

	apiV1.MapRoutes(v1, store)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			panic(err) // TODO
		}
		log.Info("shut down server")
		if viper.GetString(config.Get("database.{env}.mode")) == "memory" {
			viper.Set("database.generated", false)
			err := viper.WriteConfig()
			if err != nil {
				fmt.Println(err)
			}
		}
		close(idleConnsClosed)
	}()

	log.Info("server listening")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}

	<-idleConnsClosed
}
