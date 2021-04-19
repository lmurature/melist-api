package date_utils

import (
	"github.com/lmurature/melist-api/src/api/config"
	"time"
)

func GetNowDateFormatted() string {
	return time.Now().UTC().Format(config.DbDateLayout)
}