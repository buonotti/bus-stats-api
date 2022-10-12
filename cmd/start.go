package cmd

import (
	"fmt"
	"github.com/buonotti/bus-stats-api/controllers"
	"github.com/buonotti/bus-stats-api/util"
	"os/exec"
	"runtime"
	"time"

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
		} else {
			util.Env = util.Production
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
		panic(err)
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
	cmd := exec.Command(surrealExe)
	mode := viper.GetString(util.ConfigValue("database.{env}.mode"))
	user := viper.GetString(util.ConfigValue("database.{env}.user"))
	pass := viper.GetString(util.ConfigValue("database.{env}.pass"))
	fmt.Println(mode, pass)
	cmd.Args = []string{surrealExe, "start", mode, "-u", user, "-p", pass}
	fmt.Println(cmd.Args)
	go func() {
		_ = cmd.Run()
	}()
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
	router.SetTrustedProxies(trustedProxies)
	router.GET("/health", controllers.HealthEndpoint)

	v1 := router.Group("/api/v1")

	apiV1.MapRoutes(v1, store)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080").Error()
}
