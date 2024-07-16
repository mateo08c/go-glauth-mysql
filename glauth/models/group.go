package models

// LDAPGroup Struct for ldapgroups table
type LDAPGroup struct {
	ID        int    // internal id number, used by glauth
	Name      string `gorm:"column:name"`      // LDAP group name (i.e. cn or ou depending on context)
	GIDNumber int    `gorm:"column:gidnumber"` // LDAP GID attribute
}

// IncludeGroup Struct for includegroups table
type IncludeGroup struct {
	ID             int // internal id number, used by glauth
	ParentGroupID  int `gorm:"column:parentgroupid"`  // the LDAP group id to be included in, used by glauth
	IncludeGroupID int `gorm:"column:includegroupid"` // the LDAP group id to be included in the parent group, used by glauth
}
