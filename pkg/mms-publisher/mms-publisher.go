package mmspublisher

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/metno/go-mms/pkg/mms"
)

func PublishProduct(productEvent mms.ProductEvent) error {

	// hardcoded to test-server. Should be findable from ProductionHub?
	url := productEvent.ProductionHub + "/api/v1/events"

	// Create a json-payload from productEvent
	jsonStr, err := json.Marshal(&productEvent)
	// Create a http-request to post the payload
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	// Hardcoded Api-Key, maybe in productEvent?

	mmsApiKey, ok := os.LookupEnv("MMS_API_KEY")
	if !ok {
		log.Fatal("MMS_API_KEY not found in the environment, can't publish to MMS")
	}

	httpReq.Header.Set("Api-Key", mmsApiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	// Create a http connection to the api.
	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Failed to create http client: %v", err)
	}
	defer httpResp.Body.Close()

	// If 201 is not returned, panic with http response
	if httpResp.StatusCode != http.StatusCreated {
		log.Fatalf("Product event not posted: %s", httpResp.Status)
	}
	return nil

}
