package glauth

import "fmt"

type Context struct {
	Username string
	Password string
	Hostname string
	Port     string
	Database string
}

func (c *Context) Dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Username, c.Password, c.Hostname, c.Port, c.Database)
}
