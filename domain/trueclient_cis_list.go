package domain

import (
	"time"
)

type Cis struct {
	Cis      string
	Status   string
	StatusEx string
}

type CisSlice []*Cis

type CisesValues struct {
	Values []string `json:"values"`
}

type CisJson struct {
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

type CisJsonSlice []*CisJson

type CisPostJson struct {
	Result struct {
		RequestedCis        string        `json:"requestedCis"`
		Cis                 string        `json:"cis"`
		CisWithoutBrackets  string        `json:"cisWithoutBrackets"`
		Gtin                string        `json:"gtin"`
		ProducerInn         string        `json:"producerInn"`
		Status              string        `json:"status"`
		StatusEx            string        `json:"statusEx"`
		EmissionDate        time.Time     `json:"emissionDate"`
		ApplicationDate     time.Time     `json:"applicationDate"`
		Parent              string        `json:"parent"`
		PackageType         string        `json:"packageType"`
		OwnerInn            string        `json:"ownerInn"`
		TnVedEaesGroup      string        `json:"tnVedEaesGroup"`
		TnVedEaes           string        `json:"tnVedEaes"`
		ReceiptDate         time.Time     `json:"receiptDate"`
		EmissionType        string        `json:"emissionType"`
		ConnectDate         time.Time     `json:"connectDate"`
		ProductGroup        string        `json:"productGroup"`
		ProductGroupID      int           `json:"productGroupId"`
		ExtendedPackageType string        `json:"extendedPackageType"`
		ProducedDate        time.Time     `json:"producedDate"`
		VolumeSpecialPack   string        `json:"volumeSpecialPack"`
		Children            []interface{} `json:"children"`
		OwnerMod            struct {
			ModID int `json:"modId"`
		} `json:"ownerMod"`
		Ogvs        []interface{} `json:"ogvs"`
		Expirations []struct {
			StorageConditionID    int    `json:"storageConditionId"`
			ExpirationStorageDate string `json:"expirationStorageDate"`
		} `json:"expirations"`
		Licences []struct {
			LicenceNumber string `json:"licenceNumber"`
			LicenceDate   string `json:"licenceDate"`
		} `json:"licences"`
		SpecialAttributes struct {
			MaxRetailPrice           float64   `json:"maxRetailPrice"`
			ExpirationDate           time.Time `json:"expirationDate"`
			PrVetDocument            string    `json:"prVetDocument"`
			Capacity                 string    `json:"capacity"`
			TurnoverType             string    `json:"turnoverType"`
			RetType                  string    `json:"retType"`
			ExpNum                   string    `json:"expNum"`
			ExpName                  string    `json:"expName"`
			RemainsImport            string    `json:"remainsImport"`
			FtsDecisionCode          string    `json:"ftsDecisionCode"`
			QuantityInPack           int       `json:"quantityInPack"`
			SoldCount                int       `json:"soldCount"`
			EliminationReasonOther   string    `json:"eliminationReasonOther"`
			ProductWeightGr          int       `json:"productWeightGr"`
			ManufacturerSerialNumber string    `json:"manufacturerSerialNumber"`
			IntroducedDate           time.Time `json:"introducedDate"`
			StatusEx                 string    `json:"statusEx"`
			ProductGroup             string    `json:"productGroup"`
			ProductGroupID           int       `json:"productGroupId"`
			ExtendedPackageType      string    `json:"extendedPackageType"`
			WithdrawReason           string    `json:"withdrawReason"`
			NextCis                  []string  `json:"nextCis"`
			PrevCis                  []string  `json:"prevCis"`
			ApprovementDocument      struct {
				CertDoc []struct {
					Type       string    `json:"type"`
					Number     string    `json:"number"`
					Date       time.Time `json:"date"`
					WellNumber string    `json:"wellNumber"`
				} `json:"certDoc"`
				DeclarationDate      string `json:"declarationDate"`
				DeclarationRegNumber string `json:"declarationRegNumber"`
				DeclarationID        string `json:"declarationId"`
			} `json:"approvementDocument"`
			IsVarQuantity bool `json:"isVarQuantity"`
		} `json:"specialAttributes"`
	} `json:"result"`
	ErrorMessage string `json:"errorMessage"`
	ErrorCode    string `json:"errorCode"`
}
