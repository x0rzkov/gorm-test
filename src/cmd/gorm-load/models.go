package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var testTables = []interface{}{
	&Star{},
	&Tag{},
}

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
}

type User struct {
	Model
	Email string `validator:"email" gorm:"unique_index" json:"email"`
}

type Category struct {
	Model
	Name string `json:"name"`
}

type Topic struct {
	Model
	CategoryID uint
	Category   Category
	Title      string `json:"title"`
	Acl        []User `json:"acl" gorm:"many2many:topic_acl"`
}

// Star represents a starred repository
type Star struct {
	// gorm.Model
	RemoteID    int       `json:"id" gorm:"index;"`
	Name        *string   `gorm:"type:varchar(255);" json:"name"`
	FullName    *string   `gorm:"type:varchar(255);" json:"full_name"`
	Description *string   `gorm:"type:longtext;" json:"description"`
	Homepage    *string   `gorm:"type:varchar(255);" json:"homepage"`
	URL         *string   `gorm:"type:varchar(255);" json:"svn_url"`
	Language    *string   `gorm:"type:varchar(64);" json:"language"`
	Topics      []string  `gorm:"-" json:"topics"`
	Stargazers  int       `json:"stargarzers_count"`
	StarredAt   time.Time `json:"starred_at"`
	ServiceID   uint      `json:"service_id" gorm:"index;"`
	Tags        []Tag     `json:"-" gorm:"many2many:star_tags;"`
}

/*
func (s *Star) setTagsOld(db *gorm.DB, tags []string) error {
	var topicList []Tag
	for _, tag := range tags {

		var tagModel Tag
		err := db.FirstOrCreate(&tagModel, Tag{Name: tag}).Error
		if err != nil {
			return err
		}
		topicList = append(topicList, tagModel)
	}
	s.Topics = topicList
	return nil
}
*/

func (s *Star) setTags(db *gorm.DB, topics []string, prefix string, noEmpty bool) error {
	topics = RemoveSliceDuplicates(topics, true)
	tags := make([]Tag, 0)
	for _, topic := range topics {
		name := fmt.Sprintf("%s%s", prefix, topic)
		if len(strings.TrimSpace(name)) == 0 {
			continue
		}
		tag := Tag{Name: fmt.Sprintf("%s%s", prefix, topic)}
		var tagModel Tag
		err := db.FirstOrCreate(&tagModel, tag).Error
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	s.Tags = tags
	return nil
}

type Tag struct {
	// gorm.Model
	// ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique_index;type:varchar(255);not null"`
	// Name string `gorm:"primary_key;type:varchar(255);not null"`
	// Name string `gorm:"primary_key;type:varchar(255);not null;"`
	// Name      string `gorm:"type:varchar(255);unique_index;not null;"`
	Stars     []Star `gorm:"many2many:star_tags;"`
	StarCount int    `gorm:"-"`
}

func truncateTables(db *gorm.DB, tables ...interface{}) {
	for _, table := range tables {
		if err := db.DropTableIfExists(table).Error; err != nil {
			log.Fatalln("error while dropping table: ", err)
		}
		if err := db.DropTableIfExists(table).Error; err != nil {
			log.Fatalln("error while dropping table: ", err)
		}
		if err := db.AutoMigrate(table).Error; err != nil {
			log.Fatalln("error while auto-migrating table: ", err)
		}
	}
}
