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
	Recipes = [10]int{
		ItemAxe, ItemHoe, ItemWateringCan, ItemBasket, ItemSeedTomato,
		ItemClock, ItemPottedSunflower, ItemPottedBlueflower, ItemPinkRug, ItemBlueRug,
	}
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
		ItemBasket: {
			{ID: ItemWeed, Count: 3},
		},
		ItemSeedTomato: {
			{ID: ItemTomato, Count: 1},
		},
		ItemClock: {
			{ID: ItemWood2, Count: 2},
			{ID: ItemRock1, Count: 2},
		},
		ItemPottedSunflower: {
			{ID: ItemSunflower, Count: 1},
			{ID: ItemRock1, Count: 1},
		},
		ItemPottedBlueflower: {
			{ID: ItemBlueflower, Count: 1},
			{ID: ItemRock1, Count: 1},
		},
		ItemPinkRug: {
			{ID: ItemWeed, Count: 3},
			{ID: ItemPinkDyeFlower, Count: 1},
		},
		ItemBlueRug: {
			{ID: ItemWeed, Count: 3},
			{ID: ItemBlueDyeFlower, Count: 1},
		},
	}
)
