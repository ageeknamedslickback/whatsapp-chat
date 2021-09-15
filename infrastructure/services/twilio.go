package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ageeknamedslickback/whatsapp-chat/domain"
)

// MakePostRequest makes a post request to Twilio's conversatinal API
func MakePostRequest(payload url.Values, target domain.Message) error {
	url := fmt.Sprintf(
		"https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json",
		os.Getenv("TWILIO_ACCOUNT_SID"),
	)
	req, err := http.NewRequest(
		http.MethodPost,
		url,
		strings.NewReader(payload.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to make a new request: %w", err)
	}
	req.SetBasicAuth(
		os.Getenv("TWILIO_ACCOUNT_SID"),
		os.Getenv("TWILIO_AUTH_TOKEN"),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"status code %v returned with data %s",
			resp.StatusCode,
			string(b),
		)
	}

	if err := json.Unmarshal(b, &target); err != nil {
		return fmt.Errorf(
			"failed to unmarshall response to target message: %w",
			err,
		)
	}

	return nil
}
