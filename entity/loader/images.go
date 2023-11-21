package loader

import (
	"embed"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageCollection struct {
	CraftingTable *ebiten.Image
	CraftingUI    *ebiten.Image
	FarmItems     *ebiten.Image
	SelectedTool  *ebiten.Image
	ToolsUI       *ebiten.Image
	TreeSprites   *ebiten.Image
}

func NewImageCollection(EmbeddedAssets embed.FS) (images ImageCollection) {
	return ImageCollection{
		CraftingTable: loadImage(EmbeddedAssets, path.Join("assets", "items", "crafting_table.png")),
		CraftingUI:    loadImage(EmbeddedAssets, path.Join("assets", "ui", "crafting_ui.png")),
		FarmItems:     loadImage(EmbeddedAssets, path.Join("assets", "items", "farm_items.png")),
		SelectedTool:  loadImage(EmbeddedAssets, path.Join("assets", "ui", "selected_tool.png")),
		ToolsUI:       loadImage(EmbeddedAssets, path.Join("assets", "ui", "tools_ui.png")),
		TreeSprites:   loadImage(EmbeddedAssets, path.Join("assets", "items", "tree_sprites.png")),
	}
}

func loadImage(EmbeddedAssets embed.FS, filepath string) *ebiten.Image {
	embeddedFile, err := EmbeddedAssets.Open(filepath)
	if err != nil {
		return nil
	}
	image, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		return nil
	}
	return image
}
