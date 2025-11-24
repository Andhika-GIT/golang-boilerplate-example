package seeders

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Seeder interface {
	Seed(db *gorm.DB) error
	GetName() string
}

type SeederRegistry struct {
	seeders []Seeder
}

func NewSeederRegistry() *SeederRegistry {
	return &SeederRegistry{
		seeders: make([]Seeder, 0),
	}
}

func (r *SeederRegistry) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

func (r *SeederRegistry) RunAll(db *gorm.DB) error {
	for _, seeder := range r.seeders {
		err := seeder.Seed(db)
		log.Printf("Running seeder: %s\n", seeder.GetName())
		if err != nil {
			return fmt.Errorf("error running seeder %s: %w", seeder.GetName(), err)
		}
		log.Printf("✓ Seeder %s completed\n", seeder.GetName())

	}

	return nil
}

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

func (r *SeederRegistry) GetAllSeeders() []string {
	names := make([]string, len(r.seeders))
	for i, seeder := range r.seeders {
		names[i] = seeder.GetName()
	}
	return names
}

func InitSeeders() *SeederRegistry {
	registry := NewSeederRegistry()

	registry.Register(&UserSeeder{})

	return registry
}
