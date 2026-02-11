package dictionary

import "gorm.io/gorm"

type Dictionary struct {
	gorm.Model
	Key     string           `gorm:"unique;comment:字典键"`
	Comment string           `gorm:"comment:字典注释"`
	Items   []DictionaryItem `gorm:"foreignKey:DictionaryID;references:ID"`
}

func (d *Dictionary) TableName() string {
	return "dictionary"
}

type DictionaryItem struct {
	gorm.Model
	DictionaryID uint   `gorm:"comment:字典ID"`
	Key          string `gorm:"comment:字典项键"`
	Comment      string `gorm:"comment:字典项注释"`
	Order        int    `gorm:"comment:字典项排序"`
}

func (d *DictionaryItem) TableName() string {
	return "dictionary_item"
}

func InitDictionaryTable(db *gorm.DB) {
	db.AutoMigrate(&Dictionary{})
	db.AutoMigrate(&DictionaryItem{})

	var count int64
	db.Model(&Dictionary{}).Count(&count)
	if count > 0 {
		return
	}
}
