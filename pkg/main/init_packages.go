package initpck

import (
	"github.com/trucktrace/pkg/hsh"
	redisService "github.com/trucktrace/pkg/redis"
	"github.com/trucktrace/pkg/sort"
	"github.com/trucktrace/pkg/validation"
)

type Packages struct {
	Hsh        *hsh.HashService
	Redis      *redisService.RedisService
	Validation *validation.Validation
	Sort       *sort.Sort
}

var Redis = new(redisService.RedisService)

func GetPackages() *Packages {
	return &Packages{
		Redis:      Redis,
		Hsh:        hsh.GetHashService(),
		Validation: validation.GetValidation(),
		Sort:       sort.GetSortPackage(),
	}
}
