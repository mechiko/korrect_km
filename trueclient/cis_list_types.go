package trueclient

import "time"

type CisesValues struct {
	Values []string `json:"values"`
}

type Cises struct {
	ApplicationDate   time.Time `json:"applicationDate"`
	IntroducedDate    time.Time `json:"introducedDate"`
	VolumeSpecialPack string    `json:"volumeSpecialPack"`
	IsVarQuantity     bool      `json:"isVarQuantity"`
	Expirations       []struct {
		ExpirationStorageDate time.Time `json:"expirationStorageDate"`
		StorageConditionID    int       `json:"storageConditionId"`
		StorageConditionName  string    `json:"storageConditionName"`
	} `json:"expirations"`
	Kpp      string `json:"kpp"`
	OwnerMod struct {
		ModID   int    `json:"modId"`
		Kpp     string `json:"kpp"`
		Address string `json:"address"`
	} `json:"ownerMod"`
	Ogvs               []interface{} `json:"ogvs"`
	RequestedCis       string        `json:"requestedCis"`
	Cis                string        `json:"cis"`
	CisWithoutBrackets string        `json:"cisWithoutBrackets"`
	Status             string        `json:"status"`
	StatusEx           string        `json:"statusEx"`
	Gtin               string        `json:"gtin"`
	ProductName        string        `json:"productName"`
	ProductGroup       string        `json:"productGroup"`
	ProductGroupID     int           `json:"productGroupId"`
	ProducedDate       time.Time     `json:"producedDate"`
	PackageType        string        `json:"packageType"`
	GeneralPackageType string        `json:"generalPackageType"`
	ProducerInn        string        `json:"producerInn"`
	ProducerName       string        `json:"producerName"`
	EmissionDate       time.Time     `json:"emissionDate"`
	EmissionType       string        `json:"emissionType"`
	OwnerInn           string        `json:"ownerInn"`
	OwnerName          string        `json:"ownerName"`
	TnVedEaes          string        `json:"tnVedEaes"`
	TnVedEaesGroup     string        `json:"tnVedEaesGroup"`
	Child              []interface{} `json:"child"`
	MaxRetailPrice     interface{}   `json:"maxRetailPrice"`
}
