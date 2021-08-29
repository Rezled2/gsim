package tartaglia

import (
	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterCharFunc("tartaglia", NewChar)
}

type char struct {
	*character.Tmpl
	gundamMode bool
}

func NewChar(s *core.Core, p core.CharacterProfile) (core.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = core.WeaponClassBow
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 6

	return &c, nil
}
