package mcmodel

type Entity struct {
	ID    int
	Name  string
	Files []File `gorm:"many2many:entity2file"`
}
