package ressources

type CapabilityAction string

const (
	CapabilityActionSearch CapabilityAction = "search"

	/**
	all capabilities are not implemented yet in glauth : https://glauth.github.io/docs/capabilities.html
	*/
)

func (c CapabilityAction) String() string {
	switch c {
	case CapabilityActionSearch:
		return "search"
	default:
		return "unknown"
	}
}

type Capability struct {
	ID     int              // internal id number, used by glauth
	UserID int              // internal user id number, used by glauth
	Action CapabilityAction // string representing an allowed action, e.g. “search”
	Object string           // string representing scope of allowed action, e.g. “ou=superheros,dc=glauth,dc=com”
}
