package tactic

import (
	"github.com/golangconf/gophers-and-dragons/game"
)

func ChooseCard(s game.State) game.CardType {
	return tacticBasic(s)
}

// tacticBasic will only fight with the easiest kind of
// monsters and will run away if wounded.
// It can also use heal when it's available.
func tacticBasic(s game.State) game.CardType {
	if s.Avatar.HP < 20 {
		// or just heal here
		if s.Can(game.CardHeal) {
			return game.CardHeal
		} else if s.Creep.Traits.Has(game.TraitCoward) && s.Can(game.CardRest) {
			return game.CardRest
		}
		if s.Avatar.HP < 12 && !canKillFromOneTurn(s) {
			return game.CardRetreat
		}
	}

	// Special track on last chance to heal
	if needHealAtLast(s) {
		if s.Can(game.CardHeal) {
			return game.CardHeal
		}
	}

	if useParry(s) {
		return game.CardParry
	}

	// Stun - better to check other cards too
	if useStun(s) {
		return game.CardStun
	}

	if useFirebolt(s) {
		return game.CardFirebolt
	}
	if useAttack(s) {
		return game.CardAttack
	}
	if useMagicArrow(s) {
		return game.CardMagicArrow
	}
	if usePowerAttack(s) {
		return game.CardPowerAttack
	}

	return game.CardAttack
}

func useFirebolt(s game.State) bool {
	// Skip magic-immune creeps
	if s.Creep.Traits.Has(game.TraitMagicImmunity) {
		return false
	}

	force := s.Round > 8

	// save for Mummy
	if s.Creep.Type != game.CreepMummy && !force {
		return false
	}
	// save for mummy if it's next!
	if s.NextCreep == game.CreepMummy && force {
		return false
	}

	if s.Avatar.MP < 8 && !force {
		return false
	}

	if s.Creep.HP < 4 {
		// just kill with hands
		return false
	}

	return s.Can(game.CardFirebolt)
}

func useStun(s game.State) bool {
	if s.Creep.Type != game.CreepDragon && s.Deck[game.CardStun].Count <= 1 {
		// keep 1 for dragon
		return false
	}

	if s.Creep.MaxHP > 9 && s.Can(game.CardStun) && !s.Creep.IsStunned() && s.Creep.HP > 5 {
		return true
	}
	return false
}

func useParry(s game.State) bool {
	if s.Creep.Type != game.CreepDragon && s.Deck[game.CardParry].Count <= 1 {
		// keep 1 for dragon
		return false
	}
	switch s.Creep.Type {
	case game.CreepDragon, game.CreepKubus:
		// do nothing
	default:
		// keep cards for dragon or Kubus
		return false
	}

	if s.Creep.IsStunned() {
		return false
	}

	if s.Creep.Traits.Has(game.TraitIncrementalComplexity) && s.RoundTurn < 5 {
		return false
	}

	return s.Can(game.CardParry)
}

func useMagicArrow(s game.State) bool {
	if s.Creep.Traits.Has(game.TraitMagicImmunity) {
		return false
	}

	if s.Creep.HP == 3 && s.Can(game.CardMagicArrow) {
		return true
	}
	return false
}

func canKillFromOneTurn(s game.State) (res bool) {
	// TODO: single logic for stun
	if s.Can(game.CardStun) && s.Creep.HP > 7 {
		return true
	}
	if s.Creep.IsStunned() && s.Creep.Traits.Has(game.TraitSlow) {
		return true
	}

	if useFirebolt(s) || useAttack(s) || useMagicArrow(s) {
		return true
	}

	if useParry(s) {
		return true
	}

	if s.Creep.HP >= 4 && s.Can(game.CardPowerAttack) {
		return true
	}

	return false
}

func useAttack(s game.State) bool {
	return s.Creep.HP <= 2
}

func usePowerAttack(s game.State) bool {
	if s.Creep.Type != game.CreepDragon && s.Deck[game.CardPowerAttack].Count <= 1 {
		// keep 1 for dragon
		return false
	}
	if s.Creep.HP >= 5 && s.Can(game.CardPowerAttack) {
		return true
	}

	if s.Creep.HP >= 4 && s.Can(game.CardPowerAttack) && (s.Creep.Traits.Has(game.TraitBloodlust) || s.Round > 10) {
		return true
	}

	return false
}

func needHealAtLast(s game.State) bool {
	if s.Round < 10 && s.Round != 12 {
		return false
	}
	if s.Avatar.HP >= 32 {
		return false
	}
	if s.Creep.Traits.Has(game.TraitCoward) {
		return true
	}

	if s.Creep.Type == game.CreepKubus && s.RoundTurn <= 2 {
		return true
	}

	return false
}
