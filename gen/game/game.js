// generated file, do not edit!

export class Player {
	static get TypeID() { 49920; }
	static __deserialize = create_static_deserializer(this)
	__serialize = (b) => {
		b.write_uint16(49920);
		const start_index = b.length;
		b.write_uint32(0);
		if(this.inventory !== 'undefined') {
			b.write_uint16(10);
			b.list_writer(b.write_struct)(this.inventory)
		}
		if(this.idk !== 'undefined') {
			b.write_uint16(11);
			b.write_struct(this.idk);
		}
		if(this.nested !== 'undefined') {
			b.write_uint16(12);
			b.list_writer(b.list_writer(b.write_struct))(this.nested)
		}
		if(this.id !== 'undefined') {
			b.write_uint16(13);
			b.write_uint32(this.id);
		}
		if(this.name !== 'undefined') {
			b.write_uint16(14);
			b.write_string(this.name);
		}
		const end_index = b.length;
		b.setuint32(start_index, end_index - (start_index + 4))
		return b;
	}
}

export class Foo {
	static get TypeID() { 32471; }
	static __deserialize = create_static_deserializer(this)
	__serialize = (b) => {
		b.write_uint16(32471);
		const start_index = b.length;
		b.write_uint32(0);
		if(this.bar !== 'undefined') {
			b.write_uint16(10);
			b.write_int32(this.bar);
		}
		const end_index = b.length;
		b.setuint32(start_index, end_index - (start_index + 4))
		return b;
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

const new_list_deserializer = (item_deserializer) => (data, offset, struct, field) => {
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

