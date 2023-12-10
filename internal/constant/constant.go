package constant

import "time"

// length of short url alias
const AliasLength int = 10

// db tables
const (
	URLSTable string = "urls"
)

// available databases
const (
	POSTGRES string = "postgres"
	REDIS    string = "redis"
)

const ZeroTTL time.Duration = 0
