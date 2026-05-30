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

// GatherChainValidators pages through all on-chain validators and returns their
// descriptions keyed by operator address.
func (f CosmosFetcher) GatherChainValidators() map[string]ValidatorDescription {
	result := make(map[string]ValidatorDescription)
	nextKey := ""
	for {
		api := f.APIs[rand.Intn(len(f.APIs))]
		page := fetchValoperList(api, nextKey)
		if page == nil {
			break
		}
		for _, v := range page.Validators {
			result[v.OperatorAddress] = v.Description
		}
		nextKey = page.Pagination.NextKey
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
