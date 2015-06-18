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
	ID		bson.ObjectId `bson:"_id,omitempty"`
	
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
	    "object_type": "Sample"}

	q := c.Find(ServiceQuery)
	count, _ := q.Count()
	fmt.Printf("Total Services Found: ", count)
	if count > 0 {
		result := AnalysisResults{}
		iter := q.Iter()
		for iter.Next(&result) {
			fmt.Printf("Result:" , result.ID)
		}
		if iter.Err != nil {
				log.Fatal(iter.Err)
		}
	}

	fmt.Println("Finished")

}