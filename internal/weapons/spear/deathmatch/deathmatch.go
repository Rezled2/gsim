package deathmatch

import (
	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("deathmatch", weapon)
}

//If there are at least 2 opponents nearby, ATK is increased by 16% and DEF is increased by 16%.
//If there are fewer than 2 opponents nearby, ATK is increased by 24%.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	single := make([]float64, core.EndStatType)
	single[core.ATKP] = .12 + .04*float64(r)
	single[core.DEFP] = .12 + .04*float64(r)

	multiple := make([]float64, core.EndStatType)
	multiple[core.ATKP] = .18 + .06*float64(r)

	char.AddMod(core.CharStatMod{
		Key:    "deathmatch",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if len(c.Targets) > 1 {
				return multiple, true
			}
			return single, true
		},
	})

}
