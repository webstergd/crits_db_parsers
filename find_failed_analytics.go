package main 

import (
		"fmt"
		"flag"
		"log"
		"gopkg.in/mgo.v2"
		"gopkg.in/mgo.v2/bson"
)

var ServerName string
var DatabaseName string
var ServiceName string
func init() {
	const (
		defaultServer 			= "mongodb.example.com"
		defaultServerHelp		= "MongoDB Server"
		defaultDatabase 		= "crits"
		defaultDatabaseHelp		= "Database name for CRITs"
		defaultService			= "cuckoo_w_api"
		defaultServiceHelp		= "Which Service to clean up (separate with , for multiple"		
	)
	flag.StringVar(&ServerName, "server", defaultServer, defaultServerHelp)
	flag.StringVar(&DatabaseName, "database", defaultDatabase, defaultDatabaseHelp)
	flag.StringVar(&ServiceName, "service", defaultService, defaultServiceHelp)
}

type AnalysisResults struct {
	ID				bson.ObjectId	`bson:"_id,omitempty"`
	status			string 			`bson:"status"`
	service_name	string 			`bson:"service_name"`
}

func main() {
	flag.Parse()

	fmt.Println("Connecting to Server ", ServerName)
	session, err := mgo.Dial(ServerName)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// set mote to monotonic
	session.SetMode(mgo.Monotonic, true)

	// create session for analysis_results
	c := session.DB(DatabaseName).C("analysis_results")

	ServiceQuery := bson.M{
	    "service_name": bson.M{ "$in": []string{ServiceName} },
	    "status": bson.M{ "$in": []string{"started", "error"} },
	    "object_type": "Sample",
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
			fmt.Println("Processing: %v | %v | %v" , result.ID, result.service_name, result.status)
			//c.RemoveId(result.ID)
		}
		if err := iter.Close(); err != nil {
				log.Fatal(iter.Err)
		}
	}

	fmt.Println("Finished")

}