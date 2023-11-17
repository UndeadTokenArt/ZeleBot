package main

type player struct {
	name          string
	initiative    int
	currentHealth int
	AC            int
}

func (p *player) getCurrentDmgDone() int {
	return 0 // Implement the logic if needed
}

func (p *player) attack(target entity, dmg int) {
	switch v := target.(type) {
	case *player:
		v.currentHealth = v.currentHealth - dmg
	case *monster:
		v.damage = v.damage + dmg
	}
}
func (p *player) getInitiative() int {
	return p.initiative
}
