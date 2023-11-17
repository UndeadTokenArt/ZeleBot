package main

import (
)

type InitiativeTracker struct {
	entities []entity
}

func (it *InitiativeTracker) addEntity(e entity) {
	it.entities = append(it.entities, e)
}

func (it *InitiativeTracker) sortEntities() {
	// sort logic needs to be set up
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

func (it *InitiativeTracker) clearEntities() {
	it.entities = nil
}