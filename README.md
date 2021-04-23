### About
This library gives convenient access to MySQL.  
#####Setup to run test application
Create following table DB name = "testdb"
```sql
CREATE TABLE `cities` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `population` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8;
```

##### Config file
You can use this library where you can create Config object or read it from yaml file 
```yaml
databases:
  configs:
    mysql_master:
      type: mysql
      user: root
      password: root
      url: [ localhost ]
      port: 3306
      db: testdb
```
##### How to use 
```go
package main

import (
	"fmt"
	"github.com/harishb2k/gox-base"
	"github.com/harishb2k/gox-base/serialization"
	goxdb "github.com/harishb2k/gox-database"
	_ "github.com/harishb2k/gox-mysql"
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
	if result, err := db.Persist("metrics", "INSERT INTO cities(name, population) VALUES(?, ?)", "india", 11); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		fmt.Printf("Result: %v \n", result)
	}

	// 4 - Find all records
	if data, err := db.Find("a", "SELECT * FROM cities WHERE ID >= ?", 1); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		for k, v := range data {
			fmt.Println(k, v)
		}
	}

	// 5 - Find a record
	if data, err := db.FindOne("a", "SELECT * FROM cities WHERE ID = ?", 1); err != nil {
		fmt.Printf("%v\n", err)
		return
	} else {
		fmt.Println(data)
	}
}

```