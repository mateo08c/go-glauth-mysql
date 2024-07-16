package glauth

import (
	"errors"
	"github.com/mateo08c/go-glauth-mysql/glauth/models"
	"github.com/mateo08c/go-glauth-mysql/glauth/ressources"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func (g *Glauth) GetGroupByGID(gid int) (*ressources.Group, error) {
	var group models.LDAPGroup
	err := g.db.Where("gidnumber = ?", gid).Table("ldapgroups").First(&group).Error
	if err != nil {
		return nil, err
	}

	return &ressources.Group{
		ID:        group.ID,
		Name:      group.Name,
		GIDNumber: group.GIDNumber,
	}, nil
}

func (g *Glauth) GetGroupByName(name string) (*ressources.Group, error) {
	var group models.LDAPGroup
	err := g.db.Where("name = ?", name).Table("ldapgroups").First(&group).Error
	if err != nil {
		return nil, err
	}

	return &ressources.Group{
		ID:        group.ID,
		Name:      group.Name,
		GIDNumber: group.GIDNumber,
	}, nil
}

func (g *Glauth) GetIncludeGroupsByIncludeGroupGID(gid int) ([]*ressources.Group, error) {
	var includeGroups []*models.IncludeGroup
	err := g.db.Where("includegroupid = ?", gid).Table("includegroups").Find(&includeGroups).Error
	if err != nil {
		return nil, err
	}

	var resGroups []*ressources.Group
	for _, ig := range includeGroups {
		var gr *ressources.Group
		gr, err = g.GetGroupByGID(ig.ParentGroupID)
		if err != nil {
			continue
		}

		if gr == nil {
			continue
		}

		resGroups = append(resGroups, gr)
	}

	return resGroups, nil
}

func (g *Glauth) CreateGroup(gr *ressources.CreateGroup) error {
	exists, err := g.GroupExistByGID(gr.GIDNumber)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("group with GID " + strconv.Itoa(gr.GIDNumber) + " already exists")
	}

	exists, err = g.GroupExistByName(gr.Name)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("group with name " + gr.Name + " already exists")
	}

	group := &models.LDAPGroup{
		Name:      gr.Name,
		GIDNumber: gr.GIDNumber,
	}

	if gr.GIDNumber == 0 {
		id, err := g.FindNextGroupID()
		if err != nil {
			return err
		}

		exists, err = g.GroupExistByGID(id)
		if err != nil {
			return err
		}

		group.GIDNumber = id
	}

	err = g.db.Table("ldapgroups").Create(group).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *Glauth) UpdateGroup(name string, gr *ressources.UpdateGroup) error {
	var group models.LDAPGroup
	err := g.db.Table("ldapgroups").Where("name = ?", name).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("group not found")
		}
		return err
	}

	if gr.Name != nil {
		group.Name = *gr.Name
	}

	if gr.GIDNumber != nil {
		group.GIDNumber = *gr.GIDNumber
	}

	err = g.db.Table("ldapgroups").Save(&group).Error
	if err != nil {
		return err
	}

	return nil
}

func (g *Glauth) DeleteGroup(gid int) error {
	var group models.LDAPGroup
	err := g.db.Table("ldapgroups").Where("gidnumber = ?", gid).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("group not found")
		}
		return err
	}

	err = g.db.Table("ldapgroups").Delete(&group).Error
	if err != nil {
		return err
	}

	//search for includegroups
	var includeGroups []*models.IncludeGroup
	err = g.db.Table("includegroups").Where("parentgroupid = ?", group.GIDNumber).Find(&includeGroups).Error
	if err != nil {
		return err
	}

	for _, ig := range includeGroups {
		err = g.db.Table("includegroups").Delete(&ig).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Glauth) FindNextGroupID() (int, error) {
	var group models.LDAPGroup
	err := g.db.Table("ldapgroups").Order("gidnumber desc").First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 10000, nil
		}
		return 0, err
	}

	return group.GIDNumber + 1, nil
}

func (g *Glauth) GroupExistByGID(gid int) (bool, error) {
	var group models.LDAPGroup
	err := g.db.Where("gidnumber = ?", gid).Table("ldapgroups").First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g *Glauth) GroupExistByName(name string) (bool, error) {
	var group models.LDAPGroup
	err := g.db.Where("name = ?", name).Table("ldapgroups").First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (g *Glauth) GetGroups() ([]*ressources.Group, error) {
	var groups []*models.LDAPGroup
	err := g.db.Find(&groups).Error
	if err != nil {
		return nil, err
	}

	var resGroups []*ressources.Group
	for _, g := range groups {
		resGroups = append(resGroups, &ressources.Group{
			ID:        g.ID,
			Name:      g.Name,
			GIDNumber: g.GIDNumber,
		})
	}

	return resGroups, nil
}

func GroupExistsInList(groups []*ressources.Group, target *ressources.Group) bool {
	for _, g := range groups {
		if g.GIDNumber == target.GIDNumber {
			return true
		}
	}
	return false
}

// ToCommaSeparatedString converts a list of groups to a comma separated string of GIDNumbers
func ToCommaSeparatedString(ints []int) string {
	var res []string
	for _, g := range ints {
		res = append(res, strconv.Itoa(g))
	}
	return strings.Join(res, ",")
}
