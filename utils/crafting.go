package utils

const (
	CraftingUIFirstBoxX  = 176
	CraftingUIFirstBoxY  = 144
	CraftingUIFirstSlotX = 192
	CraftingUIFirstSlotY = 160
	CraftingUISpacing    = 96
	CraftingUIColumns    = 5
	CraftingUIBoxCount   = 10

	CraftingUIBoxCollisionX      = 182
	CraftingUIBoxCollisionY      = 150
	CraftingUIBoxCollisionWidth  = 52
	CraftingUIBoxCollisionHeight = 52

	RecipeItemX = 224
	RecipeItemY = 384
)

type Material struct {
	ID    int
	Count int
}

var (
	Recipes       = [10]int{ItemAxe, ItemHoe, ItemWateringCan}
	RecipeDetails = map[int][]Material{
		ItemAxe: {
			{ID: ItemWood2, Count: 1},
			{ID: ItemRock1, Count: 2},
		},
		ItemHoe: {
			{ID: ItemWood2, Count: 1},
			{ID: ItemRock1, Count: 1},
		},
		ItemWateringCan: {
			{ID: ItemRock1, Count: 3},
		},
	}
)
