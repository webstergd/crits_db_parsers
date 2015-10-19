package main 

import (
		"fmt"
		"flag"
		"log"
		"gopkg.in/mgo.v2"
		"gopkg.in/mgo.v2/bson"
		"strings"
)

var ServerName string
var DatabaseName string
var ServiceName string
var ServiceStatus string

// type to store analysis results from db
type AnalysisResults struct {
	ID				bson.ObjectId	`bson:"_id,omitempty"`
	Status			string 			`bson:"status"`
	Service_name	string 			`bson:"service_name"`
}

func init() {
	const (
		defaultServer 			= "mongodb.example.com"
		defaultServerHelp		= "MongoDB Server"
		defaultDatabase 		= "crits"
		defaultDatabaseHelp		= "Database name for CRITs"
		defaultService			= "cuckoo_w_api"
		defaultServiceHelp		= "Which service to clean up (separate with ','' for multiple)?"		
		defaultStatus			= "started,error"
		defaultStatusHelp		= "What is the status flag should be used to remove? (separate with ','' for multiple)"		
	)
	flag.StringVar(&ServerName, "server", defaultServer, defaultServerHelp)
	flag.StringVar(&DatabaseName, "database", defaultDatabase, defaultDatabaseHelp)
	flag.StringVar(&ServiceName, "service", defaultService, defaultServiceHelp)
	flag.StringVar(&ServiceStatus, "status", defaultStatus, defaultStatusHelp)
}

func main() {
	flag.Parse()

	fmt.Println("Connecting to Server: ", ServerName)
	session, err := mgo.Dial(ServerName)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// set mote to monotonic
	session.SetMode(mgo.Monotonic, true)

	// create session for analysis_results
	fmt.Println("Connecting to Database: ", DatabaseName)
	c := session.DB(DatabaseName).C("analysis_results")

	services := strings.Split(ServiceName, ",")
	status := strings.Split(ServiceStatus, ",")
	fmt.Println("Searching for services: %v with status: %v", services, status)
	ServiceQuery := bson.M{
		"service_name": bson.M{ "$in": services },
		"status": bson.M{ "$in": status },
	}
	ServiceQueryRestrict := bson.M{
		"_id": 1,
		"status": 1,
		"service_name": 1,
	}
	q := c.Find(ServiceQuery).Select(ServiceQueryRestrict)
	count, _ := q.Count()
	fmt.Println("Total Services Found: ", count)
	if count > 0 {
		result := AnalysisResults{}
		iter := q.Iter()
		for iter.Next(&result) {
			fmt.Println("Processing: %v | %v | %v" , result.ID, result.Service_name, result.Status)
			c.RemoveId(result.ID)
		}
		if err := iter.Close(); err != nil {
				log.Fatal(iter.Err)
		}
	}

	fmt.Println("Finished cleaning %v services", count)

}