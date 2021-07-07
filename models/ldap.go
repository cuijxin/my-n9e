package models

type LdapSection struct {
	Enable          bool           `yaml:"enable"`
	Host            string         `yaml:"host"`
	Port            int            `yaml:"port"`
	BaseDn          string         `yaml:"baseDn"`
	BindUser        string         `yaml:"bindUser"`
	BindPass        string         `yaml:"bindPass"`
	AuthFilter      string         `yaml:"authFilter"`
	Attributes      ldapAttributes `yaml:"attributes"`
	CoverAttributes bool           `yaml:"coverAttributes"`
	TLS             bool           `yaml:"tls"`
	StartTLS        bool           `yaml:"startTLS"`
}

type ldapAttributes struct {
	Nickname string `yaml:"nickname"`
	Phone    string `yaml:"phone"`
	Email    string `yaml:"email"`
}

var LDAP LdapSection

func InitLdap(ldap LdapSection) {
	LDAP = ldap
}
