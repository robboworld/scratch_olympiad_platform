package consts

type Mode string

const (
	Production  Mode = "production"
	Development Mode = "development"
)

const (
	KeyId   = "keyId"
	KeyRole = "keyRole"
)

const (
	MaxSolutionFileSize = 100 * 1024 * 1024 // 100 MB
)
