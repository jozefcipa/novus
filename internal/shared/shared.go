package shared

import "time"

type Route struct {
	Domain   string `yaml:"domain" json:"domain"` // TODO: validate config to ensure http is always present
	Upstream string `yaml:"upstream" json:"upstream"`
}

type Certificate struct {
	CertFilePath string    `json:"certFilePath"`
	KeyFilePath  string    `json:"keyFilePath"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type DomainCertificates map[string]Certificate
