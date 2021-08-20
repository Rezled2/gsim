package solar

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("solar pearl", weapon)
}

//Normal Attack hits increase Elemental Skill and Elemental Burst DMG by 20/25/30/35/40% for 6s.
//Likewise, Elemental Skill or Elmental Burst hits increase Normal Attack DMG by 20/25/30/35/40% for 6s.
func weapon(char core.Character, c *core.Core, r int, param map[string]int) {
	skill := 0
	attack := 0

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag == core.AttackTagElementalArt || ds.AttackTag == core.AttackTagElementalBurst {
			skill = c.F + 300
			return false
		}
		if ds.AttackTag == core.AttackTagNormal {
			skill = c.F + 300
		}
		return false
	}, fmt.Sprintf("solar-%v", char.Name()))

	val := make([]float64, core.EndStatType)
	val[core.DmgP] = 0.15 + float64(r)*0.05
	char.AddMod(core.CharStatMod{
		Key:    "solar",
		Expiry: -1,
		Amount: func(a core.AttackTag) ([]float64, bool) {
			if a == core.AttackTagElementalArt || a == core.AttackTagElementalBurst {
				return val, attack > c.F
			}
			if a == core.AttackTagNormal {
				return val, skill > c.F
			}
			return nil, false
		},
	})
}
