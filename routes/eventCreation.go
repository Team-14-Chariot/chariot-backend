package routes

import (
	//"net/http"
	"fmt"
	"math/rand"
	//"github.com/labstack/echo/v5"
	//"github.com/pocketbase/pocketbase"
	//"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func addEventCode(e *core.RecordCreateEvent) error {
	// generate the id for the event
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id := ""

	// generate a random character 6 times to make an ID
	for i := 0; i < 6; i++ {
		c := alphabet[rand.Intn(len(alphabet))]
		id += string(c)
	}

	fmt.Println(id)

	// get the record being inserted into the table
	e.Record.SetDataValue("event_id", id)

	return nil
}
