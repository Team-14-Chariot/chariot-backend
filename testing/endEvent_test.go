package testing

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	. "github.com/Team-14-Chariot/chariot-backend/routes"
	"github.com/pocketbase/pocketbase/tests"
)

func TestEndEventEndpoint(t *testing.T) {
	var goodBodyBuf bytes.Buffer
	goodBody := EndEventBody{
		Event_id: "USA",
	}
	json.NewEncoder(&goodBodyBuf).Encode(goodBody)

	setupTestApp := func() (*tests.TestApp, error) {
		app, err := tests.NewTestApp("../cmd/chariot-backend/pb_data")
		if err != nil {
			return nil, err
		}

		Routes(app)

		return app, nil
	}

	scenarios := []tests.ApiScenario{
		{
			Name:            "Validate Correct Operation",
			Method:          http.MethodPost,
			Url:             "/api/endEvent",
			Body:            &goodBodyBuf,
			ExpectedEvents: []{},
			ExpectedStatus:  200,
			ExpectedContent: nil,
			TestAppFactory:  setupTestApp,
		},
	}

	for _, scenario := range scenarios {
		scenario.Test(t)
	}
}
