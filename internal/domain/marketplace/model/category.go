package model

type Category struct {
	ID       int        `gorm:"primaryKey"`
	Name     string     `gorm:"not null"`
	ImageURL string     `gorm:"not null"`
	ParentID *int       `gorm:"foreignKey:CategoryID"`
	Children []Category `gorm:"foreignKey:ParentID"`
	Level    int        `gorm:"-"`
}
