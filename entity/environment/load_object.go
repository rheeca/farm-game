package environment

import (
	"guion-2d-project3/entity/model"
	"guion-2d-project3/utils"

	"github.com/lafriks/go-tiled"
)

func loadObject32(objects []model.Object, mapObjectType, itemType int, gMap *tiled.Map) []model.Object {
	for _, mapObj := range gMap.Groups[0].ObjectGroups[mapObjectType].Objects {
		obj := model.Object{
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
		objects = append(objects, obj)
	}
	return objects
}
