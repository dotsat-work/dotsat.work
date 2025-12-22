package main

import (
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"dotsat.work/internal/model"
	"dotsat.work/internal/repository"
)

func main() {
	// Connect to database
	db, err := sqlx.Connect("pgx", "postgres://postgres:postgres@localhost:5432/dotsat?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	repo := repository.NewTenantRepository(db)

	// Create test tenants
	tenants := []*model.Tenant{
		{
			ID:        uuid.New(),
			Name:      "Hewlett Packard Enterprise",
			Subdomain: "hpe",
			Status:    "active",
			Tier:      "premium",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Honeywell International",
			Subdomain: "honeywell",
			Status:    "active",
			Tier:      "standard",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Rockwell Automation",
			Subdomain: "rockwell",
			Status:    "active",
			Tier:      "premium",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Emerson Electric",
			Subdomain: "emerson",
			Status:    "active",
			Tier:      "standard",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Name:      "Schneider Electric",
			Subdomain: "schneider",
			Status:    "trial",
			Tier:      "standard",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, tenant := range tenants {
		err := repo.Create(tenant)
		if err != nil {
			log.Printf("failed to create tenant %s: %v", tenant.Name, err)
			continue
		}
		log.Printf("✅ Created tenant: %s (%s)", tenant.Name, tenant.Subdomain)
	}

	log.Println("\n🎉 Test data seeded successfully!")
}
