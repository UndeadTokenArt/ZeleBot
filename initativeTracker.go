package main

import (
	"sort"
)

type InitiativeTracker struct {
	entities []entity
}

func (it *InitiativeTracker) AddEntity(e entity) {
	it.entities = append(it.entities, e)
}

func (it *InitiativeTracker) SortEntities() {
	sort.Slice(it.entities, func(i, j int) bool {
		return it.entities[i].getInitiative() > it.entities[j].getInitiative()
	})
}

func (it *InitiativeTracker) GetCurrentEntity() entity {
	if len(it.entities) > 0 {
		return it.entities[0]
	}
	return nil
}

func (it *InitiativeTracker) NextTurn() {
	if len(it.entities) > 0 {
		// Move the current entity to the end of the slice
		it.entities = append(it.entities[1:], it.entities[0])
	}
}

func (it *InitiativeTracker) ClearEntities() {
	it.entities = nil
}
