package seeders

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Seeder interface defines contract for all seeders
type Seeder interface {
	Seed(db *gorm.DB) error
	GetName() string
}

// SeederRegistry manages collection of seeders
type SeederRegistry struct {
	seeders []Seeder
}

// NewSeederRegistry creates new registry instance
func NewSeederRegistry() *SeederRegistry {
	return &SeederRegistry{
		seeders: make([]Seeder, 0),
	}
}

// Register adds seeder to registry
func (r *SeederRegistry) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// RunAll executes all registered seeders
func (r *SeederRegistry) RunAll(db *gorm.DB) error {
	for _, seeder := range r.seeders {
		log.Printf("Running seeder: %s\n", seeder.GetName())
		err := seeder.Seed(db)
		if err != nil {
			return fmt.Errorf("error running seeder %s: %w", seeder.GetName(), err)
		}
		log.Printf("✓ Seeder %s completed\n", seeder.GetName())
	}
	return nil
}

// RunSpecific executes single seeder by name
func (r *SeederRegistry) RunSpecific(db *gorm.DB, name string) error {
	for _, seeder := range r.seeders {
		if seeder.GetName() == name {
			log.Printf("Running seeder: %s\n", name)
			if err := seeder.Seed(db); err != nil {
				return fmt.Errorf("error running seeder %s: %w", name, err)
			}
			log.Printf("✓ Seeder %s completed\n", name)
			return nil
		}
	}
	return fmt.Errorf("seeder '%s' not found", name)
}

// GetAllSeeders returns names of all registered seeders
func (r *SeederRegistry) GetAllSeeders() []string {
	names := make([]string, len(r.seeders))
	for i, seeder := range r.seeders {
		names[i] = seeder.GetName()
	}
	return names
}

// InitSeeders initializes registry with all seeders
func InitSeeders() *SeederRegistry {
	registry := NewSeederRegistry()
	registry.Register(&UserSeeder{})
	return registry
}
