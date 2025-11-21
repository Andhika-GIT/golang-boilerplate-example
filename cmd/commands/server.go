package commands

import (
	"app/internal/config"
	"app/internal/infrastructure"
	"app/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server",
	Long:  "Start the HTTP server to handle API requests",
	Run:   runServer,
}

func init() {
	ServerCmd.Flags().String("port", "3005", "Server Port")
	ServerCmd.Flags().String("env", "", "Environment file path (default from APP_ENV_FILE)")

}

func runServer(cmd *cobra.Command, args []string) {

	// LOAD ENV
	envFile := os.Getenv("APP_ENV_FILE")

	if flagEnv, _ := cmd.Flags().GetString("env"); flagEnv != "" {
		envFile = flagEnv
	}

	var cfg *config.Config

	if envFile != "" {
		cfg = config.Load(envFile)
	} else {
		cfg = config.Load()
	}

	// LOAD APP PORT
	if envPort := os.Getenv("APP_ENV_PORT"); envPort != "" {
		cfg.Server.Port = envPort
	}

	if flagPort, _ := cmd.Flags().GetString("port"); flagPort != "" {
		cfg.Server.Port = flagPort
	}

	connections, err := infrastructure.InitializeConnections(cfg)

	if err != nil {
		log.Fatal(err)
	}

	deps, err := infrastructure.InitializeDependencies(connections, cfg)

	router := gin.Default()
	routes.NewRouteRegistry(router, deps.ControllerRegistry)

	err = router.Run(":" + cfg.Server.Port)

	if err != nil {
		log.Fatal(err)
	}
}
