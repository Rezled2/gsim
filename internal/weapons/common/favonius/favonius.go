package favonius

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("favonius warbow", weapon)
	core.RegisterWeaponFunc("favonius sword", weapon)
	core.RegisterWeaponFunc("favonius lance", weapon)
	core.RegisterWeaponFunc("favonius greatsword", weapon)
	core.RegisterWeaponFunc("favonius codex", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	p := 0.50 + float64(r)*0.1
	cd := 810 - r*90
	icd := 0
	//add on crit effect
	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		crit := args[3].(bool)
		if !crit {
			return false
		}
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		if icd > c.F {
			return false
		}

		if c.Rand.Float64() > p {
			return false
		}
		c.Log.Debugw("favonius proc'd", "frame", c.F, "event", core.LogWeaponEvent, "char", char.CharIndex())

		char.QueueParticle("favonius", 3, core.NoElement, 150)

		icd = c.F + cd

		return false
	}, fmt.Sprintf("favo-%v", char.Name()))

}
