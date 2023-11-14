package main

type InitiativeTracker struct {
	entities []entity
}

func (it *InitiativeTracker) addEntity(e entity) {
	it.entities = append(it.entities, e)
}

func (it *InitiativeTracker) sortEntities() {
<<<<<<< HEAD
	// figure out sorting logic
=======
  // sort algorithm
>>>>>>> 81e7a89d770dd06b697144b471b3f04b8ab74710
}

func (it *InitiativeTracker) getCurrentEntity() entity {
	if len(it.entities) > 0 {
		return it.entities[0]
	}
	return nil
}

func (it *InitiativeTracker) nextTurn() {
	if len(it.entities) > 0 {
		// Move the current entity to the end of the slice
		it.entities = append(it.entities[1:], it.entities[0])
	}
}
