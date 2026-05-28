package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var keybaseApi = "https://keybase.io/_/api/1.0/user/lookup.json"

func QueryKeybase(identity string) (*KeybaseResponse, error) {
	url := fmt.Sprintf("%s?key_suffix=%s&fields=pictures", keybaseApi, identity)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result KeybaseResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
