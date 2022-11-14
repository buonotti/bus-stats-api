package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
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

var isDev bool

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the API and the database",
	Run: func(cmd *cobra.Command, args []string) {
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
		startDatabase()

		// start the web api
		startApi(cmd, args)
	},
}

func init() {
	startCmd.Flags().BoolVar(&isDev, "dev", false, "Run the api in development mode")
	rootCmd.AddCommand(startCmd)
}

/*
Starts the surrealDb database in a goroutine. The executable is searched in ./bin. If the database is not reachable the
api tries three times to connect to it. If the database schema is not created the cli also creates the table
*/
func startDatabase() {
	surrealExe := viper.GetString("database.executable")

	if runtime.GOOS == "windows" {
		surrealExe = surrealExe + ".exe"
	}
	log.Debug(fmt.Sprintf("Db executable is: %s", surrealExe))
	cmd := exec.Command(surrealExe)
	mode := viper.GetString(config.Get("database.{env}.mode"))
	user := viper.GetString(config.Get("database.{env}.user"))
	pass := viper.GetString(config.Get("database.{env}.pass"))
	cmd.Args = []string{surrealExe, "start", "--user", user, "--pass", pass, mode}
	go func() {
		err := cmd.Run()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()

	for i := 3; i >= 1 && !isDbOnline(); i-- {
		log.Warn(fmt.Sprintf("database seems not to be online. retrying to connect. tries left: %d", i))
		cmd := exec.Command("sleep", "2")
		err := cmd.Run()
		if err != nil {
			log.Error("error while waiting for db")
			os.Exit(1)
		}
	}

	if !isDbOnline() {
		log.Error("could not read database")
		os.Exit(1)
	}

	log.Info(fmt.Sprintf("started database with authentication in %s", mode))
	isDefined := viper.GetBool("database.generated")
	if !isDefined {
		err := surreal.ScaffoldDB()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("generated database tables")

		viper.Set("database.generated", true)
		err = viper.WriteConfig()
		if err != nil {
			log.Error(err)
		}
	}
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

func isDbOnline() bool {
	return surreal.PingDB() == nil
}
