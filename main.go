package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

/*
	Create Database
	Read from DB
*/

// Mock data structure for Trains
type Train struct {
	TrainNumber   string             `json:"train_number" bson:"train_number"`
	TrainName     string             `json:"train_name" bson:"train_name"`
	FromStation   string             `json:"from_station" bson:"from_station"`
	ToStation     string             `json:"to_station" bson:"to_station"`
	DepartureTime string             `json:"departure_time" bson:"departure_time"`
	ArrivalTime   string             `json:"arrival_time" bson:"arrival_time"`
	TravelTime    string             `json:"travel_time" bson:"travel_time"`
	Classes       []string           `json:"classes" bson:"classes"`
	Availability  map[string]string  `json:"availability" bson:"availability"` // class => status
	Fare          map[string]float64 `json:"fare" bson:"fare"`                 // class => fare
}


var mysqlDB *sql.DB

/*
	var mockTrains = []Train{
		{
			TrainNumber:   "12345",
			TrainName:     "Rajdhani Express",
			FromStation:   "KSR",
			ToStation:     "BBS",
			DepartureTime: "16:00",
			ArrivalTime:   "08:00",
			TravelTime:    "16h",
			Classes:       []string{"1A", "2A", "3A", "SL"},
			Availability:  map[string]string{"1A": "Available", "2A": "Waiting 2", "3A": "Available", "SL": "RAC"},
			Fare:          map[string]float64{"1A": 3200, "2A": 2100, "3A": 1400, "SL": 500},
		},
		{
			TrainNumber:   "67890",
			TrainName:     "Duronto Express",
			FromStation:   "NDLS",
			ToStation:     "BCT",
			DepartureTime: "22:00",
			ArrivalTime:   "14:00",
			TravelTime:    "16h",
			Classes:       []string{"1A", "2A", "3A", "SL"},
			Availability:  map[string]string{"1A": "Waiting 5", "2A": "Available", "3A": "RAC", "SL": "Available"},
			Fare:          map[string]float64{"1A": 3500, "2A": 2200, "3A": 1500, "SL": 600},
		},
	}
*/
func initMySQL() {
	var err error
	mysqlDB, err = sql.Open("mysql", "root:2010@tcp(localhost:3306)/irctc")
	if err != nil {
		log.Fatal("MySQL connection failed:", err)
	}
}

func fetchFromMySQL(from string, to string) []Train {
	// Fetch data from MySQL
	availableTrains := []Train{}
	query := "select a.train_number, train_name, from_station, to_station, departure_time, arrival_time, travel_time, GROUP_CONCAT(b.class_code) AS Classes,JSON_OBJECTAGG(b.class_code, b.availability), JSON_OBJECTAGG(b.class_code, b.fare) from trains a left join train_classes b on a.train_number = b.train_number where a.from_station LIKE ? AND a.to_station LIKE ? GROUP BY a.train_number"

	rows, err := mysqlDB.Query(query, "%"+from+"%", "%"+to+"%")
	if err != nil {
		log.Fatal("MySQL query failed:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var classCSV string
		var availabilityJSON, fareJSON []byte
		var train Train
		//	fmt.Println(rows)
		err := rows.Scan(&train.TrainNumber, &train.TrainName, &train.FromStation, &train.ToStation, &train.DepartureTime, &train.ArrivalTime, &train.TravelTime,
			&classCSV, &availabilityJSON, &fareJSON)
		if err != nil {
			log.Fatal("MySQL scan failed:", err)
		}

		train.Classes = strings.Split(classCSV, ",")

		if err := json.Unmarshal(availabilityJSON, &train.Availability); err != nil {
			log.Fatal("JSON unmarshal failed:", err)
		}

		if err := json.Unmarshal(fareJSON, &train.Fare); err != nil {
			log.Fatal("JSON unmarshal failed:", err)
		}

		availableTrains = append(availableTrains, train)
	}

	return availableTrains
}

func searchTrains(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	journeyDate := c.Query("date")
	//class := c.Query("class")

	if from == "" || to == "" || journeyDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide from, to, and date parameters"})
		return
	}

	availableTrains := []Train{}
	initMySQL()

	availableTrains = fetchFromMySQL(from, to)
	/*	for _, train := range mockTrains {
			if train.FromStation == from && train.ToStation == to {
				availableTrains = append(availableTrains, train)
			}
		}
	*/
	if len(availableTrains) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No trains found for the given route"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":   journeyDate,
		"from":   from,
		"to":     to,
		"trains": availableTrains,
	})
}

func main() {
	r := gin.Default()

	r.GET("/api/train-search", searchTrains)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
