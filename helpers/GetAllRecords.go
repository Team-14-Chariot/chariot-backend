package helpers

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

func GetAllRecords(app core.App, collection *models.Collection) []models.Record {
	ids := []models.Record{}
	app.Dao().RecordQuery(collection).All(&ids)

	records := []models.Record{}
	for _, id := range ids {
		record, _ := app.Dao().FindFirstRecordByData(collection, "id", id.GetId())
		records = append(records, *record)
	}

	return records
}
