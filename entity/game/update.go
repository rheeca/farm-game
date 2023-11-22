package game

import "guion-2d-project3/utils"

func updateAnimals(g *Game) {
	for i := range g.Chickens {
		g.Chickens[i].UpdateFrame(g.CurrentFrame)
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
