package game

import (
	"github.com/bitly/go-simplejson"
	"github.com/wonderivan/logger"
	"sync"
)

var globalsTaxPercent struct {
	mu                 sync.RWMutex
	platformTaxPercent *simplejson.Json
}

func SetPlatformTaxPercent(allTax *simplejson.Json) {
	globalsTaxPercent.mu.Lock()
	defer globalsTaxPercent.mu.Unlock()
	globalsTaxPercent.platformTaxPercent = allTax

}

func GetPlatformTaxPercent(pkgId int) float64 {

	globalsTaxPercent.mu.Lock()
	allTax := globalsTaxPercent.platformTaxPercent
	globalsTaxPercent.mu.Unlock()

	var tax int
	for i := 0; ; i++ {
		index := allTax.GetIndex(i)
		if index.Interface() == nil {
			break
		}

		if index.Get("package_id").MustInt() == pkgId {
			tax = index.Get("platform_tax_percent").MustInt()
			break
		}
	}

	if tax == 0 {
		logger.Debug("没有对应id tax,以0.05计算")
		return 0.05
	}

	return float64(tax) * 0.01
}
