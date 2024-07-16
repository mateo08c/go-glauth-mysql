package ressources

type User struct {
	ID            int           // internal id number, used by glauth
	Name          string        // LDAP name (i.e. cn, uid)
	UIDNumber     int           // LDAP UID attribute
	PrimaryGroup  *Group        // An LDAP group’s GID attribute; also used to build ou attribute; used to build memberOf
	OtherGroups   []*Group      // A comma-separated list of GID attributes; used to build memberOf
	Capabilities  []*Capability // A list of capabilities
	GivenName     string        // LDAP GivenName attribute, i.e. an account’s first name
	SN            string        // LDAP sn attribute, i.e. an account’s last name
	Mail          string        // LDAP mail attribute, i.e. email address; also used as userPrincipalName
	LoginShell    string        // LDAP loginShell attribute, pushed to the client, may be ignored
	HomeDirectory string        // LDAP homeDirectory attribute, pushed to the client, may be ignored
	Disabled      bool          // LDAP accountStatus attribute, if non-zero returns “inactive”
	PassSHA256    string        // SHA256 account password
	PassBCrypt    string        // BCRYPT-encrypted account password
	OTPSecret     string        // OTP secret, for two-factor authentication
	Yubikey       string        // UBIKey, for two-factor authentication
	SSHKeys       string        // A comma-separated list of sshPublicKey attributes
	CustAttr      string        // A JSON-encoded string, containing arbitrary additional attributes; must be {} by default
}

type CreateUser struct {
	ID            int
	Name          string
	UIDNumber     int
	PrimaryGroup  int
	OtherGroups   []int
	Capabilities  []*Capability
	GivenName     string
	SN            string
	Mail          string
	LoginShell    string
	HomeDirectory string
	Disabled      bool
	Password      string
	OTPSecret     string
	Yubikey       string
	SSHKeys       string
	CustAttr      string
}

type UpdateUser struct {
	ID            int
	UIDNumber     *int
	PrimaryGroup  *int
	OtherGroups   *[]int
	Capabilities  *[]*Capability
	GivenName     *string
	SN            *string
	Mail          *string
	LoginShell    *string
	HomeDirectory *string
	Disabled      *bool
	Password      *string
	OTPSecret     *string
	Yubikey       *string
	SSHKeys       *string
	CustAttr      *string
}
