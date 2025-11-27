const createStaticDeserializer = (cls) => (view, offset, struct, field) => {
	const s = new cls();
	offset = parseStruct(s, view, offset);
	struct[field] = s;
	return offset;
}

export class Foo {
	static get TypeID() { return 2; }
	static __deserialize = createStaticDeserializer(this)

	bar

	__serialize = (buf) => {
		buf.writeuint16(Foo.TypeID)
		const lenIndex = buf.length;
		buf.writeuint32(0)

		if (this.bar !== undefined) {
			buf.writeuint16(0)
			buf.writestring(this.bar);
		}

		const endIndex = buf.length;
		buf.setuint32(lenIndex, endIndex - (lenIndex + 4))
		return buf;
	}

	__deserializeField = (view, fieldID, offset) => {
		switch (fieldID) {
			case 0:
			console.log("Deserializing field 'bar'");
			return deserializestring(view, offset, this, "bar");
			default: return unknownField;
		}
	}
}

export class MyStruct {
	static get TypeID() { return 1; }
	static __deserialize = createStaticDeserializer(this)

	name
	age
	dead
	numList
	foo

	__serialize = (buf) => {
		buf.writeuint16(MyStruct.TypeID)
		const lenIndex = buf.length;
		buf.writeuint32(0)

		if (this.name !== undefined) {
			buf.writeuint16(0)
			buf.writestring(this.name);
		}

		if (this.age !== undefined) {
			buf.writeuint16(1)
			buf.writeuint32(this.age);
		}

		if (this.dead !== undefined) {
			buf.writeuint16(2)
			buf.writebool(this.dead);
		}

		if (this.numList !== undefined) {
			buf.writeuint16(3)
			buf.listWriter(buf.writeuint32)(this.numList);
		}

		if (this.foo !== undefined) {
			buf.writeuint16(4)
			buf.writestruct(this.foo);
		}

		const endIndex = buf.length;
		buf.setuint32(lenIndex, endIndex - (lenIndex + 4))
		return buf;
	}

	__deserializeField = (view, fieldID, offset) => {
		switch (fieldID) {
			case 0: return deserializestring(view, offset, this, "name");
			case 1: return deserializeuint32(view, offset, this, "age");
			case 2: return deserializebool(view, offset, this, "dead");
			case 3: return newListDeserializer(deserializeuint32)(view, offset, this, "numList");
			case 4: return Foo.__deserialize(view, offset, this, "foo");
			default: return unknownField;
		}
	}
}

const unknownField = new Error("Unknown Field")

class ByteBuffer {
	get length() { return this.#len; }

	write = (value) => {
		this.#resize(this.#len + value.length);
		this.#view.set(value, this.#len - value.length);
	}

	writebool = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setUint8(this.#len - 1, value ? 1 : 0);
	}

	writeint8 = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setInt8(this.#len - 1, value);
	}

	writeuint8 = (value) => {
		this.#resize(this.#len + 1);
		this.#dview.setUint8(this.#len - 1, value);
	}

	writeint16 = (value) => {
		this.#resize(this.#len + 2);
		this.#dview.setInt16(this.#len - 2, value);
	}

	writeuint16 = (value) => {
		this.#resize(this.#len + 2);
		this.#dview.setUint16(this.#len - 2, value);
	}

	writeint32 = (value) => {
		this.#resize(this.#len + 4);
		this.#dview.setInt32(this.#len - 4, value);
	}

	writeuint32 = (value) => {
		this.#resize(this.#len + 4);
		this.#dview.setUint32(this.#len - 4, value);
	}

	writestring = (value) => {
		const encoded = this.#encoder.encode(value);
		this.setuint32(this.#len, encoded.length);
		this.write(encoded);
	}

	writestruct = (value) => {
		value.__serialize(this);
	}

	listWriter = (s) => (value) => {
		this.writeuint32(0); // Placeholder for size
		const sizeIndex = this.length;
		for (const item of value) {
			s(item);
		}
		this.setuint32(sizeIndex - 4, this.#len - sizeIndex);
	}

	setuint8(offset, value) {
		this.#resize(offset + 1);
		this.#dview.setUint8(offset, value);
	}

	setuint16(offset, value) {
		this.#resize(offset + 2);
		this.#dview.setUint16(offset, value);
	}

	setuint32(offset, value) {
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

const deserializebool = (view, offset, struct, field) => {
	struct[field] = view.getUint8(offset) !== 0;
	return offset + 1;
}

const deserializeint8 = (data, offset, struct, field) => {
	struct[field] = data.getInt8(offset);
	return offset + 1;
}

const deserializeuint8 = (data, offset, struct, field) => {
	struct[field] = data.getUint8(offset);
	return offset + 1;
}

const deserializeint16 = (data, offset, struct, field) => {
	struct[field] = data.getInt16(offset);
	return offset + 2;
}

const deserializeuint16 = (data, offset, struct, field) => {
	struct[field] = data.getUint16(offset);
	return offset + 2;
}

const deserializeint32 = (data, offset, struct, field) => {
	struct[field] = data.getInt32(offset);
	return offset + 4;
}

const deserializeuint32 = (data, offset, struct, field) => {
	struct[field] = data.getUint32(offset);
	return offset + 4;
}

const deserializestring = (data, offset, struct, field) => {
	const length = data.getUint32(offset);
	offset += 4;
	const bytes = new Uint8Array(data.buffer, data.byteOffset + offset, length);
	const decoder = new TextDecoder();
	struct[field] = decoder.decode(bytes);
	return offset + length;
}

const newListDeserializer = (itemDeserializer) => (data, offset, struct, field) => {
	const length = data.getUint32(offset);
	offset += 4;
	const endOffset = offset + length;
	const list = [];
	let i = 0;
	while (offset < endOffset) {
		offset = itemDeserializer(data, offset, list, i);
		i++;
	}
	struct[field] = list;
	return offset;
}

function parseStruct(struct, view, offset) {
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
		const next = struct.__deserializeField(view, fieldID, offset)
		if (next === unknownField) {
			return totalSize;
		}
		offset = next;
	}
	return offset;
}

const s = new MyStruct();
s.name = "Ben";
s.age = 30;
s.dead = false;
s.numList = [1, 2, 3, 4, 5];
s.foo = new Foo();
s.foo.bar = "Hello, World!";

const buf = new ByteBuffer();
s.__serialize(buf);
const bytes = buf.bytes();

console.log("Serialized Bytes:", bytes);

const dataView = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
const s2 = new MyStruct();
parseStruct(s2, dataView, 0);

console.log("Deserialized Struct:", s2);
