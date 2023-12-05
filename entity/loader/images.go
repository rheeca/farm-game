package loader

import (
	"embed"
	"fmt"
	"guion-2d-project3/utils"
	"log"
	"path"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ImageCollection struct {
	BedPink                  *ebiten.Image
	BlackScreen              *ebiten.Image
	ButtonDelete             *ebiten.Image
	CharacterCustomizationUI *ebiten.Image
	Characters               []*ebiten.Image
	CraftingTable            *ebiten.Image
	CraftingUI               *ebiten.Image
	DoorSprites              *ebiten.Image
	FarmingPlants            *ebiten.Image
	FarmItems                *ebiten.Image
	SelectedCharacter        *ebiten.Image
	SelectedItem             *ebiten.Image
	SelectedTool             *ebiten.Image
	Tilesets                 map[string]*ebiten.Image
	ToolsUI                  *ebiten.Image
	TreeSprites              *ebiten.Image
}

func NewImageCollection(EmbeddedAssets embed.FS) (images ImageCollection) {
	characters := []*ebiten.Image{
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_white.png")),
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_purple.png")),
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_pink.png")),
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_aqua.png")),
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_green.png")),
		loadImage(EmbeddedAssets, path.Join("client", "assets", "player", "player_blue.png")),
	}
	return ImageCollection{
		BedPink:                  loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "bed_pink.png")),
		BlackScreen:              loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "black_screen.png")),
		ButtonDelete:             loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "button_delete.png")),
		CharacterCustomizationUI: loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "character_customization_ui.png")),
		Characters:               characters,
		CraftingTable:            loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "crafting_table.png")),
		CraftingUI:               loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "crafting_ui.png")),
		DoorSprites:              loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "door_sprites.png")),
		FarmingPlants:            loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "farming_plants.png")),
		FarmItems:                loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "farm_items.png")),
		SelectedCharacter:        loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "selected_character.png")),
		SelectedItem:             loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "selected_item.png")),
		SelectedTool:             loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "selected_tool.png")),
		Tilesets:                 loadTilesets(EmbeddedAssets),
		ToolsUI:                  loadImage(EmbeddedAssets, path.Join("client", "assets", "ui", "tools_ui.png")),
		TreeSprites:              loadImage(EmbeddedAssets, path.Join("client", "assets", "items", "tree_sprites.png")),
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

func loadTilesets(embeddedAssets embed.FS) map[string]*ebiten.Image {
	tilesets := map[string]*ebiten.Image{}
	for _, tsPath := range utils.Tilesets {
		embeddedFile, err := embeddedAssets.Open(path.Join("client", "assets", "tilesets", tsPath))
		if err != nil {
			log.Fatal("failed to load embedded image:", embeddedFile, err)
		}
		tsImage, _, err := ebitenutil.NewImageFromReader(embeddedFile)
		if err != nil {
			fmt.Println("error loading tileset image")
		}
		tilesets[strings.Split(tsPath, ".")[0]] = tsImage
	}
	return tilesets
}
