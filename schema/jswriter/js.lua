local function fmt(s, tab)
  return (s:gsub('($%b{})', function(w) return tab[w:sub(3, -2)] or w end))
end
getmetatable("").__mod = fmt

local function to_pascal_case(str)
    local result = str:gsub("[%w]+", function(word)
        return word:sub(1,1):upper() .. word:sub(2):lower()
    end)
    return (result:gsub("[^%w]", ""))
end

local tab = "\t"
local brk = "\n"

local desField = "__deserializeField"
local serField = "__serialize"
local staticDes = "__deserialize"

local function print_list_serializer(field)
	if field.type.kind == "primitive" then
		return "buf.write${type}" % {type = field.type.elementType.name}
	elseif field.type.elementType.kind == "struct" then
		return "buf.writestruct"
	elseif field.type.elementType.kind == "list" then
		return "buf.listwriter(${w})" % {w = print_list_serializer({type = field.type.elementType})}
	else
		error("Unknown list element type")
	end
end

local function print_serializer_fn(field)
	print(field.type.kind)
	if field.type.kind == "primitive" then
		return "buf.write${type}(this.${field});" % {type = field.type.name, field = field.name}
	elseif field.type.kind == "struct" then
		return "buf.writestruct(this.${field});" % {field = field.name}
	elseif field.type.kind == "list" then
		return "buf.listwriter(${w})(this.${field})" % {w = print_list_serializer(field.type.of), field = field.name}
	else
		error("Unknown field type")
	end
end

local function print_str_serializer(struct)
	local out =  "\t${name} = (buf) => {\n" % {name = serField}
	out = out .. "\t\tbuf.writeuint16(${typeid});\n" % {typeid = struct.id}
	out = out .. "\t\tconst startLen = buf.length;\n"
	out = out .. "\t\tbuf.writeuint32(0);\n"
	for _, field in pairs(struct.fields) do
		out = out .. "\t\tif(this.${field} !== 'undefined') {\n" % {field = field.name}
		out = out .. "\t\t\tbuf.writeuint16(${fieldid});\n" % {fieldid = field.id}
		out = out .. "\t\t\t" .. print_serializer_fn(field) .. "\n"
		out = out .. "\t\t}\n"
	end
	return out .. "\t}\n"
end

local function print_struct(struct)
	local out = "export class ${name} {\n" % {name = to_pascal_case(struct.name)}
	out = out .. "\tstatic get TypeID() { ${typeid}; }\n" % {typeid = struct.id}
	out = out .. "\tstatic ${f} = createStaticDeserializer(this)\n" % {f = staticDes}
	out = out .. print_str_serializer(struct)

	return out .. "}\n\n"
end

-- CODEGEN STEP

output = "// generated file, do not edit!\n\n"

for _, v in pairs(structs) do
	output = output .. print_struct(v)
end

-- APPEND INCLDUES AND SET OUTPUT
local include = [[
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

function createStaticDeserializer(cls) {
	return (view, offset, struct, field) => {
		const s = new cls();
		offset = parseStruct(s, view, offset);
		struct[field] = s;
		return offset;
	}
}

const unknownField = new Error("Unknown Field")
]]

output = output .. "\n" .. include .. "\n"
