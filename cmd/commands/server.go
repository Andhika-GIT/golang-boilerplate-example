package commands

import (
	"app/internal/config"
	"app/internal/infrastructure"
	"app/internal/routes"
	"log"

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
	ServerCmd.Flags().String("port", "8080", "Server Port")
}

func runServer(cmd *cobra.Command, args []string) {
	cfg := config.Load()

	if port, _ := cmd.Flags().GetString("port"); port != "" {
		cfg.Server.Port = port
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
