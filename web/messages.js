// Auto-generated code for schema: messages v1

export class Item {
	static get TypeID() { return 13286; }
	constructor(props = {}){ Object.assign(this, props) }
	toBytes() { return Item_serialize(this, new ByteBuffer()).bytes(); }
	fromBytes(bytes) {
		if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
		parse_struct(this, Item_deserialize_field, new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength), 0);
		return this;
	}
}
function Item_serialize(it, b) {
	write_uint16(13286, b);
	const start_index = b.length;
	write_uint32(0, b);
	if(typeof it.name !== 'undefined') {
		write_uint16(1, b);
		write_string(it.name, b);
	}
	const end_index = b.length;
	b.set_uint32(start_index, end_index - (start_index + 4))
	return b;
}

function Item_deserialize(view, offset, struct, field){
	const s = new Item();
	offset = parse_struct(s, Item_deserialize_field, view, offset);
	struct[field] = s;
	return offset;
}

function Item_deserialize_field(it, view, fieldID, offset) {
	switch(fieldID) {
		case 1: return deserialize_string(view, offset, it, 'name')
		default:
			return unknown_field;
	}
}

export class Player {
	static get TypeID() { return 49920; }
	constructor(props = {}){ Object.assign(this, props) }
	toBytes() { return Player_serialize(this, new ByteBuffer()).bytes(); }
	fromBytes(bytes) {
		if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
		parse_struct(this, Player_deserialize_field, new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength), 0);
		return this;
	}
}
function Player_serialize(it, b) {
	write_uint16(49920, b);
	const start_index = b.length;
	write_uint32(0, b);
	if(typeof it.name !== 'undefined') {
		write_uint16(11, b);
		write_string(it.name, b);
	}
	if(typeof it.inventory !== 'undefined') {
		write_uint16(12, b);
		lw1(it.inventory, b)
	}
	if(typeof it.foo !== 'undefined') {
		write_uint16(13, b);
		write_string(it.foo, b);
	}
	if(typeof it.dead !== 'undefined') {
		write_uint16(14, b);
		write_bool(it.dead, b);
	}
	if(typeof it.lol !== 'undefined') {
		write_uint16(15, b);
		lw2(it.lol, b)
	}
	if(typeof it.lol2 !== 'undefined') {
		write_uint16(16, b);
		lw3(it.lol2, b)
	}
	if(typeof it.id !== 'undefined') {
		write_uint16(10, b);
		write_uint32(it.id, b);
	}
	const end_index = b.length;
	b.set_uint32(start_index, end_index - (start_index + 4))
	return b;
}

function Player_deserialize(view, offset, struct, field){
	const s = new Player();
	offset = parse_struct(s, Player_deserialize_field, view, offset);
	struct[field] = s;
	return offset;
}

function Player_deserialize_field(it, view, fieldID, offset) {
	switch(fieldID) {
		case 11: return deserialize_string(view, offset, it, 'name')
		case 12: return ld1(view, offset, it, 'inventory')
		case 13: return deserialize_string(view, offset, it, 'foo')
		case 14: return deserialize_bool(view, offset, it, 'dead')
		case 15: return ld2(view, offset, it, 'lol')
		case 16: return ld3(view, offset, it, 'lol2')
		case 10: return deserialize_uint32(view, offset, it, 'id')
		default:
			return unknown_field;
	}
}

export function deserialize(bytes) {
	if (!('buffer' in bytes)) bytes = new Uint8Array(bytes);
	const view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);
	const typeID = view.getUint16(0, true);
	switch(typeID) {
		case 13286: return new Item().fromBytes(bytes);
		case 49920: return new Player().fromBytes(bytes);
		default: throw new Error(`Unknown TypeID: ${typeID}`);
	}
}

const lw1 = make_list_writer(Item_serialize);
const lw2 = make_list_writer(make_list_writer(write_uint32));
const lw3 = make_list_writer(make_list_writer(make_list_writer(Item_serialize)));
const ld1 = make_list_deserializer(Item_deserialize);
const ld2 = make_list_deserializer(make_list_deserializer(deserialize_uint32));
const ld3 = make_list_deserializer(make_list_deserializer(make_list_deserializer(Item_deserialize)));

let tmp;

class ByteBuffer {
	get length() { return this.len; }
	encoder = new TextEncoder();
	buffer = new ArrayBuffer(0xFFF)
	view = new Uint8Array(this.buffer, 0)
	dview = new DataView(this.buffer, 0)
	len = 0;

	write(value) {
		tmp = this.length;
		this.resize(this.len + value.length);
		this.view.set(value, tmp);
	}

	set_uint8(offset, value) {
		this.resize(offset + 1);
		this.dview.setUint8(offset, value, true);
	}

	set_uint16(offset, value) {
		this.resize(offset + 2);
		this.dview.setUint16(offset, value, true);
	}

	set_uint32(offset, value) {
		this.resize(offset + 4);
		this.dview.setUint32(offset, value, true);
	}

	bytes() {
		return new Uint8Array(this.buffer, 0, this.len);
	}

	resize = (length) => {
		if (this.len < length) {
			this.len = length;
			if (this.view.length < length) {
				const newBuffer = new ArrayBuffer(Math.max(this.view.length * 2, length));
				const newView = new Uint8Array(newBuffer);
				newView.set(this.view, 0);
				this.buffer = newBuffer;
				this.view = newView;
				this.dview = new DataView(this.buffer, 0);
			}
		}
	}
}

function write_bool(value, b) {
	tmp = b.length;
	b.resize(b.len + 1);
	b.dview.setUint8(tmp, value ? 1 : 0, true);
}

function write_int8(value, b) {
	tmp = b.length;
	b.resize(b.len + 1);
	b.dview.setInt8(tmp, value, true);
}

function write_uint8(value, b) {
	tmp = b.length;
	b.resize(b.len + 1);
	b.dview.setUint8(tmp, value, true);
}

function write_int16(value, b) {
	tmp = b.length;
	b.resize(b.len + 2);
	b.dview.setInt16(tmp, value, true);
}

function write_uint16(value, b) {
	tmp = b.length;
	b.resize(b.len + 2);
	b.dview.setUint16(tmp, value, true);
}

function write_int32(value, b) {
	tmp = b.length;
	b.resize(b.len + 4);
	b.dview.setInt32(tmp, value, true);
}

function write_uint32(value, b) {
	tmp = b.length;
	b.resize(b.len + 4);
	b.dview.setUint32(tmp, value, true);
}

function write_string(value, b) {
	const stringLength = value.length;
	if (stringLength > 300) {
		const encoded = b.encoder.encode(value);
		b.set_uint32(b.length, encoded.length);
		b.write(encoded);
		return;
	}
	const lengthPos = b.length;
	write_uint32(0, b);
	const start = b.length;
	if (stringLength === 0) {
		return;
	}
	let codePoint;
	for (let i = 0; i < stringLength; i++) {
		// decode UTF-16
		const a = value.charCodeAt(i);
		if (i + 1 === stringLength || a < 0xD800 || a >= 0xDC00) {
			codePoint = a;
		} else {
			const b2 = value.charCodeAt(++i);  // Renamed to avoid shadowing
			codePoint = (a << 10) + b2 + (0x10000 - (0xD800 << 10) - 0xDC00);
		}
		if (codePoint < 0x80) {
			write_uint8(codePoint, b);
		} else {
			if (codePoint < 0x800) {
				write_uint8(((codePoint >> 6) & 0x1F) | 0xC0, b);
			} else {
				if (codePoint < 0x10000) {
					write_uint8(((codePoint >> 12) & 0x0F) | 0xE0, b);
				} else {
					write_uint8(((codePoint >> 18) & 0x07) | 0xF0, b);
					write_uint8(((codePoint >> 12) & 0x3F) | 0x80, b);
				}
				write_uint8(((codePoint >> 6) & 0x3F) | 0x80, b);
			}
			write_uint8((codePoint & 0x3F) | 0x80, b);
		}
	}
	b.set_uint32(lengthPos, b.length - start);
}

function make_list_writer(s) {
	return (value, b) => {
		write_uint32(0, b);
		const sizeIndex = b.len;
		for (const item of value) s(item, b);
		b.set_uint32(sizeIndex - 4, b.len - sizeIndex);
	}
}

function deserialize_bool(view, offset, struct, field) {
	struct[field] = view.getUint8(offset, true) !== 0;
	return offset + 1;
}

function deserialize_int8(data, offset, struct, field) {
	struct[field] = data.getInt8(offset, true);
	return offset + 1;
}

function deserialize_uint8(data, offset, struct, field) {
	struct[field] = data.getUint8(offset, true);
	return offset + 1;
}

function deserialize_int16(data, offset, struct, field) {
	struct[field] = data.getInt16(offset, true);
	return offset + 2;
}

function deserialize_uint16(data, offset, struct, field) {
	struct[field] = data.getUint16(offset, true);
	return offset + 2;
}

function deserialize_int32(data, offset, struct, field) {
	struct[field] = data.getInt32(offset, true);
	return offset + 4;
}

function deserialize_uint32(data, offset, struct, field) {
	struct[field] = data.getUint32(offset, true);
	return offset + 4;
}

const text_decoder = new TextDecoder();
function deserialize_string(data, offset, struct, field) {
	const length = data.getUint32(offset, true);
	offset += 4;
	if (length > 300) {
		const bytes = new Uint8Array(data.buffer, data.byteOffset + offset, length);
		struct[field] = text_decoder.decode(bytes);
		return offset + length;
	} else {
		const end = offset + length;
		let result = "";
		let codePoint;
		while (offset < end) {
			const a = data.getUint8(offset++);
			if (a < 0xC0) {
				codePoint = a;
			} else {
				const b = data.getUint8(offset++);
				if (a < 0xE0) {
					codePoint = ((a & 0x1F) << 6) | (b & 0x3F);
				} else {
					const c = data.getUint8(offset++);
					if (a < 0xF0) {
						codePoint = ((a & 0x0F) << 12) | ((b & 0x3F) << 6) | (c & 0x3F);
					} else {
						const d = data.getUint8(offset++);
						codePoint = ((a & 0x07) << 18) | ((b & 0x3F) << 12) | ((c & 0x3F) << 6) | (d & 0x3F);
					}
				}
			}
			if (codePoint < 0x10000) {
				result += String.fromCharCode(codePoint);
			} else {
				codePoint -= 0x10000;
				result += String.fromCharCode((codePoint >> 10) + 0xD800, (codePoint & ((1 << 10) - 1)) + 0xDC00);
			}
		}
		offset = end;
		struct[field] = result;
		return offset;
	}
}

function make_list_deserializer(item_deserializer) {
	return (data, offset, struct, field) => {
		const length = data.getUint32(offset, true);
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
	};
}

function parse_struct(struct, field_method, view, offset) {
	const typeID = view.getUint16(offset, true);
	offset += 2;
	if (typeID !== struct.constructor.TypeID) {
		throw new Error(`Type Mismatch: Expected ${struct.constructor.TypeID} got ${typeID} for ${struct.constructor.name}`);
	}
	const length = view.getUint32(offset, true);
	const totalSize = length + 6;
	offset += 4;
	const endOffset = offset + length;
	while (offset < endOffset) {
		const fieldID = view.getUint16(offset, true);
		offset += 2;
		const next = field_method(struct, view, fieldID, offset)
		if (next === unknown_field) {
			return totalSize;
		}
		offset = next;
	}
	return offset;
}

const unknown_field = new Error("Unknown Field")

