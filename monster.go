package main

type monster struct {
<<<<<<< HEAD
	name       string
	damage     int
	initiative int
=======
  name      string
  damage    int
  initiative int
>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
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
func (m *monster) getInitiative() int {
	return m.initiative
}
<<<<<<< HEAD
=======
func (m *monster) getInitiative() int {
  return m.initiative
}
>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
