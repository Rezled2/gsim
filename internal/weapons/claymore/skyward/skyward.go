package skyward

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterWeaponFunc("skyward pride", weapon)
}

func weapon(char core.Character, c *core.Core, r int, param map[string]int) {

	m := make([]float64, core.EndStatType)
	m[core.DmgP] = 0.06 + float64(r)*0.02
	char.AddMod(core.CharStatMod{
		Key: "skyward pride",
		Amount: func(a core.AttackTag) ([]float64, bool) {
			return m, true
		},
		Expiry: -1,
	})

	counter := 0
	dur := 0

	dmg := 0.6 + float64(r)*0.2

	c.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
		if c.ActiveChar != char.CharIndex() {
			return false
		}
		dur = c.F + 1200
		counter = 0
		c.Log.Debugw("Skyward Pride activated", "frame", c.F, "event", core.LogWeaponEvent, "expiring ", dur)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Name()))

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		ds := args[1].(*core.Snapshot)
		//check if char is correct?
		if ds.ActorIndex != char.CharIndex() {
			return false
		}
		if ds.AttackTag != core.AttackTagNormal && ds.AttackTag != core.AttackTagExtra {
			return false
		}
		//check if cd is up
		if c.F > dur {
			return false
		}
		if counter > 8 {
			return false
		}

		counter++
		d := char.Snapshot(
			"Skyward Pride Proc",
			core.AttackTagWeaponSkill,
			core.ICDTagNone,
			core.ICDGroupDefault,
			core.StrikeTypeDefault,
			core.Physical,
			100,
			dmg,
		)
		char.QueueDmg(&d, 1)
		return false
	}, fmt.Sprintf("skyward-pride-%v", char.Name()))

}
