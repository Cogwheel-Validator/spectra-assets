package fetch

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func NewCosmosFetcher(chainId string, APIs []string) CosmosFetcher {
	return CosmosFetcher{
		ChainId: chainId,
		APIs:    APIs,
	}
}

func testApi(chainId, api string) (healthy bool) {
	url := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/node_info", api)
	resp, err := http.Get(url)

	if err != nil {
		return false
	}

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

	if nodeInfo.DefaultNodeInfo.Network != chainId {
		return false
	}

	return true

}
