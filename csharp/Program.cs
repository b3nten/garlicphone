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

// read from file
var readBytes = File.ReadAllBytes("player.dat");
var deserializedPlayer = new Player().Deserialize(readBytes);
Console.WriteLine($"Player Name: {deserializedPlayer.Name}");
Console.WriteLine($"Player Inventory Item Name: {deserializedPlayer.Inventory?[0].Name}");
Console.WriteLine($"Player Foo: {deserializedPlayer.Foo}");
Console.WriteLine($"Player Lol: {string.Join(", ", deserializedPlayer.Lol?[0] ?? Array.Empty<uint>())}");
Console.WriteLine($"Player Lol2 Item Name: {deserializedPlayer.Lol2?[0]?[0]?[0].Name}");
Console.WriteLine($"Player Dead: {deserializedPlayer.Dead}");
