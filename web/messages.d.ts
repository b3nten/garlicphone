// Auto-generated code for schema: messages v1

export class Player {
	id?: number;
	name?: string;
	inventory?: Item[];
	foo?: string;
	dead?: boolean;
	lol?: number[][];

	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBuffer | ArrayBufferView): Player;
}

export class Item {
	name?: string;

	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBuffer | ArrayBufferView): Item;
}

export function deserialize(bytes: ArrayBuffer | ArrayBufferView): Player | Item;

