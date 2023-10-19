package testlib

import (
	"encoding/json"
)

type NuoDBCertInfo struct {
	CAPathLength     int    `json:"caPathLength"`
	CertificatePem   string `json:"certificatePem"`
	ExpiresTimestamp string `json:"expiresTimestamp"`
	IssuerName       string `json:"issuerName"`
	SubjectName      string `json:"subjectName"`
}

type NuoDBCertificateInfo struct {
	ServerCertificates  map[string]NuoDBCertInfo `json:"serverCertificates"`
	ProcessCertificates map[string]NuoDBCertInfo `json:"processCertificates"`
	TrustedCertificates map[string]NuoDBCertInfo `json:"trustedCertificates"`
}

func UnmarshalCertificateInfo(s string) (err error, info NuoDBCertificateInfo) {
	err = json.Unmarshal([]byte(s), &info)
	return
}
