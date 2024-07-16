package glauth

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mateo08c/go-glauth-mysql/glauth/models"
	"github.com/mateo08c/go-glauth-mysql/glauth/ressources"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (g *Glauth) UserExistByName(name string) (bool, error) {
	var user models.User
	err := g.db.Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g *Glauth) UserExistByUID(uid int) (bool, error) {
	var user models.User
	err := g.db.Where("uidnumber = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func (g *Glauth) CreateCapability(c *ressources.Capability) error {
	ca := &models.Capability{
		UserID: c.UserID,
		Action: string(c.Action),
		Object: c.Object,
	}

	err := g.db.Create(ca).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *Glauth) FindNextUserID() (int, error) {
	var user models.User
	err := g.db.Order("uidnumber desc").First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 20000, nil
		}
		return 0, err
	}

	return user.UIDNumber + 1, nil
}

func (g *Glauth) userModelToResource(u *models.User) (*ressources.User, error) {
	// Création de la ressource User avec les champs de base
	r := &ressources.User{
		ID:            u.ID,
		Name:          u.Name,
		UIDNumber:     u.UIDNumber,
		GivenName:     u.GivenName,
		SN:            u.SN,
		Mail:          u.Mail,
		LoginShell:    u.LoginShell,
		HomeDirectory: u.HomeDirectory,
		Disabled:      u.Disabled,
		PassSHA256:    u.PassSHA256,
		PassBCrypt:    u.PassBCrypt,
		OTPSecret:     u.OTPSecret,
		Yubikey:       u.Yubikey,
		SSHKeys:       u.SSHKeys,
		CustAttr:      u.CustAttr,
	}

	// Ajout du primary group s'il est défini dans le modèle User
	if u.PrimaryGroup != 0 {
		pg, err := g.GetGroupByGID(u.PrimaryGroup)
		if err == nil {
			if pg != nil {
				r.PrimaryGroup = &ressources.Group{
					ID:        pg.ID,
					Name:      pg.Name,
					GIDNumber: pg.GIDNumber,
				}
			}
		}
	}

	// Ajout des other groups s'ils sont définis dans le modèle User
	if u.OtherGroups != nil && len(u.OtherGroups) > 0 {
		ogIDs := strings.Split(strings.Trim(string(u.OtherGroups), ","), ",")
		for _, gID := range ogIDs {
			if i, err := strconv.Atoi(gID); err == nil {
				og, err := g.GetGroupByGID(i)
				if err != nil {
					continue
				}
				if og != nil {
					if r.PrimaryGroup != nil && r.PrimaryGroup.GIDNumber == og.GIDNumber {
						continue
					}

					if !GroupExistsInList(r.OtherGroups, og) {
						r.OtherGroups = append(r.OtherGroups, &ressources.Group{
							ID:        og.ID,
							Name:      og.Name,
							GIDNumber: og.GIDNumber,
						})
					}
				}
			}
		}
	}

	// Ajout des include groups en utilisant la fonction existante
	includeGroups, err := g.GetIncludeGroupsByIncludeGroupGID(u.PrimaryGroup)
	if err == nil {
		for _, ig := range includeGroups {
			if r.PrimaryGroup != nil && r.PrimaryGroup.GIDNumber == ig.GIDNumber {
				continue
			}

			if !GroupExistsInList(r.OtherGroups, ig) {
				r.OtherGroups = append(r.OtherGroups, ig)
			}
		}
	}

	// Ajout des capabilities en utilisant la fonction existante
	capabilities, err := g.GetCapabilitiesByUserUIDNumber(u.UIDNumber)
	if err == nil {
		if capabilities != nil {
			r.Capabilities = capabilities
		}
	}

	return r, nil
}

func (g *Glauth) GetCapabilitiesByUserUIDNumber(uid int) ([]*ressources.Capability, error) {
	var capabilities []*models.Capability
	err := g.db.Where("userid = ?", uid).Table("capabilities").Find(&capabilities).Error
	if err != nil {
		return nil, err
	}

	var resCaps []*ressources.Capability
	for _, c := range capabilities {
		resCaps = append(resCaps, &ressources.Capability{
			ID:     c.ID,
			UserID: c.UserID,
			Action: ressources.CapabilityAction(c.Action),
			Object: c.Object,
		})
	}

	return resCaps, nil
}

func (g *Glauth) GetUserByName(s string) (*ressources.User, error) {
	var user models.User
	err := g.db.Where("name = ?", s).First(&user).Error
	if err != nil {
		return nil, err
	}

	return g.userModelToResource(&user)
}

func (g *Glauth) UpdateUserPassword(name, password string) error {
	h := sha256.New()
	h.Write([]byte(password))
	pass := fmt.Sprintf("%x", h.Sum(nil))

	err := g.db.Table("users").Where("name = ?", name).Update("passsha256", pass).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *Glauth) UpdateUserPasswordByUID(uid int, password string) error {
	h := sha256.New()
	h.Write([]byte(password))
	pass := fmt.Sprintf("%x", h.Sum(nil))

	err := g.db.Table("users").Where("uidnumber = ?", uid).Update("passsha256", pass).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *Glauth) UpdateUser(name string, u *ressources.UpdateUser) error {
	var user models.User
	err := g.db.Table("users").Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Update fields from the UpdateUser struct if they are not nil
	if u.GivenName != nil {
		user.GivenName = *u.GivenName
	}
	if u.SN != nil {
		user.SN = *u.SN
	}
	if u.Mail != nil {
		user.Mail = *u.Mail
	}
	if u.LoginShell != nil {
		user.LoginShell = *u.LoginShell
	}
	if u.HomeDirectory != nil {
		user.HomeDirectory = *u.HomeDirectory
	}
	if u.Disabled != nil {
		user.Disabled = *u.Disabled
	}
	if u.OTPSecret != nil {
		user.OTPSecret = *u.OTPSecret
	}
	if u.Yubikey != nil {
		user.Yubikey = *u.Yubikey
	}
	if u.SSHKeys != nil {
		user.SSHKeys = *u.SSHKeys
	}

	// Update the OtherGroups field if it is provided
	if u.OtherGroups != nil {
		//check if the groups exist
		for _, gID := range *u.OtherGroups {
			exist, err := g.GroupExistByGID(gID)
			if err != nil {
				return err
			}

			if !exist {
				return errors.New("group with GID " + strconv.Itoa(gID) + " does not exist")
			}
		}

		//check if the primary group is in the other groups
		if user.PrimaryGroup != 0 {
			for _, gID := range *u.OtherGroups {
				if user.PrimaryGroup == gID {
					return errors.New("primary group cannot be in the other groups")
				}
			}
		}

		user.OtherGroups = []byte(ToCommaSeparatedString(*u.OtherGroups))
	}

	// Update the password if provided
	if u.Password != nil {
		h := sha256.New()
		h.Write([]byte(*u.Password))
		user.PassSHA256 = fmt.Sprintf("%x", h.Sum(nil))
	}

	// Validate and update Custom Attributes
	if u.CustAttr != nil {
		err := json.Unmarshal([]byte(*u.CustAttr), &map[string]interface{}{})
		if err != nil {
			return err
		}
		user.CustAttr = *u.CustAttr
	}

	// Update the user in the database
	err = g.db.Save(&user).Error
	if err != nil {
		return err
	}

	// Update capabilities if provided
	if u.Capabilities != nil {
		// Optionally clear existing capabilities or handle updates accordingly
		if err := g.db.Where("user_id = ?", user.UIDNumber).Delete(&models.Capability{}).Error; err != nil {
			return err
		}

		for _, c := range *u.Capabilities {
			c.UserID = user.UIDNumber
			err = g.CreateCapability(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Glauth) GetUsers() ([]*ressources.User, error) {
	var users []*models.User
	err := g.db.Find(&users).Error
	if err != nil {
		return nil, err
	}

	var resUsers []*ressources.User
	for _, u := range users {
		var r *ressources.User
		r, err = g.userModelToResource(u)
		if err != nil {
			continue
		}

		if r == nil {
			continue
		}

		resUsers = append(resUsers, r)
	}

	return resUsers, nil
}

func (g *Glauth) GetUserByUID(uid int) (*ressources.User, error) {
	var user *models.User
	err := g.db.Where("uidnumber = ?", uid).First(&user).Error
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return g.userModelToResource(user)
}

func (g *Glauth) CreateUser(u *ressources.CreateUser) error {
	user := &models.User{
		Name:          u.Name,
		OtherGroups:   []byte(ToCommaSeparatedString(u.OtherGroups)),
		GivenName:     u.GivenName,
		SN:            u.SN,
		Mail:          u.Mail,
		LoginShell:    u.LoginShell,
		HomeDirectory: u.HomeDirectory,
		Disabled:      u.Disabled,
		OTPSecret:     u.OTPSecret,
		Yubikey:       u.Yubikey,
		SSHKeys:       u.SSHKeys,
		CustAttr:      u.CustAttr,
	}

	//check if user already exists
	exists, err := g.UserExistByName(u.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("user already exists")
	}

	if u.UIDNumber == 0 {
		id, err := g.FindNextUserID()
		if err != nil {
			return err
		}

		exists, err = g.UserExistByUID(id)
		if err != nil {
			return err
		}

		user.UIDNumber = id
	}

	if u.PrimaryGroup != 0 {
		user.PrimaryGroup = u.PrimaryGroup
	}

	if u.Password != "" {
		h := sha256.New()
		h.Write([]byte(u.Password))
		user.PassSHA256 = fmt.Sprintf("%x", h.Sum(nil))
	}

	if u.CustAttr != "" {
		err := json.Unmarshal([]byte(u.CustAttr), &map[string]interface{}{})
		if err != nil {
			return err
		}

		user.CustAttr = u.CustAttr
	} else {
		user.CustAttr = "{}"
	}

	err = g.db.Create(user).Error
	if err != nil {
		return err
	}

	if u.Capabilities != nil {
		for _, c := range u.Capabilities {
			c.UserID = user.UIDNumber
			err = g.CreateCapability(c)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Glauth) DeleteUser(uid int) error {
	var user models.User
	err := g.db.Table("users").Where("uidnumber = ?", uid).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	err = g.db.Table("users").Delete(&user).Error
	if err != nil {
		return err
	}

	//delete capabilities
	err = g.db.Table("capabilities").Where("userid = ?", user.UIDNumber).Delete(&models.Capability{}).Error
	if err != nil {
		return err
	}

	return nil
}
