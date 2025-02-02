# 2D Game Design - Project 3
By Rheeca Guion

![alt-text](client/assets/icons/chicken.png)
![alt-text](client/assets/icons/character.png)
![alt-text](client/assets/icons/character_purple.png)
![alt-text](client/assets/icons/potted_sunflower.png)

This project is a farming and crafting game. All requirements for Project 3 have been completed.

For Part 2, a multiplayer feature has been added. Players can work together to gather resources, craft materials and plant crops.

## How to run
To run the game, there must be a player that serves as the host. Other players can join as clients.
### Server
To start a host player, the game must be run from the root directory. This runs the main function in `server.go`.

### Client
To play the game as a client, the game must be run from the client directory. This runs the main function in `/client/client.go`.


## Controls
#### WASD
* Player movement
#### Left-click mouse
* For using tools and planting seeds
#### Right-click mouse
* Interact with objects, plants and animals on the map
#### Number keys (1-9)
* Can be used to equip items from the backpack; clicking on the backpack slot with the mouse also equips the item. The delete button on the right deletes the currently equipped item.

## How to play
* Start by choosing a character.
* The player starts on the farm map where they can pick up debris. These items can be used to craft tools and other items using the crafting table.
* The forest map has more plants that can be foraged and used for crafting. To gather wood from trees, an axe must be equipped.
* The animal map has chickens and cows that the player can interact with to show a heart. Chickens walk randomly around the map.
* Sleeping refreshes the resources in the forest.

### Farming
* Seeds can only be planted on tillable land on the farm map.
* The hoe tool must first be used to till the grassland. The resulting tilled plot can be watered and a seed can be planted in it.
* The plant will grow the next day and can be harvested.


## Credits
Art: _Sprout Lands_ (c) [Cup Nooble](https://cupnooble.itch.io/)

Background music: _First Town_ (c) Karl Kopio
