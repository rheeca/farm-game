package loader

import (
	"embed"
	"fmt"
	"guion-2d-project3/utils"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageCollection struct {
	FarmItems    *ebiten.Image
	ToolsUI      *ebiten.Image
	SelectedTool *ebiten.Image
}

func NewImageCollection(EmbeddedAssets embed.FS) (images ImageCollection, err error) {
	return loadImages(EmbeddedAssets)
}

func loadImages(EmbeddedAssets embed.FS) (images ImageCollection, err error) {

	embeddedFile, err := EmbeddedAssets.Open(path.Join("assets", "items", "farm_items.png"))
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEmbeddedImage, err)
	}
	farmItemsImg, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEbitenImage, err)
	}

	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", "ui", "tools_ui.png"))
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEmbeddedImage, err)
	}
	toolsUIImg, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEbitenImage, err)
	}

	embeddedFile, err = EmbeddedAssets.Open(path.Join("assets", "ui", "selected_tool.png"))
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEmbeddedImage, err)
	}
	selToolImg, _, err := ebitenutil.NewImageFromReader(embeddedFile)
	if err != nil {
		return ImageCollection{}, fmt.Errorf(utils.ErrorLoadEbitenImage, err)
	}

	return ImageCollection{
		FarmItems:    farmItemsImg,
		ToolsUI:      toolsUIImg,
		SelectedTool: selToolImg,
	}, nil
}
