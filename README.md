# Zero Allocation JSON Logger for GORM

This package supports the use of [Zerolog](https://github.com/rs/zerolog) with [Gorm](https://gorm.io/)

## Usage

```go
package main

import (
	"time"

	"github.com/truongkma/gormzerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger := gormzerolog.NewLogger(gormzerolog.Config{
		SlowThreshold:        time.Second,
		ParameterizedQueries: true,
	})

	dsn := "host=localhost user=postgres password=postgres dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	database, err := gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{
			Logger: logger,
		},
	)
}

```
