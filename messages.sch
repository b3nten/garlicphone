version = 1

item = struct {
	name = str{ id = 1 },
}

player = struct {
	id = uint32{ id = 10 },
	name = str{ id = 11 },
	inventory = list(item){ id = 12 },
	foo = str{ id = 13 },
	dead = bool{ id = 14 },
	lol = list(list(uint32)){ id = 15 },
}
