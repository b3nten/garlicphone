// Auto-generated code for schema: messages v1

export class Player {
	id?: number;
	name?: string;
	inventory?: Foo[];
	idk?: Foo;
	nested?: Foo[][];

	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBufferView): Player;
}

export class Foo {
	bar?: number;

	static readonly TypeID: number;
	toBytes(): Uint8Array;
	fromBytes(bytes: ArrayBufferView): Foo;
}

