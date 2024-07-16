package ressources

type Group struct {
	ID        int    // internal id number, used by glauth
	Name      string // LDAP group name (i.e. cn or ou depending on context)
	GIDNumber int    // LDAP GID attribute
}

type CreateGroup struct {
	Name      string
	GIDNumber int
}

type UpdateGroup struct {
	Name      *string
	GIDNumber *int
}
