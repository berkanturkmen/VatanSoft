package database

import (
	"log"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/berkanturkmen/VatanSoft/models"
	"github.com/joho/godotenv"
)

var DB * gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

		dsn := "host=" + host + " user=" + username + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=disable TimeZone=Europe/Istanbul"

		DB,
	err = gorm.Open(postgres.Open(dsn), & gorm.Config {})
	if err != nil {
		log.Fatal(err)
	}

	Migrate()
	SeedDB()
}

func Migrate() {
	err := DB.AutoMigrate( & models.City {}, & models.Region {}, & models.Hospital {}, & models.Employee {}, & models.Session {}, & models.Polyclinic {}, & models.HospitalPolyclinic {}, & models.Job {}, & models.Title {}, & models.Personnel {}, )

	if err != nil {
		log.Fatal(err)
	}
}

func SeedDB() {
	cities := [] models.City {
		{
			Name: "İstanbul",
		}, {
			Name: "Ankara",
		}, {
			Name: "İzmir",
		},
	}

		for _,
	city := range cities {
		result := DB.FirstOrCreate( & city, models.City {
			Name: city.Name,
		})
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
	}

		regions := [] models.Region {
		{
			Name: "Pendik",
			CityID: 1,
		}, {
			Name: "Çankaya",
			CityID: 2,
		}, {
			Name: "Bornova",
			CityID: 3,
		},
	}

		for _,
	region := range regions {
		result := DB.FirstOrCreate( & region, models.Region {
			Name: region.Name,
			CityID: region.CityID,
		})
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
	}

		polyclinics := [] models.Polyclinic {
		{
			Name: "Üroloji",
		}, {
			Name: "Nöroloji",
		}, {
			Name: "İç Hastalıkları",
		},
	}

		for _,
	polyclinic := range polyclinics {
		result := DB.FirstOrCreate( & polyclinic, models.Polyclinic {
			Name: polyclinic.Name,
		})
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
	}

		jobs := [] models.Job {
		{
			Name: "Doktor",
		}, {
			Name: "İdari Personel",
		}, {
			Name: "Hizmet Personeli",
		},
	}

		for _,
	job := range jobs {
		result := DB.FirstOrCreate( & job, models.Job {
			Name: job.Name,
		})
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
	}

		titles := [] models.Title {
		{
			Name: "Asistan",
			JobID: 1,
		}, {
			Name: "Başhekim",
			JobID: 2,
		}, {
			Name: "Müdür",
			JobID: 2,
		}, {
			Name: "Danışman",
			JobID: 3,
		}, {
			Name: "Temizlikçi",
			JobID: 3,
		}, {
			Name: "Güvenlik",
			JobID: 3,
		},
	}

		for _,
	title := range titles {
		result := DB.FirstOrCreate( & title, models.Title {
			Name: title.Name,
			JobID: title.JobID,
		})
		if result.Error != nil {
			log.Fatalf("Error: %v", result.Error)
		}
	}
}
