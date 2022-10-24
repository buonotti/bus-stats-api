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

	"github.com/buonotti/bus-stats-api/controllers"
	"github.com/buonotti/bus-stats-api/util"

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
	Short: "Start the API",
	Run: func(cmd *cobra.Command, args []string) {
		if isDev {
			util.Env = util.Development
			log.SetLevel(log.DebugLevel)
		} else {
			util.Env = util.Production
			log.SetLevel(log.InfoLevel)
		}
		loadConfig()
		configLogger()
		startDatabase()
		runApi(cmd, args)

	},
}

func init() {
	startCmd.Flags().BoolVar(&isDev, "dev", false, "Run the api in development mode")
	rootCmd.AddCommand(startCmd)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func configLogger() {
	if util.Env == util.Development {
		log.SetFormatter(&log.TextFormatter{
			ForceColors:  true,
			PadLevelText: true,
		})
	} else {
		log.SetFormatter(&log.JSONFormatter{})
	}
}

func startDatabase() {
	surrealExe := "./bin/surreal-v1.0.0-beta.8." + runtime.GOOS + "-" + runtime.GOARCH

	if runtime.GOOS == "windows" {
		surrealExe = surrealExe + ".exe"
	}
	log.Debug(fmt.Sprintf("Db executable is: %s", surrealExe))
	cmd := exec.Command(surrealExe)
	mode := viper.GetString(util.ConfigValue("database.{env}.mode"))
	user := viper.GetString(util.ConfigValue("database.{env}.user"))
	pass := viper.GetString(util.ConfigValue("database.{env}.pass"))
	cmd.Args = []string{surrealExe, "start", "-u", user, "-p", pass, mode}
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
		_, err := util.RestClient.R().SetBody(`
DEFINE TABLE user SCHEMAFULL;
DEFINE FIELD email ON user TYPE string;
DEFINE FIELD password ON user TYPE string;
DEFINE FIELD image ON user TYPE object;
DEFINE FIELD image.name ON user TYPE string;
DEFINE FIELD image.type ON user TYPE string;
DEFINE TABLE stop SCHEMAFULL;
DEFINE FIELD name ON stop TYPE string;
DEFINE FIELD location ON stop TYPE array;
DEFINE TABLE line SCHEMAFULL;
DEFINE FIELD name ON line TYPE string;
`).Post(util.DatabaseUrl())
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("generated dabase tables")

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
func runApi(cmd *cobra.Command, args []string) {
	docs.SwaggerInfo.BasePath = viper.GetString(util.ConfigValue("api.base_path"))
	gin.SetMode(viper.GetString(util.ConfigValue("gin.{env}.mode")))
	trustedProxies := viper.GetStringSlice(util.ConfigValue("gin.{env}.trusted_proxies"))

	store := persist.NewMemoryStore(2 * time.Minute)

	router := gin.New()
	router.Use(util.LogrusLogger())
	router.Use(gin.Recovery())
	if util.Env == util.Development {
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
		if viper.GetString(util.ConfigValue("database.{env}.mode")) == "memory" {
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
	_, err := util.RestClient.R().SetBody("INFO FOR DB;").Post(util.DatabaseUrl())
	if err != nil {
		return false
	}
	return true
}
