package kazuha

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) ActionFrames(a core.ActionType, p map[string]int) int {
	switch a {
	case core.ActionAttack:
		f := 0
		switch c.NormalCounter {
		//TODO: need to add atkspd mod
		case 0:
			//add frames if last action is also attack
			if c.Core.LastAction.Target == "kazuha" && c.Core.LastAction.Typ == core.ActionAttack {
				f += 60
			}
			f = 14
		case 1:
			f = 34 - 14
		case 2:
			f = 70 - 34 //hit at 60, 70
		case 3:
			f = 97 - 70
		case 4:
			f = 140 - 97
		}
		f = int(float64(f) / (1 + c.Stats[core.AtkSpd]))
		return f
	case core.ActionHighPlunge:
		if c.Core.LastAction.Target == "kazuha" && c.Core.LastAction.Typ == core.ActionAttack {
			_, ok := c.Core.LastAction.Param["hold"]
			if ok {
				return 63
			}
			return 55
		}
		c.Core.Log.Warnw("invalid plunge", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0
	case core.ActionSkill:
		_, ok := p["hold"]
		if ok {
			return 69
		}
		return 36
	case core.ActionBurst:
		return 93
	default:
		c.Core.Log.Warnw("unknown action", "event", core.LogActionEvent, "frame", c.Core.F, "action", a)
		return 0
	}
}
