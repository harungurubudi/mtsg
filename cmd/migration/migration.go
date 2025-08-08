package migration

import (
	"fmt"
	"log"

	"github.com/harungurubudi/mtsg/internal/di/provider"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

var (
	// Root command for migrations
	RootCmd = &cobra.Command{
		Use:   "migration",
		Short: "Database migration commands",
		Long:  "Manage database migrations using sql-migrate",
	}

	// Up command
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run database migrations up",
		Long:  "Apply all pending migrations to the database",
		Run:   runMigrationUp,
	}

	// Down command
	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Run database migrations down",
		Long:  "Rollback the last migration from the database",
		Run:   runMigrationDown,
	}
)

func init() {
	RootCmd.AddCommand(upCmd)
	RootCmd.AddCommand(downCmd)
}

// runMigrationUp executes migrations up
func runMigrationUp(cmd *cobra.Command, args []string) {
	fmt.Println("🔄 Running database migrations up...")

	// Get database connection
	db := getDatabaseConnection()

	// Create migration source
	migrations := &migrate.FileMigrationSource{
		Dir: "migration",
	}

	// Execute migrations
	applied, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("❌ Failed to run migrations up: %v", err)
	}

	fmt.Printf("✅ Successfully applied %d migrations\n", applied)
}

// runMigrationDown executes migrations down
func runMigrationDown(cmd *cobra.Command, args []string) {
	fmt.Println("🔄 Rolling back database migrations...")

	// Get database connection
	db := getDatabaseConnection()

	// Create migration source
	migrations := &migrate.FileMigrationSource{
		Dir: "migration",
	}

	// Execute migrations down
	rolledBack, err := migrate.Exec(db.DB, "postgres", migrations, migrate.Down)
	if err != nil {
		log.Fatalf("❌ Failed to rollback migrations: %v", err)
	}

	fmt.Printf("✅ Successfully rolled back %d migrations\n", rolledBack)
}

// getDatabaseConnection provides database connection using DI
func getDatabaseConnection() *sqlx.DB {
	// Get configuration
	cfg := provider.ProvideConfig()

	// Get database connection
	db := provider.ProvideSqlx(cfg)

	return db
}
