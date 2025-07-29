package trueclient

type Mod struct {
	Kpp          string `json:"kpp"`
	ProductGroup string `json:"productGroup"`
	IsBlocked    bool   `json:"isBlocked"`
	Address      string `json:"address"`
	FiasId       string `json:"fiasId"`
}

type ModSlice []Mod

type ProductGroupInfo struct {
	ProductGroup string `json:"productGroup"`
	Status       string `json:"status"`
}
type ProductGroupInfoSlice []ProductGroupInfo

type ListMods struct {
	Result   ModSlice `json:"result"`
	Total    int      `json:"total"`
	NextPage bool     `json:"nextPage"`
}

type ModInfo struct {
	ProductGroupInfo ProductGroupInfoSlice `json:"productGroupInfo"`
	Mods             ModSlice              `json:"mods"`
}

type ModsListPostParam struct {
	Pg  []string `json:"pg"`
	Inn string   `json:"inn"`
	// Kpp []string `json:"kpp"`
}
