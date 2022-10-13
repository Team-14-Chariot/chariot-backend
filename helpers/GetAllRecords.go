package helpers

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

func GetAllRecords(app *pocketbase.PocketBase, collection *models.Collection) []models.Record {
	ids := []models.Record{}
	app.Dao().RecordQuery(collection).All(&ids)

	records := []models.Record{}
	for _, id := range ids {
		record, _ := app.Dao().FindFirstRecordByData(collection, "id", id.GetId())
		records = append(records, *record)
	}

	return records
}
