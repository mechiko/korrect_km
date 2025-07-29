package trueclient

import "time"

type filterSearch struct {
	ProductGroups []string `json:"productGroups"`
	Gtins         []string `json:"gtins,omitempty"`
	OrderIds      []string `json:"orderIds,omitempty"`
}

type paginationSearch struct {
	PerPage          int    `json:"perPage,omitempty"`
	LastEmissionDate string `json:"lastEmissionDate"`
	Sgtin            string `json:"sgtin"`
	Direction        int    `json:"direction,omitempty"`
}

type searchQueryFilter struct {
	Filter     *filterSearch     `json:"filter"`
	Pagination *paginationSearch `json:"pagination,omitempty"`
}

type SearchResult struct {
	Result []struct {
		Sgtin              string    `json:"sgtin"`
		Cis                string    `json:"cis"`
		CisWithoutBrackets string    `json:"cisWithoutBrackets"`
		Gtin               string    `json:"gtin"`
		ProducerInn        string    `json:"producerInn"`
		Status             string    `json:"status"`
		EmissionDate       time.Time `json:"emissionDate"`
		ApplicationDate    time.Time `json:"applicationDate,omitempty"`
		ProductionDate     time.Time `json:"productionDate,omitempty"`
		GeneralPackageType string    `json:"generalPackageType"`
		OwnerInn           string    `json:"ownerInn"`
		EmissionType       string    `json:"emissionType"`
		ProductGroup       string    `json:"productGroup"`
		HaveChildren       bool      `json:"haveChildren"`
		ExpirationDate     time.Time `json:"expirationDate,omitempty"`
		IntroducedDate     time.Time `json:"introducedDate,omitempty"`
		ModID              string    `json:"modId,omitempty"`
		Expiration         []struct {
			ExpirationStorageDate time.Time `json:"expirationStorageDate"`
			StorageConditionID    string    `json:"storageConditionId"`
		} `json:"expiration"`
		OrderID                 string    `json:"orderId,omitempty"`
		OperationIntroducedDate time.Time `json:"operationIntroducedDate,omitempty"`
		Parent                  string    `json:"parent,omitempty"`
	} `json:"result"`
	IsLastPage bool `json:"isLastPage"`
}
