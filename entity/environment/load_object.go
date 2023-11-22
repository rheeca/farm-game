package environment

import (
	"github.com/lafriks/go-tiled"
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"
)

func loadObject32(objects []model.Object, mapObjectType, itemType int, gMap *tiled.Map) []model.Object {
	for _, mapObj := range gMap.Groups[0].ObjectGroups[mapObjectType].Objects {
		wood := model.Object{
			Type: itemType,
			XLoc: int(mapObj.X),
			YLoc: int(mapObj.Y),
			Sprite: model.SpriteBody{
				X:      int(mapObj.X),
				Y:      int(mapObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			Collision: model.CollisionBody{
				X:      int(mapObj.X),
				Y:      int(mapObj.Y),
				Width:  utils.UnitSize,
				Height: utils.UnitSize,
			},
			IsCollision: true,
		}
		objects = append(objects, wood)
	}
	return objects
}
