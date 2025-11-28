// Auto-generated code for schema: messages v1

export class Player {
	static get TypeID() { return 49920; }
	id; 	name; 	inventory;
	toBytes() { return this.__serialize(new ByteBuffer()).bytes(); }
	fromBytes(bytes) {
		if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
		parse_struct(this, new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength), 0);
		return this;
	}
}
Object.defineProperty(Player, '__deserialize', { value: create_static_deserializer(Player), enumerable: false })
function Player___serialize(b) {
	b.write_uint16(49920);
	const start_index = b.length;
	b.write_uint32(0);
	if(this.id !== 'undefined') {
		b.write_uint16(10);
		b.write_uint32(this.id);
	}
	if(this.name !== 'undefined') {
		b.write_uint16(11);
		b.write_string(this.name);
	}
	if(this.inventory !== 'undefined') {
		b.write_uint16(12);
		b.list_writer(b.write_struct)(this.inventory)
	}
	const end_index = b.length;
	b.set_uint32(start_index, end_index - (start_index + 4))
	return b;
}
Object.defineProperty(Player.prototype, '__serialize', { value: Player___serialize , enumerable: false })
function Player___deserialize_field(view, fieldID, offset) {
	switch(fieldID) {
		case 10: return deserialize_uint32(view, offset, this, 'id')
		case 11: return deserialize_string(view, offset, this, 'name')
		case 12: return list_deserializer(Item.__deserialize)(view, offset, this, 'inventory')
		default:
			return unknown_field;
	}
}
Object.defineProperty(Player.prototype, '__deserialize_field', { value: Player___deserialize_field , enumerable: false })

export class Item {
	static get TypeID() { return 13286; }
	name;
	toBytes() { return this.__serialize(new ByteBuffer()).bytes(); }
	fromBytes(bytes) {
		if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
		parse_struct(this, new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength), 0);
		return this;
	}
}
Object.defineProperty(Item, '__deserialize', { value: create_static_deserializer(Item), enumerable: false })
function Item___serialize(b) {
	b.write_uint16(13286);
	const start_index = b.length;
	b.write_uint32(0);
	if(this.name !== 'undefined') {
		b.write_uint16(1);
		b.write_string(this.name);
	}
	const end_index = b.length;
	b.set_uint32(start_index, end_index - (start_index + 4))
	return b;
}
Object.defineProperty(Item.prototype, '__serialize', { value: Item___serialize , enumerable: false })
function Item___deserialize_field(view, fieldID, offset) {
	switch(fieldID) {
		case 1: return deserialize_string(view, offset, this, 'name')
		default:
			return unknown_field;
	}
}
Object.defineProperty(Item.prototype, '__deserialize_field', { value: Item___deserialize_field , enumerable: false })

export function deserialize(bytes) {
	if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
	const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
	const typeID = view.getUint16(0);
	switch(typeID) {
		case 49920: return new Player().fromBytes(bytes);
		case 13286: return new Item().fromBytes(bytes);
		default: throw new Error(`Unknown TypeID: ${typeID}`);
	}
}

class ByteBuffer {
	get length() { return this.#len; }

	write = (value) => {
		this.#resize(this.#len + value.length);
		this.#view.set(value, this.#len - value.length);
	}

	write_bool = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setUint8(this.#len - 1, value ? 1 : 0);
	}

	write_int8 = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setInt8(this.#len - 1, value);
	}

	write_uint8 = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setUint8(this.#len - 1, value);
	}

	write_int16 = (value) => {
		this.#resize(this.#len + 2);
		this.#dview.setInt16(this.#len - 2, value);
	}

	write_uint16 = (value) => {
		this.#resize(this.#len + 2);
		this.#dview.setUint16(this.#len - 2, value);
	}

	write_int32 = (value) => {
		this.#resize(this.#len + 4);
		this.#dview.setInt32(this.#len - 4, value);
	}

	write_uint32 = (value) => {
		this.#resize(this.#len + 4);
		this.#dview.setUint32(this.#len - 4, value);
	}

	write_string = (value) => {
		const encoded = this.#encoder.encode(value);
		this.set_uint32(this.#len, encoded.length);
		this.write(encoded);
	}

	write_struct = (value) => {
		value.__serialize(this);
	}

	list_writer = (s) => (value) => {
		this.write_uint32(0); // Placeholder for size
		const sizeIndex = this.length;
		for (const item of value) {
			s(item);
		}
		this.set_uint32(sizeIndex - 4, this.#len - sizeIndex);
	}

	set_uint8(offset, value) {
		this.#resize(offset + 1);
		this.#dview.setUint8(offset, value);
	}

	set_uint16(offset, value) {
		this.#resize(offset + 2);
		this.#dview.setUint16(offset, value);
	}

	set_uint32(offset, value) {
		this.#resize(offset + 4);
		this.#dview.setUint32(offset, value);
	}

	bytes() {
		return new Uint8Array(this.#buffer, 0, this.#len);
	}

	#encoder = new TextEncoder();
	#buffer = new ArrayBuffer(0xFFFF)
	#view = new Uint8Array(this.#buffer, 0)
	#dview = new DataView(this.#buffer, 0)
	#len = 0;

	#resize(length) {
		if (this.#len < length) {
			this.#len = length;
			if (this.#view.length < length) {
				this.#view = new Uint8Array(this.#buffer, 0, this.#view.length * 2);
				this.#dview = new DataView(this.#view.buffer, 0);
			}
		}
	}
}

const deserialize_bool = (view, offset, struct, field) => {
	struct[field] = view.getUint8(offset) !== 0;
	return offset + 1;
}

const deserialize_int8 = (data, offset, struct, field) => {
	struct[field] = data.getInt8(offset);
	return offset + 1;
}

const deserialize_uint8 = (data, offset, struct, field) => {
	struct[field] = data.getUint8(offset);
	return offset + 1;
}

const deserialize_int16 = (data, offset, struct, field) => {
	struct[field] = data.getInt16(offset);
	return offset + 2;
}

const deserialize_uint16 = (data, offset, struct, field) => {
	struct[field] = data.getUint16(offset);
	return offset + 2;
}

const deserialize_int32 = (data, offset, struct, field) => {
	struct[field] = data.getInt32(offset);
	return offset + 4;
}

const deserialize_uint32 = (data, offset, struct, field) => {
	struct[field] = data.getUint32(offset);
	return offset + 4;
}

const deserialize_string = (data, offset, struct, field) => {
	const length = data.getUint32(offset);
	offset += 4;
	const bytes = new Uint8Array(data.buffer, data.byteOffset + offset, length);
	const decoder = new TextDecoder();
	struct[field] = decoder.decode(bytes);
	return offset + length;
}

const list_deserializer = (item_deserializer) => (data, offset, struct, field) => {
	const length = data.getUint32(offset);
	offset += 4;
	const endOffset = offset + length;
	const list = [];
	let i = 0;
	while (offset < endOffset) {
		offset = item_deserializer(data, offset, list, i);
		i++;
	}
	struct[field] = list;
	return offset;
}

function parse_struct(struct, view, offset) {
	const typeID = view.getUint16(offset);
	offset += 2;
	if (typeID !== struct.constructor.TypeID) {
		throw new Error(`Type Mismatch: Expected ${struct.constructor.TypeID} got ${typeID} for ${struct.constructor.name}`);
	}
	const length = view.getUint32(offset);
	const totalSize = length + 6;
	offset += 4;
	const endOffset = offset + length;
	while (offset < endOffset) {
		const fieldID = view.getUint16(offset);
		offset += 2;
		const next = struct.__deserialize_field(view, fieldID, offset)
		if (next === unknown_field) {
			return totalSize;
		}
		offset = next;
	}
	return offset;
}

function create_static_deserializer(cls) {
	return (view, offset, struct, field) => {
		const s = new cls();
		offset = parse_struct(s, view, offset);
		struct[field] = s;
		return offset;
	}
}

const unknown_field = new Error("Unknown Field")
