package hutao

import (
	"fmt"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/def"
	"go.uber.org/zap"
)

func init() {
	combat.RegisterCharFunc("hutao", NewChar)
}

type char struct {
	*character.Tmpl
	paraParticleICD int
	// chargeICDCounter   int
	// chargeCounterReset int
	ppBonus    float64
	tickActive bool
	c6icd      int
}

func NewChar(s def.Sim, log *zap.SugaredLogger, p def.CharacterProfile) (def.Character, error) {
	c := char{}
	t, err := character.NewTemplateChar(s, log, p)
	if err != nil {
		return nil, err
	}
	c.Tmpl = t
	c.Energy = 60
	c.EnergyMax = 60
	c.Weapon.Class = def.WeaponClassSpear
	c.NormalHitNum = 6

	c.ppHook()
	c.onExitField()
	c.a4()

	if c.Base.Cons == 6 {
		c.c6()
	}

	return &c, nil
}

/**
[11:32 PM] sakuno | yanfei is my new maid: @gimmeabreak
https://www.youtube.com/watch?v=3aCiH2U4BjY

framecounts for 7 attempts of N2CJ (no hitlag):
83, 85, 88, 89, 77, 82, 84

first 3 not from the uploaded recording (as a n1cd player i cud barely pull it off :monkaS: )
YouTube
**/

//var normalFrames = []int{13, 16, 25, 36, 44, 39}               // from kqm lib
var normalFrames = []int{10, 13, 22, 33, 41, 36} // from kqm lib, -3 for hit lag
//var dmgFrame = [][]int{{13}, {16}, {25}, {36}, {26, 44}, {39}} // from kqm lib
var dmgFrame = [][]int{{10}, {13}, {22}, {33}, {23, 41}, {36}} // from kqm lib - 3 for hit lag

func (c *char) ActionFrames(a def.ActionType, p map[string]int) int {
	switch a {
	case def.ActionAttack:
		f := normalFrames[c.NormalCounter]
		f = int(float64(f) / (1 + c.Stats[def.AtkSpd]))
		return f
	case def.ActionCharge:
		return 9 //rough.. 11, -2 for hit lag
	case def.ActionSkill:
		return 42 // from kqm lib
	case def.ActionBurst:
		return 130 // from kqm lib
	default:
		c.Log.Warnf("%v: unknown action (%v), frames invalid", c.Base.Name, a)
		return 0
	}
}

func (c *char) ActionStam(a def.ActionType, p map[string]int) float64 {
	switch a {
	case def.ActionDash:
		return 18
	case def.ActionCharge:
		if c.Sim.Status("paramita") > 0 && c.Base.Cons >= 1 {
			return 0
		}
		return 25
	default:
		c.Log.Warnf("%v ActionStam for %v not implemented; Character stam usage may be incorrect", c.Base.Name, a.String())
		return 0
	}

}

func (c *char) a4() {
	val := make([]float64, def.EndStatType)
	val[def.PyroP] = 0.33
	c.AddMod(def.CharStatMod{
		Key:    "hutao-a4",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if c.Sim.Status("paramita") == 0 {
				return nil, false
			}
			if c.HPCurrent/c.HPMax <= 0.5 {
				return val, true
			}
			return nil, false
		},
	})
}

func (c *char) c6() {
	c.Sim.AddOnHurt(func(s def.Sim) {
		c.checkc6()
	})
}

func (c *char) checkc6() {
	if c.Base.Cons < 6 {
		return
	}
	if c.Sim.Frame() < c.c6icd && c.c6icd != 0 {
		return
	}
	//check if hp less than 25%
	if c.HPCurrent/c.HPMax > .25 {
		return
	}
	//if dead, revive back to 1 hp
	if c.HPCurrent == -1 {
		c.HPCurrent = 1
	}
	//increase crit rate to 100%
	val := make([]float64, def.EndStatType)
	val[def.CR] = 1
	c.AddMod(def.CharStatMod{
		Key:    "hutao-c6",
		Amount: func(a def.AttackTag) ([]float64, bool) { return val, true },
		Expiry: c.Sim.Frame() + 600,
	})

	c.c6icd = c.Sim.Frame() + 3600
}

func (c *char) Attack(p map[string]int) int {
	f := c.ActionFrames(def.ActionAttack, p)
	hits := len(attack[c.NormalCounter])
	//check for particles
	c.ppParticles()

	d := c.Snapshot(
		fmt.Sprintf("Normal %v", c.NormalCounter),
		def.AttackTagNormal,
		def.ICDTagNormalAttack,
		def.ICDGroupDefault,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		0,
	)

	for i := 0; i < hits; i++ {
		x := d.Clone()
		x.Mult = attack[c.NormalCounter][i][c.TalentLvlAttack()]
		c.QueueDmg(&x, dmgFrame[c.NormalCounter][i])
	}

	c.AdvanceNormalIndex()

	return f
}

func (c *char) ChargeAttack(p map[string]int) int {

	f := c.ActionFrames(def.ActionCharge, p)

	if c.Sim.Status("paramita") > 0 {
		//[3:56 PM] Isu: My theory is that since E changes attack animations, it was coded
		//to not expire during any attack animation to simply avoid the case of potentially
		//trying to change animations mid-attack, but not sure how to fully test that
		//[4:41 PM] jstern25| ₼WHO_SUPREMACY: this mostly checks out
		//her e can't expire during q as well
		if f > c.Sim.Status("paramita") {
			c.Sim.AddStatus("paramita", f)
			// c.S.Status["paramita"] = c.Sim.Frame() + f //extend this to barely cover the burst
		}

		c.applyBB()
		//charge land 182, tick 432, charge 632, tick 675
		//charge land 250, tick 501, charge 712, tick 748

		//e cast at 123, animation ended 136 should end at 664 if from cast or 676 if from animation end, tick at 748 still buffed?
	}

	//check for particles
	//TODO: assuming charge can generate particles as well
	c.ppParticles()

	d := c.Snapshot(
		"Charge Attack",
		def.AttackTagExtra,
		def.ICDTagExtraAttack,
		def.ICDGroupPole,
		def.StrikeTypeSlash,
		def.Physical,
		25,
		charge[c.TalentLvlAttack()],
	)

	c.QueueDmg(&d, f-5)

	return f
}

func (c *char) ppParticles() {
	if c.Sim.Status("paramita") > 0 {
		if c.paraParticleICD < c.Sim.Frame() {
			c.paraParticleICD = c.Sim.Frame() + 300 //5 seconds
			count := 2
			if c.Sim.Rand().Float64() < 0.5 {
				count = 3
			}
			c.QueueParticle("Hutao", count, def.Pyro, dmgFrame[c.NormalCounter][0])
		}
	}
}

func (c *char) applyBB() {
	c.Log.Debugw("Applying Blood Blossom", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "current dur", c.Sim.Status("htbb"))
	//check if blood blossom already active, if active extend duration by 8 second
	//other wise start first tick func
	if !c.tickActive {
		//TODO: does BB tick immediately on first application?
		c.AddTask(c.bbtickfunc(c.Sim.Frame()), "bb", 240)
		c.tickActive = true
		c.Log.Debugw("Blood Blossom applied", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "expected end", c.Sim.Frame()+570, "next expected tick", c.Sim.Frame()+240)
	}
	// c.CD["bb"] = c.Sim.Frame() + 570 //TODO: no idea how accurate this is, does this screw up the ticks?
	c.Sim.AddStatus("htbb", 570)
	c.Log.Debugw("Blood Blossom duration extended", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "new expiry", c.Sim.Status("htbb"))
}

func (c *char) bbtickfunc(src int) func() {
	return func() {
		c.Log.Debugw("Blood Blossom checking for tick", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "cd", c.Sim.Status("htbb"), "src", src)
		if c.Sim.Status("htbb") == 0 {
			c.tickActive = false
			return
		}
		//queue up one damage instance
		d := c.Snapshot(
			"Blood Blossom",
			def.AttackTagElementalArt,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Pyro,
			25,
			bb[c.TalentLvlSkill()],
		)

		//if cons 2, add flat dmg
		if c.Base.Cons >= 2 {
			d.FlatDmg += c.HPMax * 0.1
		}
		c.Sim.ApplyDamage(&d)
		c.Log.Debugw("Blood Blossom ticked", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "next expected tick", c.Sim.Frame()+240, "dur", c.Sim.Status("htbb"), "src", src)
		//only queue if next tick buff will be active still
		// if c.Sim.Frame()+240 > c.CD["bb"] {
		// 	return
		// }
		//queue up next instance
		c.AddTask(c.bbtickfunc(src), "bb", 240)

	}
}

func (c *char) Skill(p map[string]int) int {
	//increase based on hp at cast time
	//drains hp
	c.Sim.AddStatus("paramita", 520+20) //to account for animation
	c.Log.Debugw("Paramita acivated", "frame", c.Sim.Frame(), "event", def.LogCharacterEvent, "expiry", c.Sim.Frame()+540+20)
	//figure out atk buff
	c.ppBonus = ppatk[c.TalentLvlSkill()] * c.HPMax
	max := (c.Base.Atk + c.Weapon.Atk) * 4
	if c.ppBonus > max {
		c.ppBonus = max
	}

	//remove some hp
	c.HPCurrent = 0.7 * c.HPCurrent
	c.checkc6()

	c.SetCD(def.ActionSkill, 960)
	return c.ActionFrames(def.ActionSkill, p)
}

func (c *char) ppHook() {
	val := make([]float64, def.EndStatType)
	c.AddMod(def.CharStatMod{
		Key:    "hutao-paramita",
		Expiry: -1,
		Amount: func(a def.AttackTag) ([]float64, bool) {
			if c.Sim.Status("paramita") == 0 {
				return nil, false
			}
			val[def.ATK] = c.ppBonus
			return val, true
		},
	})
}

func (c *char) onExitField() {
	c.Sim.AddEventHook(func(s def.Sim) bool {
		c.Sim.DeleteStatus("paramita")
		return false
	}, "hutao-exit", def.PostSwapHook)
}

func (c *char) Burst(p map[string]int) int {
	low := (c.HPCurrent / c.HPMax) <= 0.5
	mult := burst[c.TalentLvlBurst()]
	regen := regen[c.TalentLvlBurst()]
	if low {
		mult = burstLow[c.TalentLvlBurst()]
		regen = regenLow[c.TalentLvlBurst()]
	}
	targets := p["targets"]
	//regen for p+1 targets, max at 5; if not specified then p = 1
	count := 1
	if targets > 0 {
		count = targets
	}
	if count > 5 {
		count = 5
	}
	c.HPCurrent += c.HPMax * float64(count) * regen

	f := c.ActionFrames(def.ActionBurst, p)

	//[2:28 PM] Aluminum | Harbinger of Jank: I think the idea is that PP won't fall off before dmg hits, but other buffs aren't snapshot
	//[2:29 PM] Isu: yes, what Aluminum said. PP can't expire during the burst animation, but any other buff can
	if f > c.Sim.Status("paramita") && c.Sim.Status("paramita") > 0 {
		c.Sim.AddStatus("paramita", f) //extend this to barely cover the burst
	}

	if c.Sim.Status("paramita") > 0 && c.Base.Cons >= 2 {
		c.applyBB()
	}

	c.AddTask(func() {
		//TODO: apparently damage is based on stats on contact, not at cast
		d := c.Snapshot(
			"Spirit Soother",
			def.AttackTagElementalBurst,
			def.ICDTagNone,
			def.ICDGroupDefault,
			def.StrikeTypeDefault,
			def.Pyro,
			50,
			mult,
		)
		d.Targets = def.TargetAll
		c.Sim.ApplyDamage(&d)
	}, "Hutao Burst", f-5) //random 5 frame

	c.Energy = 0
	c.SetCD(def.ActionBurst, 900)
	return f
}

func (c *char) Snapshot(name string, a def.AttackTag, icd def.ICDTag, g def.ICDGroup, st def.StrikeType, e def.EleType, d def.Durability, mult float64) def.Snapshot {
	ds := c.Tmpl.Snapshot(name, a, icd, g, st, e, d, mult)

	if c.Sim.Status("paramita") > 0 {
		switch ds.AttackTag {
		case def.AttackTagNormal:
		case def.AttackTagExtra:
		default:
			return ds
		}
		ds.Element = def.Pyro
	}
	return ds
}
