package maiden

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/core"
)

func init() {
	core.RegisterSetFunc("maiden beloved", New)
}

func New(c core.Character, s *core.Core, count int) {
	if count >= 2 {
		m := make([]float64, core.EndStatType)
		m[core.Heal] = 0.15
		c.AddMod(core.CharStatMod{
			Key: "maiden-2pc",
			Amount: func(a core.AttackTag) ([]float64, bool) {
				return m, true
			},
			Expiry: -1,
		})
	}
	if count >= 4 {
		dur := 0

		s.Events.Subscribe(core.PostBurst, func(args ...interface{}) bool {
			// s.s.Log.Debugw("\t\tNoblesse 2 pc","frame",s.F, "name", ds.CharName, "abil", ds.AbilType)
			if s.ActiveChar != c.CharIndex() {
				return false
			}
			dur = s.F + 600
			s.Log.Debugw("maiden 4pc proc", "frame", s.F, "event", core.LogArtifactEvent, "char", c.CharIndex(), "expiry", dur)
			return false
		}, fmt.Sprintf("maid 4pc - %v", c.Name()))

		s.Health.AddIncHealBonus(func() float64 {
			if s.F < dur {
				return 0.2
			}
			return 0
		})
	}
	//add flat stat to char
}
