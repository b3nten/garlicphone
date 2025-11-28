version = 1

item = struct {
	name = str{ id = 1 },
}

player = struct {
	id = uint32{ id = 10 },
	name = str{ id = 11 },
	inventory = list(item){ id = 12 },
}
