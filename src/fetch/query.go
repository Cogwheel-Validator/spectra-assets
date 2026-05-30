package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func fetchValoperList(api, nextKey string) *Validators {
	query, _ := url.Parse(fmt.Sprintf("%s/cosmos/staking/v1beta1/validators", api))
	params := url.Values{}
	if nextKey != "" {
		params.Set("pagination.key", nextKey)
	}
	params.Set("pagination.limit", "100")
	query.RawQuery = params.Encode()

	resp, err := http.Get(query.String())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var result Validators
	if err := json.Unmarshal(body, &result); err != nil {
		return nil
	}

	return &result
}
