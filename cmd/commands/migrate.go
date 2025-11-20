package commands

import (
	"app/internal/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

const migrationDir = "./internal/database/migrations"

var (
	// ============================================
	// MAIN COMMAND
	// ============================================
	MigrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration with golang-migrate",
	}

	// ============================================
	// SUB COMMAND
	// ============================================
	MigrateUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Run pending migrations",
		Run:   runMigrateUp,
	}

	MigrateDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Rollback migrations",
		Run:   runMigrateDown,
	}

	MigrateCreateCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create new migration files",
		Args:  cobra.ExactArgs(1),
		Run:   runMigrateCreate,
	}

	MigrateStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Run:   runMigrateStatus,
	}
)

// ============================================
// COMMAND IMPLEMENTATIONS
// ============================================

func runMigrateUp(cmd *cobra.Command, args []string) {
	log.Println("🔄 Running migrations...")
	runMigrate("up")
	log.Println("✅ Migrations completed!")
}

func runMigrateDown(cmd *cobra.Command, args []string) {
	steps := 1 // Default
	if len(args) > 0 {
		// 👇 Parse argument jika ada
		if s, err := strconv.Atoi(args[0]); err == nil && s > 0 {
			steps = s
		}
	}

	log.Printf("⏪ Rolling back %d migration(s)...", steps)
	runMigrate("down", fmt.Sprintf("%d", steps))
	log.Printf("✅ Rollback %d migration(s) completed!", steps)
}

func runMigrateCreate(cmd *cobra.Command, args []string) {
	name := args[0]
	log.Printf("📝 Creating migration: %s", name)

	// Ensure migration directory exists
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		log.Fatalf("❌ Failed to create migration directory: %v", err)
	}

	// Execute migrate create command
	createCmd := exec.Command("migrate",
		"create", "-ext", "sql",
		"-dir", "./internal/database/migrations",
		"-seq",
		name,
	)
	createCmd.Stdout = os.Stdout
	createCmd.Stderr = os.Stderr

	if err := createCmd.Run(); err != nil {
		log.Fatalf("❌ Failed to create migration: %v", err)
	}

	log.Printf("✅ Created migration: %s", name)
}

func runMigrateStatus(cmd *cobra.Command, args []string) {
	log.Println("📊 Migration status:")
	runMigrate("version")
}

// ============================================
// HELPER - DIRECT MIGRATE EXECUTION
// ============================================

func runMigrate(args ...string) {
	// Build database URL from config
	dbURL := getDatabaseURL()

	// Prepare migrate command
	migrateArgs := []string{
		"-database", dbURL,
		"-path", migrationDir,
		"-verbose",
	}
	migrateArgs = append(migrateArgs, args...)

	// Execute directly
	cmd := exec.Command("migrate", migrateArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Migration command failed: %v", err)
	}
}

// getDatabaseURL from app config
func getDatabaseURL() string {
	cfg := config.Load()

	switch cfg.Database.Type {
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.User, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port,
			cfg.Database.Name, cfg.Database.SSLMode,
		)
	case "mysql":
		return fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s",
			cfg.Database.User, cfg.Database.Password,
			cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
		)
	default:
		log.Fatalf("❌ Unsupported database: %s", cfg.Database.Type)
		return ""
	}
}

// ============================================
// INIT - REGISTER COMMANDS
// ============================================

func init() {
	MigrateCmd.AddCommand(MigrateUpCmd)
	MigrateCmd.AddCommand(MigrateDownCmd)
	MigrateCmd.AddCommand(MigrateCreateCmd)
	MigrateCmd.AddCommand(MigrateStatusCmd)

}
