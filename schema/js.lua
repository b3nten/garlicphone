Schema = Schema

local function fmt(s, tab)
	return (s:gsub('($%b{})', function(w) return tab[w:sub(3, -2)] or w end))
end
getmetatable("").__mod = fmt

local function to_pascal_case(str)
	local result = str:gsub("[%w]+", function(word)
		return word:sub(1, 1):upper() .. word:sub(2):lower()
	end)
	return (result:gsub("[^%w]", ""))
end

local desField = "__deserialize_field"
local serField = "__serialize"
local staticDes = "__deserialize"

local function print_list_serializer(type)
	if type.kind == "primitive" then
		return "b.write_${type}" % { type = type.name }
	elseif type.kind == "struct" then
		return "b.write_struct"
	elseif type.kind == "list" then
		return "b.list_writer(${w})" % { w = print_list_serializer(type.of) }
	else
		error("Unknown list element type")
	end
end

local function print_serializer_fn(field)
	if field.type.kind == "primitive" then
		return "b.write_${type}(this.${field});" % { type = field.type.name, field = field.name }
	elseif field.type.kind == "struct" then
		return "b.write_struct(this.${field});" % { field = field.name }
	elseif field.type.kind == "list" then
		return "b.list_writer(${w})(this.${field})" % { w = print_list_serializer(field.type.of), field = field.name }
	else
		error("Unknown field type")
	end
end

local function print_struct_serializer(struct)
	local out = "function ${sname}_${ser}(b) {\n" % {
		sname = to_pascal_case(struct.name),
		ser = serField
	}
	out = out .. "\tb.write_uint16(${typeid});\n" % { typeid = struct.id }
	out = out .. "\tconst start_index = b.length;\n"
	out = out .. "\tb.write_uint32(0);\n"
	for _, field in pairs(struct.fields) do
		out = out .. "\tif(this.${field} !== 'undefined') {\n" % { field = field.name }
		out = out .. "\t\tb.write_uint16(${fieldid});\n" % { fieldid = field.id }
		out = out .. "\t\t" .. print_serializer_fn(field) .. "\n"
		out = out .. "\t}\n"
	end
	out = out .. "\tconst end_index = b.length;\n"
	out = out .. "\tb.set_uint32(start_index, end_index - (start_index + 4))\n"
	out = out .. "\treturn b;\n"
	out = out .. "}\n"
	out = out ..
	"Object.defineProperty(${sname}.prototype, '__serialize', { value: ${sname}___serialize , enumerable: false })\n" % {
		sname = to_pascal_case(struct.name)
	}
	return out
end

local function print_list_deserializer(type)
	if type.kind == "primitive" then
		return "deserialize_${type}" % { type = type.name }
	elseif type.kind == "struct" then
		return "${stype}.${sdes}" % {
			stype = to_pascal_case(type.name),
			sdes = staticDes
		}
	elseif type.kind == "list" then
		return "list_deserializer(${w})" % { w = print_list_deserializer(type.of) }
	else
		error("Unknown list element type")
	end
end

local function print_deserializer_fn(field)
	if field.type.kind == "primitive" then
		return "deserialize_${type}(view, offset, this, '${field}')" % { type = field.type.name, field = field.name }
	elseif field.type.kind == "struct" then
		return "${stype}.${sdes}(view, offset, this, '${field}')"
				% {
					stype = to_pascal_case(field.type.name),
					field = field.name,
					sdes = staticDes
				}
	elseif field.type.kind == "list" then
		return "list_deserializer(${w})(view, offset, this, '${field}')" % {
			w = print_list_deserializer(field.type.of),
			field = field.name
		}
	else
		error("Unknown field type")
	end
end

local function print_deserializer_switch(struct)
	local out = "function ${sname}_${name}(view, fieldID, offset) {\n" %
	{ name = desField, sname = to_pascal_case(struct.name) }
	out = out .. "\tswitch(fieldID) {\n"
	for _, field in pairs(struct.fields) do
		out = out .. "\t\tcase ${fieldid}: return ${fn}\n" % {
			fieldid = field.id,
			fn = print_deserializer_fn(field)
		}
	end
	out = out .. "\t\tdefault:\n"
	out = out .. "\t\t\treturn unknown_field;\n"
	out = out .. "\t}\n"
	out = out .. "}\n"
	out = out ..
	"Object.defineProperty(${sname}.prototype, '__deserialize_field', { value: ${sname}___deserialize_field , enumerable: false })" %
	{
		sname = to_pascal_case(struct.name)
	}
	return out
end

local function print_static_deserializer(struct)
	local out =
	"Object.defineProperty(${sname}, '__deserialize', { value: create_static_deserializer(${sname}), enumerable: false })\n" %
	{
		sname = to_pascal_case(struct.name)
	}
	return out
end

local function print_struct(struct)
	local out = "export class ${name} {\n" % { name = to_pascal_case(struct.name) }
	out = out .. "\tstatic get TypeID() { return ${typeid}; }\n" % { typeid = struct.id }
	out = out .. "\n\ttoBytes() { return this.__serialize(new ByteBuffer()).bytes(); }"
	out = out .. "\n\tfromBytes(bytes) {"
	out = out .. "\n\t\tif (!('buffer' in bytes)) bytes = new Uint8Array(bytes);"
	out = out .. "\n\t\tparse_struct(this, new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength), 0);\n"
	out = out .. "\t\treturn this;\n"
	out = out .. "\t}\n"
	out = out .. "}\n"
	out = out .. print_static_deserializer(struct)
	out = out .. print_struct_serializer(struct)
	out = out .. print_deserializer_switch(struct)
	return out .. "\n\n"
end

local function print_deserialize(structs)
	local out = ""
	out = out .. "export function deserialize(bytes) {"
	out = out .. "\n\tif (!('buffer' in bytes)) bytes = new Uint8Array(bytes);"
	out = out .. "\n\tconst view = new DataView(bytes.buffer, bytes.byteOffset, bytes.byteLength);"
	out = out .. "\n\tconst typeID = view.getUint16(0);"
	out = out .. "\n\tswitch(typeID) {"
	for _, v in pairs(structs) do
		out = out .. "\n\t\tcase ${typeid}: return new ${sname}().fromBytes(bytes);" % {
			typeid = v.id,
			sname = to_pascal_case(v.name)
		}
	end
	out = out .. "\n\t\tdefault: throw new Error(`Unknown TypeID: ${typeID}`);"
	out = out .. "\n\t}"
	out = out .. "\n}\n"
	return out
end

local function print_ts_type(type)
	if type.kind == "primitive" then
		if type.name == "bool" then
			return "boolean"
		elseif type.name == "string" then
			return "string"
		else
			return "number"
		end
	elseif type.kind == "struct" then
		return to_pascal_case(type.name)
	elseif type.kind == "list" then
		return print_ts_type(type.of) .. "[]"
	else
		error("Unknown type for TS defs")
	end
end

local function print_ts_defs(structs)
	local out = ""
	for _, struct in pairs(structs) do
		out = out .. "export class ${sname} {\n" % { sname = to_pascal_case(struct.name) }
		for _, field in pairs(struct.fields) do
			out = out .. "\t${field}?: ${type};\n" % {
				field = field.name,
				type = print_ts_type(field.type)
			}
		end
		out = out .. "\n\tstatic readonly TypeID: number;\n"
		out = out .. "\ttoBytes(): Uint8Array;\n"
		out = out .. "\tfromBytes(bytes: ArrayBuffer | ArrayBufferView): ${sname};\n" % { sname = to_pascal_case(struct.name) }
		out = out .. "}\n\n"
	end
	out = out .. "export function deserialize(bytes: ArrayBuffer | ArrayBufferView): "
	for _, struct in pairs(structs) do
		out = out .. "${sname} | " % { sname = to_pascal_case(struct.name) }
	end
	out = out:sub(1, -4) .. ";\n\n"
	return out
end

local function print_prelude(name, version)
	local out = "// Auto-generated code for schema: ${name} v${version}\n\n" % {
		name = name,
		version = version
	}
	return out
end

-- CODEGEN STEP

local js_file = print_prelude(Schema.name, Schema.version)
for _, v in pairs(Schema.structs) do
	js_file = js_file .. print_struct(v)
end
js_file = js_file .. print_deserialize(Schema.structs)

local ts_file = print_prelude(Schema.name, Schema.version) .. print_ts_defs(Schema.structs)

-- APPEND INCLDUES AND SET js_file
local include = [[
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
]]

Output = {}

Output[Schema.name .. ".js"] = js_file .. "\n" .. include .. "\n"
Output[Schema.name .. ".d.ts"] = ts_file
