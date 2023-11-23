package game

import (
	"guion-2d-project3/utils"
)

func updateAnimals(g *Game) {
	for i := range g.Chickens {
		g.Chickens[i].UpdateFrame(g.CurrentFrame)
		g.Chickens[i].RandomMovement(g.CurrentFrame)
		if g.Chickens[i].State == utils.ChickenWalkState {
			moveChickenTowardsDestination(g, i)
		}
	}
	for i := range g.Cows {
		g.Cows[i].UpdateFrame(g.CurrentFrame)
	}
}

func updateTrees(g *Game) {
	for i, t := range g.Environment.Trees {
		if t.IsNil {
			continue
		}
		g.Environment.Trees[i].UpdateFrame(g.CurrentFrame)
	}
}

func updateObjects(g *Game) {
	for i, o := range g.Environment.Objects[g.CurrentMap] {
		if o.IsNil {
			continue
		}
		if o.Type == utils.ItemDoor {
			g.Environment.Objects[g.CurrentMap][i].UpdateFrame(g.CurrentFrame)
		}
	}
}

func moveChickenTowardsDestination(g *Game, c int) {
	if g.Chickens[c].XLoc > g.Chickens[c].Destination.X*utils.TileWidth {
		g.Chickens[c].Direction = utils.AnimalLeft
		g.Chickens[c].Dx -= utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Chickens[c].UpdateLocation()
		} else {
			g.Chickens[c].Dx = 0
		}
	} else if g.Chickens[c].XLoc < g.Chickens[c].Destination.X*utils.TileWidth {
		g.Chickens[c].Direction = utils.AnimalRight
		g.Chickens[c].Dx += utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Chickens[c].UpdateLocation()
		} else {
			g.Chickens[c].Dx = 0
		}
	} else if g.Chickens[c].YLoc > g.Chickens[c].Destination.Y*utils.TileHeight {
		g.Chickens[c].Dy -= utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Chickens[c].UpdateLocation()
		} else {
			g.Chickens[c].Dy = 0
		}
	} else if g.Chickens[c].YLoc < g.Chickens[c].Destination.Y*utils.TileHeight {
		g.Chickens[c].Dy += utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Chickens[c].UpdateLocation()
		} else {
			g.Chickens[c].Dy = 0
		}
	}
}
