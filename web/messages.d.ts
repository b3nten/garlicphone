// Auto-generated code for schema: messages v1

export class Item {
	name?: string;
	constructor(props?: Omit<Partial<Item>, 'fromBytes' | 'toBytes'>);
	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBuffer | ArrayBufferView): Item;
}

export class Player {
	name?: string;
	inventory?: Item[];
	foo?: string;
	dead?: boolean;
	lol?: number[][];
	lol2?: Item[][][];
	id?: number;
	constructor(props?: Omit<Partial<Player>, 'fromBytes' | 'toBytes'>);
	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBuffer | ArrayBufferView): Player;
}

export function deserialize(bytes: ArrayBuffer | ArrayBufferView): Item | Player;

