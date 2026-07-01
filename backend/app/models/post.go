package models

import (
	"time"

	"github.com/goravel/framework/database/orm"
	"gorm.io/datatypes"
)

type Post struct {
	orm.Model
	Slug        string            `gorm:"uniqueIndex" json:"slug"`
	TitleEn     string            `json:"title_en"`
	TitleId     string            `json:"title_id"`
	ExcerptEn   *string           `json:"excerpt_en"`
	ExcerptId   *string           `json:"excerpt_id"`
	ContentEn   *string           `json:"content_en"`
	ContentId   *string           `json:"content_id"`
	Category    *string           `json:"category"`
	Tags        datatypes.JSONMap `gorm:"type:json" json:"tags"`
	Thumbnail   *string           `json:"thumbnail"`
	Status      string            `json:"status"`
	PublishedAt *time.Time        `json:"published_at"`
	ReadTime    int               `json:"read_time"`
}
