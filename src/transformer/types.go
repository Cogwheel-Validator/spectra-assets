package transformer

type ValidatorMeta struct {
	Moniker         string `json:"moniker,omitempty"`
	Identity        string `json:"identity,omitempty"`
	Description     string `json:"description,omitempty"`
	Website         string `json:"website,omitempty"`
	SecurityContact string `json:"security_contact,omitempty"`
}

type ValidatorEntry struct {
	Address  string
	DataDir  string
	Meta     ValidatorMeta
	LogoPath string // empty when no logo file is present
}

type ValidatorAsset struct {
	Address         string `json:"address"`
	Moniker         string `json:"moniker,omitempty"`
	Identity        string `json:"identity,omitempty"`
	Description     string `json:"description,omitempty"`
	Website         string `json:"website,omitempty"`
	SecurityContact string `json:"security_contact,omitempty"`
}
