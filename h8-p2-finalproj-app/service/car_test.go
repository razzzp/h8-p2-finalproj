package service_test

import (
	"encoding/json"
	"fmt"
	"h8-p2-finalproj-app/model"
	"h8-p2-finalproj-app/service"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateTestDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME")+"_test",
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Car{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Rental{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.Payment{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.AutoMigrate(&model.TopUp{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func RunFindOnQuery(db *gorm.DB) ([]map[string]any, error) {
	var results []map[string]any
	err := db.Find(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func PrintAsJSON(m any) {
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(string(jsonData))
	}
}

func TestCountRentalsPerCar(t *testing.T) {
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}
	carService := service.NewCarService(db)
	q := carService.CountRentalsPerCarQuery(&service.GetCarsQueryParams{})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 5)
}

func TestCountRentalsPerCarWithDatesLargerThanRental(t *testing.T) {
	// r is rental item
	// q is query
	//  r.StartDate               r.EndDate
	// 	----|-----------|--------|---|----
	//            q.StartDate q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	q := carService.CountRentalsPerCarQuery(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 1)
}

func TestCountRentalsPerCarWithStartDateOverlap(t *testing.T) {
	// r is rental item
	// q is query
	//               r.StartDate     r.EndDate
	// 	----|-----------|--------|---|----
	//    q.StartDate       q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-08-30")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	q := carService.CountRentalsPerCarQuery(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 1)
}

func TestCountRentalsPerCarWithEndDateOverlap(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-02")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	q := carService.CountRentalsPerCarQuery(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 1)
}

func TestCountRentalsPerCarNotFound(t *testing.T) {
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-08-02")
	endDate, _ := time.Parse(time.DateOnly, "2024-08-04")
	carService := service.NewCarService(db)
	q := carService.CountRentalsPerCarQuery(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 0)
}

func TestCarsWithRentals(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-02")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	q := carService.CarsWithRentalQuery(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	PrintAsJSON(res)
	assert.Len(t, res, 5)
}

func TestCarsWithRentalsWithSeats(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	seats := 5
	carService := service.NewCarService(db)
	q := carService.CarsWithRentalQuery(&service.GetCarsQueryParams{
		Seats: &seats,
	})
	res, err := RunFindOnQuery(q)
	if err != nil {
		t.Error(err)
	}
	PrintAsJSON(res)
	assert.Len(t, res, 3)
}

func TestGetCarsWithRentals(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-02")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	res, err := carService.GetCarsWithRentals(&service.GetCarsQueryParams{
		StartDate: &startDate,
		EndDate:   &endDate,
	})
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.Len(t, res, 5)
	assert.Equal(t, uint(1), res[0].NumOfRentals)
}

func TestIsCarAvailableTrue(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-02")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-04")
	carService := service.NewCarService(db)
	res, err := carService.IsCarAvailable(
		1,
		startDate,
		endDate,
	)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.True(t, res)
}

func TestIsCarAvailableFalse(t *testing.T) {
	// r is rental item
	// q is query
	//   r.StartDate   r.EndDate
	// 	----|----|------|--------|--------
	//         q.StartDate      q.EndDate
	godotenv.Load("../.env")
	db := CreateTestDB()
	if db == nil {
		t.FailNow()
	}

	startDate, _ := time.Parse(time.DateOnly, "2024-09-21")
	endDate, _ := time.Parse(time.DateOnly, "2024-09-23")
	carService := service.NewCarService(db)
	res, err := carService.IsCarAvailable(
		5,
		startDate,
		endDate,
	)
	if err != nil {
		t.Error(err)
	}
	// PrintAsJSON(res)
	assert.False(t, res)
}
