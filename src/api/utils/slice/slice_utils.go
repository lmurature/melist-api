package slice

import "github.com/lmurature/melist-api/src/api/domain/share"

func ShareConfigUserExists(configs share.ShareConfigs, userId int64) bool {
	for _, c := range configs {
		if c.UserId == userId {
			return true
		}
	}
	return false
}
