package loader

import (
	"embed"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageCollection struct {
	FarmItems    *ebiten.Image
	ToolsUI      *ebiten.Image
	SelectedTool *ebiten.Image
	TreeSprites  *ebiten.Image
}

func NewImageCollection(EmbeddedAssets embed.FS) (images ImageCollection, err error) {
	return loadImages(EmbeddedAssets)
}

func loadImages(EmbeddedAssets embed.FS) (images ImageCollection, err error) {
	return ImageCollection{
		FarmItems:    loadImage(EmbeddedAssets, path.Join("assets", "items", "farm_items.png")),
		ToolsUI:      loadImage(EmbeddedAssets, path.Join("assets", "ui", "tools_ui.png")),
		SelectedTool: loadImage(EmbeddedAssets, path.Join("assets", "ui", "selected_tool.png")),
		TreeSprites:  loadImage(EmbeddedAssets, path.Join("assets", "items", "tree_sprites.png")),
	}, nil
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
