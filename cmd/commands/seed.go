package commands

import (
	"app/internal/config"
	"app/internal/database/seeders"
	"app/internal/infrastructure"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// Command flags
var spesificSeeder string
var listSeeder bool

// SeederCmd handles database seeding operations
var SeederCmd = &cobra.Command{
	Use:   "Seed",
	Short: "Start Database Seed",
	Run:   runSeeder,
}

func init() {
	SeederCmd.Flags().StringVarP(&spesificSeeder, "name", "n", "", "run specific seeder")
	SeederCmd.Flags().BoolVarP(&listSeeder, "list", "l", false, "get all seeder list")
}

// runSeeder executes seeder command with provided flags
func runSeeder(cmd *cobra.Command, args []string) {
	// Load environment configuration
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

	// Initialize database connection
	connections, err := infrastructure.InitializeConnections(cfg)
	if err != nil {
		log.Printf("error when running seeder : %s", err.Error())
	}

	registry := seeders.InitSeeders()

	// Handle list flag - show available seeders
	if listSeeder {
		fmt.Println("Available seeders:")
		for _, name := range registry.GetAllSeeders() {
			fmt.Printf("  - %s\n", name)
		}
		return
	}

	// Handle specific seeder flag - run single seeder
	if spesificSeeder != "" {
		err := registry.RunSpecific(connections.DB, spesificSeeder)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}

	// Default behavior - run all seeders
	fmt.Println("Running all seeders...")
	err = registry.RunAll(connections.DB)
	if err != nil {
		log.Print(err.Error())
	}

	log.Println("\n✓ All seeders completed successfully!")
}
