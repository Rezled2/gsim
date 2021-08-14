package monster

import (
	"testing"

	"github.com/genshinsim/gsim/pkg/character"
	"github.com/genshinsim/gsim/pkg/core"
)

func TestPyroAura(t *testing.T) {

	dmgCount := 0
	shdCount := 0
	var target *Target

	c, err := core.New()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	target = New(0, c, core.EnemyProfile{
		Level:  88,
		HP:     0,
		Resist: defaultResMap(),
	})
	c.Targets = append(c.Targets, target)

	c.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		dmgCount++
		return false
	}, "atk-count")

	c.Events.Subscribe(core.OnShielded, func(args ...interface{}) bool {
		shdCount++
		return false
	}, "shield-count")

	char, err := character.NewTemplateChar(c, testChar)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	c.Chars = append(c.Chars, char)

	c.Init()

	//TEST ATTACH

	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 25,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	if target.aura.Durability() != 20 {
		expect("initial durability", 20, target.aura.Durability())
		t.Error("intial attach: invalid durability")
		t.FailNow()
	}

	//TEST DECAY
	c.Skip(285)

	if !durApproxEqual(10, target.aura.Durability(), 0.01) {
		expect("decay durability after 4.75 seconds (tolerance 0.01)", 10, target.aura.Durability())
		t.Error("decay test: invalid durability")
		t.FailNow()
	}

	//TEST REFRESH
	c.Combat.ApplyDamage(&core.Snapshot{
		Durability: 50,
		Element:    core.Pyro,
		ICDTag:     core.ICDTagNone,
		ICDGroup:   core.ICDGroupDefault,
		Stats:      make([]float64, core.EndStatType),
		Targets:    core.TargetAll,
		DamageSrc:  -1,
	})

	if !durApproxEqual(60, target.aura.Durability(), 0.01) {
		expect("refresh 50 units on 10 existing (tolerance 0.01)", 60, target.aura.Durability())
		t.Error("refresh test: invalid durability")
		t.FailNow()
	}

}
