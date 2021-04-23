package main

import (
	"fmt"
	"github.com/harishb2k/gox-base"
	"github.com/harishb2k/gox-base/metrics"
	"github.com/harishb2k/gox-base/serialization"
	goxdb "github.com/harishb2k/gox-database"
	_ "github.com/harishb2k/gox-mysql"
	"sync/atomic"
	"time"
)

type mysqlTestAppParsingConfig struct {
	Databases goxdb.Configs `yaml:"databases"`
}

func main() {
	// 1 - Read config from file
	// You can chose to create your own config object from code
	config := mysqlTestAppParsingConfig{}
	if err := serialization.ReadYaml("./testapp/config.yaml", &config); err != nil {
		fmt.Printf("Error %v", err)
		return
	}
	config.Databases.SetupDefaults()

	// 2 - Create a new DB instance. Here we will create DB instance for "mysql_master"
	var db goxdb.Db
	var err error
	cfg := config.Databases.Configs["mysql_master"]
	if db, err = goxdb.GetOrCreate("mysql_master", &cfg, gox.NewNoOpCrossFunction()); err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// 3 - Persist a record
	if result, err := db.Persist(metrics.LabeledMetric{Name: "test"}, "INSERT INTO cities(name, population) VALUES(?, ?)", "india", 11); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		fmt.Printf("Result: %v \n", result)
	}

	// 4 - Find all records
	if data, err := db.Find(metrics.LabeledMetric{Name: "test"}, "SELECT * FROM cities WHERE ID >= ?", 1); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		for k, v := range data {
			fmt.Println(k, v)
		}
	}

	// 5 - Find a record
	if data, err := db.FindOne(metrics.LabeledMetric{Name: "test"}, "SELECT * FROM cities WHERE ID = ?", 1); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		fmt.Println(data)
	}

	fmt.Println("--- Start ---")
	perf(db)
}

func perf(db goxdb.Db) {
	var count int32 = 0
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100000000; j++ {
				atomic.AddInt32(&count, 1)
				if _, err := db.Persist(metrics.LabeledMetric{Name: "test"}, "INSERT INTO cities(name, population) VALUES(?, ?)", "india", 11); err != nil {
					fmt.Printf("Error  %v \n", err)
				} else {
					if count%1000 == 0 {
						fmt.Println(count)
					}
				}
			}
		}()
	}
	time.Sleep(1 * time.Minute)
}
