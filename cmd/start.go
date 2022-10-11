package cmd

import (
	"time"

	"github.com/buonotti/bus-stats-api/api"
	apiV1 "github.com/buonotti/bus-stats-api/api/v1"
	"github.com/buonotti/bus-stats-api/docs"
	"github.com/buonotti/bus-stats-api/services"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
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
			services.Env = services.Development
		} else {
			services.Env = services.Production
		}
		loadConfig()
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

/*
Entry point for the api
*/
func runApi(cmd *cobra.Command, args []string) {
	docs.SwaggerInfo.BasePath = viper.GetString(services.ConfigValue("api.base_path"))
	gin.SetMode(viper.GetString(services.ConfigValue("gin.{env}.mode")))
	trustedProxies := viper.GetStringSlice(services.ConfigValue("gin.{env}.trusted_proxies"))

	store := persist.NewMemoryStore(2 * time.Minute)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.SetTrustedProxies(trustedProxies)
	router.GET("/health", api.HealthEndpoint)

	v1 := router.Group("/api/v1")

	apiV1.MapRoutes(v1, store)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080")
}
