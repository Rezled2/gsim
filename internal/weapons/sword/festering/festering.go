package festering

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("festering desire", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.CR] = .045 + .015*float64(r)
	m[core.DmgP] = .12 + 0.04*float64(r)
	char.AddMod(core.CharStatMod{
		Key:    "festering",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, a == core.AttackTagElementalArt
		},
	})

}
