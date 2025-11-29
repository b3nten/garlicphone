// Auto-generated code for schema: Messages v1

namespace Messages;

interface IMessages<TSelf> where TSelf : IMessages<TSelf>
{
    public byte[] Serialize();
    public TSelf Deserialize(byte[] data);
}
class Item : IMessages<Item>
{
    public string? Name;

    public readonly static ushort TypeId = 13286;

    public static Item CreateFromBytes(byte[] data)
    {
        Item it = new Item();
        using (MemoryStream ms = new MemoryStream(data))
        using (BinaryReader r = new BinaryReader(ms))
        {
            _Item.Deserialize(it, r);
        }
        return it;
    }

    public Item Deserialize(byte[] data)
    {
        using (MemoryStream ms = new MemoryStream(data))
        using (BinaryReader r = new BinaryReader(ms))
        {
            _Item.Deserialize(this, r);
        }
        return this;
    }

    public byte[] Serialize()
    {
        using (MemoryStream ms = new MemoryStream())
        using (BinaryWriter w = new BinaryWriter(ms))
        {
            _Item.Serialize(this, w);
            return ms.ToArray();
        }
    }
}

file class _Item
{
    static public void Serialize(Item it, BinaryWriter w)
    {
        w.Write(Item.TypeId);
        var lengthPos = w.BaseStream.Position;
        w.Write((UInt32)0);
        if (it.Name != null)
        {
            w.Write((ushort)1);
            var bytes = System.Text.Encoding.UTF8.GetBytes(it.Name);
            w.Write((uint)bytes.Length);
            w.Write(bytes);
        }
        var endPos = w.BaseStream.Position;
        w.Seek((int)lengthPos, SeekOrigin.Begin);
        w.Write((UInt32)(endPos - lengthPos - 4));
        w.Seek(0, SeekOrigin.End);
    }
    static public void Deserialize(Item it, BinaryReader r)
    {
        ushort typeId = r.ReadUInt16();
        if (typeId != Item.TypeId)
        {
            throw new Exception($"TypeId mismatch: expected Item.TypeId but got {typeId}");
        }
        uint length = r.ReadUInt32();
        long startPos = r.BaseStream.Position;
        while (r.BaseStream.Position - startPos < length)
        {
            ushort fieldId = r.ReadUInt16();
            switch (fieldId)
            {
                case 1:
                    {
                        uint strLen = r.ReadUInt32();
                        var strBytes = r.ReadBytes((int)strLen);
                        it.Name = System.Text.Encoding.UTF8.GetString(strBytes);
                    }
                    break;
                default:
                    r.BaseStream.Seek(startPos + length, SeekOrigin.Begin);
                    return;
            }
        }
    }
}
class Player : IMessages<Player>
{
    public uint? Id;
    public string? Name;
    public List<Item>? Inventory;
    public string? Foo;
    public bool? Dead;
    public List<List<uint>>? Lol;
    public List<List<List<Item>>>? Lol2;

    public readonly static ushort TypeId = 49920;

    public static Player CreateFromBytes(byte[] data)
    {
        Player it = new Player();
        using (MemoryStream ms = new MemoryStream(data))
        using (BinaryReader r = new BinaryReader(ms))
        {
            _Player.Deserialize(it, r);
        }
        return it;
    }

    public Player Deserialize(byte[] data)
    {
        using (MemoryStream ms = new MemoryStream(data))
        using (BinaryReader r = new BinaryReader(ms))
        {
            _Player.Deserialize(this, r);
        }
        return this;
    }

    public byte[] Serialize()
    {
        using (MemoryStream ms = new MemoryStream())
        using (BinaryWriter w = new BinaryWriter(ms))
        {
            _Player.Serialize(this, w);
            return ms.ToArray();
        }
    }
}

file class _Player
{
    static public void Serialize(Player it, BinaryWriter w)
    {
        w.Write(Player.TypeId);
        var lengthPos = w.BaseStream.Position;
        w.Write((UInt32)0);
        if (it.Id != null)
        {
            w.Write((ushort)10);
            w.Write(it.Id.Value);
        }
        if (it.Name != null)
        {
            w.Write((ushort)11);
            var bytes = System.Text.Encoding.UTF8.GetBytes(it.Name);
            w.Write((uint)bytes.Length);
            w.Write(bytes);
        }
        if (it.Inventory != null)
        {
            w.Write((ushort)12);
            var length0 = w.BaseStream.Position;
            w.Write((uint)0);
            for (int i0 = 0; i0 < it.Inventory.Count; i0++)
            {
                var e0 = it.Inventory[i0];
                _Item.Serialize(e0, w);
            }
            var end0 = w.BaseStream.Position;
            w.Seek((int)length0, SeekOrigin.Begin);
            w.Write((uint)(end0 - length0 - 4));
            w.Seek(0, SeekOrigin.End);
        }
        if (it.Foo != null)
        {
            w.Write((ushort)13);
            var bytes = System.Text.Encoding.UTF8.GetBytes(it.Foo);
            w.Write((uint)bytes.Length);
            w.Write(bytes);
        }
        if (it.Dead != null)
        {
            w.Write((ushort)14);
            w.Write(it.Dead.Value);
        }
        if (it.Lol != null)
        {
            w.Write((ushort)15);
            var length0 = w.BaseStream.Position;
            w.Write((uint)0);
            for (int i0 = 0; i0 < it.Lol.Count; i0++)
            {
                var e0 = it.Lol[i0];
                var length1 = w.BaseStream.Position;
                w.Write((uint)0);
                for (int i1 = 0; i1 < e0.Count; i1++)
                {
                    var e1 = e0[i1];
                    w.Write(e1);
                }
                var end1 = w.BaseStream.Position;
                w.Seek((int)length1, SeekOrigin.Begin);
                w.Write((uint)(end1 - length1 - 4));
                w.Seek(0, SeekOrigin.End);
            }
            var end0 = w.BaseStream.Position;
            w.Seek((int)length0, SeekOrigin.Begin);
            w.Write((uint)(end0 - length0 - 4));
            w.Seek(0, SeekOrigin.End);
        }
        if (it.Lol2 != null)
        {
            w.Write((ushort)16);
            var length0 = w.BaseStream.Position;
            w.Write((uint)0);
            for (int i0 = 0; i0 < it.Lol2.Count; i0++)
            {
                var e0 = it.Lol2[i0];
                var length1 = w.BaseStream.Position;
                w.Write((uint)0);
                for (int i1 = 0; i1 < e0.Count; i1++)
                {
                    var e1 = e0[i1];
                    var length2 = w.BaseStream.Position;
                    w.Write((uint)0);
                    for (int i2 = 0; i2 < e1.Count; i2++)
                    {
                        var e2 = e1[i2];
                        _Item.Serialize(e2, w);
                    }
                    var end2 = w.BaseStream.Position;
                    w.Seek((int)length2, SeekOrigin.Begin);
                    w.Write((uint)(end2 - length2 - 4));
                    w.Seek(0, SeekOrigin.End);
                }
                var end1 = w.BaseStream.Position;
                w.Seek((int)length1, SeekOrigin.Begin);
                w.Write((uint)(end1 - length1 - 4));
                w.Seek(0, SeekOrigin.End);
            }
            var end0 = w.BaseStream.Position;
            w.Seek((int)length0, SeekOrigin.Begin);
            w.Write((uint)(end0 - length0 - 4));
            w.Seek(0, SeekOrigin.End);
        }
        var endPos = w.BaseStream.Position;
        w.Seek((int)lengthPos, SeekOrigin.Begin);
        w.Write((UInt32)(endPos - lengthPos - 4));
        w.Seek(0, SeekOrigin.End);
    }
    static public void Deserialize(Player it, BinaryReader r)
    {
        ushort typeId = r.ReadUInt16();
        if (typeId != Player.TypeId)
        {
            throw new Exception($"TypeId mismatch: expected Player.TypeId but got {typeId}");
        }
        uint length = r.ReadUInt32();
        long startPos = r.BaseStream.Position;
        while (r.BaseStream.Position - startPos < length)
        {
            ushort fieldId = r.ReadUInt16();
            switch (fieldId)
            {
                case 10:
                    it.Id = r.ReadUInt32();
                    break;
                case 11:
                    {
                        uint strLen = r.ReadUInt32();
                        var strBytes = r.ReadBytes((int)strLen);
                        it.Name = System.Text.Encoding.UTF8.GetString(strBytes);
                    }
                    break;
                case 12:
                    {
                        uint listLength0 = r.ReadUInt32();
                        long startPos0 = r.BaseStream.Position;
                        var list0 = new System.Collections.Generic.List<Item>();
                        while (r.BaseStream.Position - startPos0 < listLength0)
                        {
                            Item e0;
                            {
                                Item obj = new();
                                _Item.Deserialize(obj, r);
                                e0 = obj;
                            }
                            list0.Add(e0);
                        }
                        it.Inventory = list0;
                    }
                    break;
                case 13:
                    {
                        uint strLen = r.ReadUInt32();
                        var strBytes = r.ReadBytes((int)strLen);
                        it.Foo = System.Text.Encoding.UTF8.GetString(strBytes);
                    }
                    break;
                case 14:
                    it.Dead = r.ReadBoolean();
                    break;
                case 15:
                    {
                        uint listLength0 = r.ReadUInt32();
                        long startPos0 = r.BaseStream.Position;
                        var list0 = new System.Collections.Generic.List<List<uint>>();
                        while (r.BaseStream.Position - startPos0 < listLength0)
                        {
                            List<uint> e0;
                            {
                                uint listLength1 = r.ReadUInt32();
                                long startPos1 = r.BaseStream.Position;
                                var list1 = new System.Collections.Generic.List<uint>();
                                while (r.BaseStream.Position - startPos1 < listLength1)
                                {
                                    uint e1;
                                    e1 = r.ReadUInt32();
                                    list1.Add(e1);
                                }
                                e0 = list1;
                            }
                            list0.Add(e0);
                        }
                        it.Lol = list0;
                    }
                    break;
                case 16:
                    {
                        uint listLength0 = r.ReadUInt32();
                        long startPos0 = r.BaseStream.Position;
                        var list0 = new System.Collections.Generic.List<List<List<Item>>>();
                        while (r.BaseStream.Position - startPos0 < listLength0)
                        {
                            List<List<Item>> e0;
                            {
                                uint listLength1 = r.ReadUInt32();
                                long startPos1 = r.BaseStream.Position;
                                var list1 = new System.Collections.Generic.List<List<Item>>();
                                while (r.BaseStream.Position - startPos1 < listLength1)
                                {
                                    List<Item> e1;
                                    {
                                        uint listLength2 = r.ReadUInt32();
                                        long startPos2 = r.BaseStream.Position;
                                        var list2 = new System.Collections.Generic.List<Item>();
                                        while (r.BaseStream.Position - startPos2 < listLength2)
                                        {
                                            Item e2;
                                            {
                                                Item obj = new();
                                                _Item.Deserialize(obj, r);
                                                e2 = obj;
                                            }
                                            list2.Add(e2);
                                        }
                                        e1 = list2;
                                    }
                                    list1.Add(e1);
                                }
                                e0 = list1;
                            }
                            list0.Add(e0);
                        }
                        it.Lol2 = list0;
                    }
                    break;
                default:
                    r.BaseStream.Seek(startPos + length, SeekOrigin.Begin);
                    return;
            }
        }
    }
}
