package models

// User Struct for users table
type User struct {
	ID            int     // internal id number, used by glauth
	Name          string  `gorm:"column:name"`          // LDAP name (i.e. cn, uid)
	UIDNumber     int     `gorm:"column:uidnumber"`     // LDAP UID attribute
	PrimaryGroup  int     `gorm:"column:primarygroup"`  // An LDAP group’s GID attribute; also used to build ou attribute; used to build memberOf
	OtherGroups   []uint8 `gorm:"column:othergroups"`   // A comma-separated list of GID attributes; used to build memberOf
	GivenName     string  `gorm:"column:givenname"`     // LDAP GivenName attribute, i.e. an account’s first name
	SN            string  `gorm:"column:sn"`            // LDAP sn attribute, i.e. an account’s last name
	Mail          string  `gorm:"column:mail"`          // LDAP mail attribute, i.e. email address; also used as userPrincipalName
	LoginShell    string  `gorm:"column:loginshell"`    // LDAP loginShell attribute, pushed to the client, may be ignored
	HomeDirectory string  `gorm:"column:homedirectory"` // LDAP homeDirectory attribute, pushed to the client, may be ignored
	Disabled      bool    `gorm:"column:disabled"`      // LDAP accountStatus attribute, if non-zero returns “inactive”
	PassSHA256    string  `gorm:"column:passsha256"`    // SHA256 account password
	PassBCrypt    string  `gorm:"column:passbcrypt"`    // BCRYPT-encrypted account password
	OTPSecret     string  `gorm:"column:otpsecret"`     // OTP secret, for two-factor authentication
	Yubikey       string  `gorm:"column:yubikey"`       // UBIKey, for two-factor authentication
	SSHKeys       string  `gorm:"column:sshkeys"`       // A comma-separated list of sshPublicKey attributes
	CustAttr      string  `gorm:"column:custattr"`      // A JSON-encoded string, containing arbitrary additional attributes; must be {} by default
}
