package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"guion-2d-project3/utils"
	"image"
)

func drawMap(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, layer := range g.Environment.Maps[g.CurrentMap].Layers {
		if layer.Name == utils.GuideOnlyLayer {
			continue
		}
		for tileY := 0; tileY < utils.MapRows; tileY += 1 {
			for tileX := 0; tileX < utils.MapColumns; tileX += 1 {
				// find img of tile to draw
				tileToDraw := layer.Tiles[tileY*utils.MapColumns+tileX]
				if tileToDraw.IsNil() {
					continue
				}

				tileToDrawX := int(tileToDraw.ID) % tileToDraw.Tileset.Columns
				tileToDrawY := int(tileToDraw.ID) / tileToDraw.Tileset.Columns

				ebitenTileToDraw := g.Environment.Tilesets[tileToDraw.Tileset.Name].SubImage(image.Rect(tileToDrawX*utils.TileWidth,
					tileToDrawY*utils.TileHeight,
					tileToDrawX*utils.TileWidth+utils.TileWidth,
					tileToDrawY*utils.TileHeight+utils.TileHeight)).(*ebiten.Image)

				// draw tile
				drawOptions.GeoM.Reset()
				TileXpos := float64(utils.TileWidth * tileX)
				TileYpos := float64(utils.TileHeight * tileY)
				drawOptions.GeoM.Translate(TileXpos, TileYpos)
				screen.DrawImage(ebitenTileToDraw, &drawOptions)
			}
		}
	}
}

func drawTrees(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	for _, t := range g.Environment.Trees {
		drawOptions.GeoM.Reset()
		drawOptions.GeoM.Translate(float64(t.XLoc), float64(t.YLoc))
		screen.DrawImage(g.Images.TreeSprites.SubImage(image.Rect(t.Frame*t.Sprite.Width,
			0,
			t.Frame*t.Sprite.Width+t.Sprite.Width,
			t.Sprite.Height)).(*ebiten.Image), &drawOptions)
	}
}

func drawObjects(g *Game, screen *ebiten.Image, drawOptions ebiten.DrawImageOptions) {
	mapObjects := g.Environment.Objects[g.CurrentMap]
	for _, o := range mapObjects {
		var x0, y0, x1, y1 int
		var objImage *ebiten.Image

		drawOptions.GeoM.Reset()
		if o.Type == utils.ItemCraftingTable {
			objImage = g.Images.CraftingTable
			x0, y0 = 0, 0
			x1 = objImage.Bounds().Dx()
			y1 = objImage.Bounds().Dx()
		}

		drawOptions.GeoM.Translate(float64(o.XLoc), float64(o.YLoc))
		screen.DrawImage(objImage.SubImage(image.Rect(x0, y0, x1, y1)).(*ebiten.Image), &drawOptions)
	}
}

// DrawCenteredText %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
// from https://github.com/sedyh/ebitengine-cheatsheet
func DrawCenteredText(screen *ebiten.Image, font font.Face, s string, cx, cy int) {
	bounds := text.BoundString(font, s)
	x, y := cx-bounds.Min.X-bounds.Dx()/2, cy-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, colornames.Brown)
}

// %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
