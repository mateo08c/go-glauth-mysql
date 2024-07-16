package models

type Capability struct {
	ID     int    // internal id number, used by glauth
	UserID int    `gorm:"column:userid"` // internal user id number, used by glauth
	Action string `gorm:"column:action"` // string representing an allowed action, e.g. “search”
	Object string `gorm:"column:object"` // string representing scope of allowed action, e.g. “ou=superheros,dc=glauth,dc=com”
}
