package tartaglia

import "github.com/genshinsim/gsim/pkg/core"

func (c *char) Attack(p map[string]int) int {

	f := c.ActionFrames(core.ActionAttack, p)

	d := c.Snapshot(
		//fmt.Sprintf("Normal %v", c.NormalCounter),
		"Normal",
		core.AttackTagNormal,
		core.ICDTagNormalAttack,
		core.ICDGroupDefault,
		core.StrikeTypeSlash,
		core.Physical,
		25,
		attack[c.NormalCounter][c.TalentLvlAttack()],
	)
	c.QueueDmg(&d, f-1)

	c.AdvanceNormalIndex()

	return f
}
