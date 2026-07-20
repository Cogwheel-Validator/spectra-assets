package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func NewCosmosFetcher(chainId string, APIs []string) (CosmosFetcher, error) {
	var healthyApis []string
	for _, api := range APIs {
		if testApi(chainId, api) {
			healthyApis = append(healthyApis, api)
		}
	}
	if len(healthyApis) == 0 {
		return CosmosFetcher{}, fmt.Errorf("no healthy APIs found for chain %s", chainId)
	}
	return CosmosFetcher{ChainId: chainId, APIs: healthyApis}, nil
}

// GatherChainOperators pages through all on-chain validators and returns their
// descriptions keyed by operator address.
//
// It supports both validators and governors, depending on the `governors` flag.
func (f CosmosFetcher) GatherChainOperators(governors bool) map[string]OperatorDescription {
	fetchPage := fetchValoperPage
	if governors {
		fetchPage = fetchGovernorPage
	}

	result := make(map[string]OperatorDescription)
	retryCount := 0
	maxRetries := 10
	nextKey := ""
	for {
		api := f.APIs[rand.Intn(len(f.APIs))]
		entries, next, err := fetchPage(api, nextKey)
		if err != nil {
			retryCount++
			if retryCount >= maxRetries {
				break
			}
			continue
		}

		for _, e := range entries {
			result[e.Address] = e.Description
		}
		nextKey = next
		if nextKey == "" {
			break
		}
	}
	return result
}

func testApi(chainId, api string) (healthy bool) {
	url := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/node_info", api)
	resp, err := http.Get(url)

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	var nodeInfo NodeInfo
	if err := json.Unmarshal(body, &nodeInfo); err != nil {
		return false
	}

	return nodeInfo.DefaultNodeInfo.Network == chainId
}

type operatorEntry struct {
	Address     string
	Description OperatorDescription
}

func fetchValoperPage(api, nextKey string) ([]operatorEntry, string, error) {
	page, err := fetchValoperList(api, nextKey)
	if err != nil {
		return nil, "", err
	}
	entries := make([]operatorEntry, len(page.Validators))
	for i, v := range page.Validators {
		entries[i] = operatorEntry{Address: v.OperatorAddress, Description: v.Description}
	}
	return entries, page.Pagination.NextKey, nil
}

func fetchGovernorPage(api, nextKey string) ([]operatorEntry, string, error) {
	page, err := fetchGovernorList(api, nextKey)
	if err != nil {
		return nil, "", err
	}
	entries := make([]operatorEntry, len(page.Governors))
	for i, g := range page.Governors {
		entries[i] = operatorEntry{Address: g.GovernorAddress, Description: g.Description}
	}
	return entries, page.Pagination.NextKey, nil
}
