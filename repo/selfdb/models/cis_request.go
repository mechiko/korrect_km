package models

import (
	"github.com/upper/db/v4"
)

type CisRequest struct {
	Cis      string `db:"cis"`
	Status   string `db:"status"`
	StatusEx string `db:"status_ex"`
	Responce string `db:"responce"`
	Produced string `db:"produced"`
	Expired  string `db:"expired"`
}

type CisRequestSlice []*CisRequest

func (p *CisRequest) Store(sess db.Session) db.Store {
	return sess.Collection("cis_request")
}
