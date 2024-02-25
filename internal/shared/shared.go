package shared

import "time"

type Route struct {
	Domain   string `yaml:"domain" json:"domain" validate:"required,fqdn"`
	Upstream string `yaml:"upstream" json:"upstream" validate:"required,url,startswith=http://"`
}

type Certificate struct {
	CertFilePath string    `json:"certFilePath" validate:"required,filepath"`
	KeyFilePath  string    `json:"keyFilePath" validate:"required,filepath"`
	ExpiresAt    time.Time `json:"expiresAt" validate:"required"`
}

type DomainCertificates map[string]Certificate
