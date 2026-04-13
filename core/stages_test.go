package core

import (
	"math"
	"testing"
)

func TestStageMultiplier_CanonicalTable(t *testing.T) {
	cases := []struct {
		stage int
		want  float64
	}{
		{-6, 2.0 / 8.0},
		{-5, 2.0 / 7.0},
		{-4, 2.0 / 6.0},
		{-3, 2.0 / 5.0},
		{-2, 0.5},
		{-1, 2.0 / 3.0},
		{0, 1.0},
		{1, 1.5},
		{2, 2.0},
		{3, 2.5},
		{4, 3.0},
		{5, 3.5},
		{6, 4.0},
	}
	for _, c := range cases {
		got := StageMultiplier(c.stage)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("StageMultiplier(%d) = %v, want %v", c.stage, got, c.want)
		}
	}
}

func TestStageMultiplier_ClampOutOfRange(t *testing.T) {
	if StageMultiplier(10) != StageMultiplier(StageMax) {
		t.Errorf("StageMultiplier(10) must clamp to StageMultiplier(%d)", StageMax)
	}
	if StageMultiplier(-10) != StageMultiplier(StageMin) {
		t.Errorf("StageMultiplier(-10) must clamp to StageMultiplier(%d)", StageMin)
	}
}

func TestAccuracyStageMultiplier_CanonicalTable(t *testing.T) {
	cases := []struct {
		stage int
		want  float64
	}{
		{-6, 3.0 / 9.0},
		{-3, 3.0 / 6.0},
		{0, 1.0},
		{3, 2.0},
		{6, 3.0},
	}
	for _, c := range cases {
		got := AccuracyStageMultiplier(c.stage)
		if math.Abs(got-c.want) > 1e-9 {
			t.Errorf("AccuracyStageMultiplier(%d) = %v, want %v", c.stage, got, c.want)
		}
	}
}

func TestApplyStage_ClampAndCompose(t *testing.T) {
	s := StatStages{}
	s = ApplyStage(s, StatAtk, 2)
	if s.Atk != 2 {
		t.Errorf("after +2 atk: want 2, got %d", s.Atk)
	}
	s = ApplyStage(s, StatAtk, 5)
	if s.Atk != StageMax {
		t.Errorf("should clamp to %d, got %d", StageMax, s.Atk)
	}
	// Other stats untouched
	if s.Def != 0 || s.SpA != 0 || s.Spe != 0 {
		t.Errorf("unexpected side effects on other stats: %+v", s)
	}
	s = ApplyStage(s, StatAtk, -100)
	if s.Atk != StageMin {
		t.Errorf("should clamp to %d, got %d", StageMin, s.Atk)
	}
}

func TestApplyStage_Immutable(t *testing.T) {
	original := StatStages{Atk: 1}
	_ = ApplyStage(original, StatAtk, 5)
	if original.Atk != 1 {
		t.Errorf("ApplyStage mutated input: got Atk=%d", original.Atk)
	}
}

func TestStatStageEffects_KnownAndUnknown(t *testing.T) {
	if effects, ok := StatStageEffects("swords-dance"); !ok || len(effects) != 1 || effects[0].Stat != StatAtk || effects[0].Delta != 2 {
		t.Errorf("swords-dance: want [+2 Atk self], got %+v ok=%v", effects, ok)
	}
	if effects, ok := StatStageEffects("Growl"); !ok || effects[0].Target != "opponent" || effects[0].Delta != -1 {
		t.Errorf("growl (case-insensitive): got %+v ok=%v", effects, ok)
	}
	if _, ok := StatStageEffects("tackle"); ok {
		t.Error("tackle should not be a stat-stage move")
	}
}

func TestStageChangeLogFragment(t *testing.T) {
	cases := []struct {
		prev, delta int
		want        string
	}{
		{0, 1, "subió"},
		{0, 2, "subió mucho"},
		{0, -1, "bajó"},
		{0, -2, "bajó mucho"},
		{StageMax, 1, "no puede subir más"},
		{StageMin, -1, "no puede bajar más"},
	}
	for _, c := range cases {
		got := StageChangeLogFragment(c.prev, c.delta)
		if got != c.want {
			t.Errorf("StageChangeLogFragment(%d,%d) = %q, want %q", c.prev, c.delta, got, c.want)
		}
	}
}
