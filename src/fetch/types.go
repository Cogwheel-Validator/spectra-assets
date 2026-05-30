package fetch

// Keybase Response types //

type KeybaseResponse struct {
	Status KeybaseStatus `json:"status"`
	Them   []KeybaseThem `json:"them"`
}

type KeybaseStatus struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

type KeybaseThem struct {
	Id       string       `json:"id"`
	Pictures ThemPictures `json:"pictures"`
}

type ThemPictures struct {
	Primary ThemPrimary `json:"primary"`
}
type ThemPrimary struct {
	Url    string `json:"url"`
	Source any    `json:"source"`
}

// Cosmos API Response types //

type Validators struct {
	Validators []Validator `json:"validators"`
	Pagination Pagination  `json:"pagination"`
}

type Validator struct {
	OperatorAddress string               `json:"operator_address"`
	Description     ValidatorDescription `json:"description"`
}
type ValidatorDescription struct {
	Moniker  string `json:"moniker"`
	Details  string `json:"details"`
	Identity string `json:"identity"`
}

type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

type NodeInfo struct {
	DefaultNodeInfo struct {
		Network string `json:"network"`
	} `json:"default_node_info"`
}

// Cosmos fetch type //

type CosmosFetcher struct {
	ChainId string
	APIs    []string
}

type CosmosValInfo struct {
	OperatorAddress string
	Identity        string
	ImageUrl        string
}
