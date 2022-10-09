package routes

import (
	//"net/http"
	"fmt"
	"math/rand"
	"time"

	//"github.com/labstack/echo/v5"
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func addEventCode(e *core.RecordCreateEvent) error {
	if e.Record.TableName() == "events" {
		// generate the id for the event
		alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		id := ""
		
		// create a random seed each run to prevent the same event IDs being generated over and over
		rand.Seed(time.Now().Unix())

		// generate a random character 6 times to make an ID
		for i := 0; i < 6; i++ {
			c := alphabet[rand.Intn(len(alphabet))]
			id += string(c)
		}

		// debug statement
		fmt.Println(id)

		// set the event_id to the generated one
		e.Record.SetDataValue("event_id", id)
	}

	return nil
}
