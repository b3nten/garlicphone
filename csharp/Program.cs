// See https://aka.ms/new-console-template for more information

using Messages;

Item item = new();

item.Name = "Sword";

var player = new Player();
player.Id = 1;
player.Name = "Benton";
player.Inventory = new Item[] { item };
player.Foo = "Bar";
player.Lol = new uint[][] { new uint[] { 1, 2 }, new uint[] { 3, 4 } };
player.Lol2 = new Item[][][] { new Item[][] { new Item[] { item } } };
player.Dead = false;

var bytes = player.Serialize();

// write to file
File.WriteAllBytes("player.dat", bytes);
