package service

import (
	"encoding/xml"

	"github.com/bk-rim/openbanking/model"

	"time"
)

type GrpHdrType struct {
	MsgId   string `xml:"MsgId"`
	CreDtTm string `xml:"CreDtTm"`
}

type AmtType struct {
	Ccy string  `xml:"Ccy,attr"`
	Val float64 `xml:",chardata"`
}

type AccType struct {
	Nm       string `xml:"Nm"`
	CdtrAcct struct {
		Id struct {
			IBAN string `xml:"IBAN"`
		} `xml:"Id"`
	} `xml:"CdtrAcct"`
}

type CdtrAcctType struct {
	Id struct {
		IBAN string `xml:"IBAN"`
	} `xml:"Id"`
}

func generateXML(payment model.Payment) ([]byte, error) {
	xmlData := &struct {
		XMLName xml.Name   `xml:"Document"`
		GrpHdr  GrpHdrType `xml:"GrpHdr"`
		Cdtr    struct {
			Nm       string `xml:"Nm"`
			CdtrAcct struct {
				Id struct {
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			} `xml:"CdtrAcct"`
		} `xml:"Cdtr"`
		Dbtr struct {
			Nm       string `xml:"Nm"`
			CdtrAcct struct {
				Id struct {
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			} `xml:"CdtrAcct"`
		} `xml:"Dbtr"`
		Amt struct {
			Ccy string  `xml:"Ccy,attr"`
			Val float64 `xml:",chardata"`
		} `xml:"Amt"`
	}{
		GrpHdr: GrpHdrType{
			MsgId:   payment.IdempotencyUniqueKey,
			CreDtTm: time.Now().UTC().Format(time.RFC3339),
		},
		Cdtr: AccType{
			Nm: payment.CreditorName,
			CdtrAcct: struct {
				Id struct {
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			}{
				Id: struct {
					IBAN string `xml:"IBAN"`
				}{
					IBAN: payment.CreditorIBAN,
				},
			},
		},
		Dbtr: AccType{
			Nm: payment.DebtorName,
			CdtrAcct: struct {
				Id struct {
					IBAN string `xml:"IBAN"`
				} `xml:"Id"`
			}{
				Id: struct {
					IBAN string `xml:"IBAN"`
				}{
					IBAN: payment.DebtorIBAN,
				},
			},
		},
		Amt: AmtType{
			Ccy: "EUR",
			Val: payment.Amount,
		},
	}

	xmlBytes, err := xml.MarshalIndent(xmlData, "", "    ")
	if err != nil {
		return nil, err
	}

	xmlWithHeader := []byte(xml.Header + string(xmlBytes))

	return xmlWithHeader, nil
}
