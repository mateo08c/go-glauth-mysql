package glauth

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Glauth struct {
	context *Context
	db      *gorm.DB
}

func New(c *Context) (*Glauth, error) {
	g := &Glauth{
		context: c,
	}

	err := g.connect()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Glauth) connect() error {
	mo := mysql.Open(g.context.Dsn())
	db, err := gorm.Open(mo, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	g.db = db

	return nil
}
