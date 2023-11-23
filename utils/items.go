package utils

// farm items
const (
	ItemAxe              int = 2
	ItemHoe                  = 3
	ItemWateringCan          = 10
	ItemWood2                = 19
	ItemSeedTomato           = 32
	ItemTomato               = 33
	ItemRock1                = 42
	ItemSunflower            = 99
	ItemBlueflower           = 100
	ItemWeed                 = 101
	ItemPinkDyeFlower        = 102
	ItemBlueDyeFlower        = 103
	ItemChair                = 106
	ItemClock                = 107
	ItemPottedSunflower      = 108
	ItemPottedBlueflower     = 109
	ItemPinkRug              = 110
	ItemBlueRug              = 111
	ItemBasket               = 114
	ItemCraftingTable        = 121
	ItemDoor                 = 122
	ItemBedPink              = 123
	ItemMapStone3            = 124
	ItemMapWood              = 125
)

// plants
const (
	PlantTomato = 2
)

// flowers and stones
const (
	MapStone3        = 14
	MapWeed          = 27
	MapSunflower     = 38
	MapPinkDyeFlower = 43
	MapBlueDyeFlower = 55
	MapBlueflower    = 57
)

// trees
const (
	MapWood = 78
)

const (
	FarmItemsColumns = 8
)

var (
	PlantItemMapping = map[int]int{
		PlantTomato: ItemTomato,
	}
)

func IsSeed(item int) bool {
	seedItems := []int{
		ItemSeedTomato,
	}
	for _, s := range seedItems {
		if s == item {
			return true
		}
	}
	return false
}
