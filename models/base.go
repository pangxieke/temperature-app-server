package models

import (
	"github.com/jinzhu/gorm"
)

type Model interface {
	Id() uint
}

func applyOffsetAndLimit(s *gorm.DB, limitAndOffset ...int) *gorm.DB {
	if len(limitAndOffset) >= 1 {
		if limitAndOffset[0] != 0 {
			s = s.Limit(limitAndOffset[0])
		}

	}
	if len(limitAndOffset) >= 2 {
		if limitAndOffset[1] != 0 {
			s = s.Offset(limitAndOffset[1])
		}
	}
	return s
}
