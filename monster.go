package main

type monster struct {
  name      string
  damage    int
  initative int
}

func (m *monster) attack(target entity, dmg int) {
  switch v := target.(type) {
  case *player:
    v.currentHealth = v.currentHealth - dmg
  case *monster:
    v.damage = v.damage + dmg
  }
}

func (m *monster) getCurrentDmgDone() int {
  return m.damage
}
func (p *player) getInitiative() int {
	return p.initiative
}