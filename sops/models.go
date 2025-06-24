package sops

type KmsConf struct {
	ARN     string
	Profile string
}

func (c KmsConf) IsConfigured() bool {
	return len(c.ARN) > 0
}

type PgpConf struct {
	Fingerprint string
}

func (c PgpConf) IsConfigured() bool {
	return len(c.Fingerprint) > 0
}
