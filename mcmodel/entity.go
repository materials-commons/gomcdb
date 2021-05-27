package mcmodel

type Entity struct {
	ID           int
	Name         string
	Files        []File `gorm:"many2many:entity2file"`
	EntityStates []EntityState
}

type EntityState struct {
	ID       int
	EntityID int
}
