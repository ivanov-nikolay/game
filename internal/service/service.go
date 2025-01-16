package service

import (
	"math/rand"
	"time"

	"github.com/ivanov-nikolay/game/pkg/life"
)

// LifeService - хранение состояния
type LifeService struct {
	currentWorld *life.World
	nextWorld    *life.World
}

// New создает и заполняет сетку мира
func New(height, width, fill int) (*LifeService, error) {
	rand.NewSource(time.Now().UTC().UnixNano())

	currentWorld := life.NewWorld(height, width)

	// Заполним случайными показателями, чтобы упростить пример
	currentWorld.RandInit(fill)

	newWorld := life.NewWorld(height, width)

	ls := LifeService{
		currentWorld: currentWorld,
		nextWorld:    newWorld,
	}

	return &ls, nil
}

// NewState возвращает очередное состояние игры
func (ls *LifeService) NewState() *life.World {
	life.NextState(ls.currentWorld, ls.nextWorld)

	ls.currentWorld = ls.nextWorld

	return ls.currentWorld
}
