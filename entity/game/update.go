package game

import (
	"guion-2d-project3/utils"
)

func updateAnimals(g *Game) {
	for i := range g.Data.Chickens {
		g.Data.Chickens[i].UpdateFrame(g.CurrentFrame)
		g.Data.Chickens[i].RandomMovement(g.CurrentFrame)
		if g.Data.Chickens[i].State == utils.ChickenWalkState {
			moveChickenTowardsDestination(g, i)
		}
	}
	for i := range g.Data.Cows {
		g.Data.Cows[i].UpdateFrame(g.CurrentFrame)
	}
}

func updateTrees(g *Game) {
	for i, t := range g.Data.Environment.Trees {
		if t.IsNil {
			continue
		}
		g.Data.Environment.Trees[i].UpdateFrame(g.CurrentFrame)
	}
}

func updateObjects(g *Game) {
	for i, o := range g.Data.Environment.Objects[g.CurrentMap] {
		if o.IsNil {
			continue
		}
		if o.Type == utils.ItemDoor {
			g.Data.Environment.Objects[g.CurrentMap][i].UpdateFrame(g.CurrentFrame)
		}
	}
}

func moveChickenTowardsDestination(g *Game, c int) {
	if g.Data.Chickens[c].XLoc > g.Data.Chickens[c].Destination.X*utils.TileWidth {
		g.Data.Chickens[c].Direction = utils.AnimalLeft
		g.Data.Chickens[c].Dx -= utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Data.Chickens[c].UpdateLocation()
		} else {
			g.Data.Chickens[c].Dx = 0
		}
	} else if g.Data.Chickens[c].XLoc < g.Data.Chickens[c].Destination.X*utils.TileWidth {
		g.Data.Chickens[c].Direction = utils.AnimalRight
		g.Data.Chickens[c].Dx += utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Data.Chickens[c].UpdateLocation()
		} else {
			g.Data.Chickens[c].Dx = 0
		}
	} else if g.Data.Chickens[c].YLoc > g.Data.Chickens[c].Destination.Y*utils.TileHeight {
		g.Data.Chickens[c].Dy -= utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Data.Chickens[c].UpdateLocation()
		} else {
			g.Data.Chickens[c].Dy = 0
		}
	} else if g.Data.Chickens[c].YLoc < g.Data.Chickens[c].Destination.Y*utils.TileHeight {
		g.Data.Chickens[c].Dy += utils.AnimalMovementSpeed
		if !chickenHasCollisions(g, c) {
			g.Data.Chickens[c].UpdateLocation()
		} else {
			g.Data.Chickens[c].Dy = 0
		}
	}
}
