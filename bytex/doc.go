/*
Package bytex provides zero-copy string and byte slice conversion utilities along with a Bytes type.

The package offers functions for converting between strings and byte slices without memory allocation,
which can significantly improve performance in certain scenarios. However, these functions use unsafe
operations and come with important safety considerations that must be understood before use.

Safety Warning:
The zero-copy conversion functions (StringToBytes and BytesToString) use unsafe operations to avoid
memory allocation. This means the returned slice or string shares memory with the original data.
Modifying the result of StringToBytes will cause undefined behavior since strings in Go are immutable.
Similarly, if the original byte slice passed to BytesToString is modified after the conversion, the
resulting string may also change, leading to undefined behavior. Only use these functions when you
can guarantee the data won't be modified or when you only need read-only access.

Exported Functions:
  - StringToBytes(s string) []byte - Convert string to byte slice without memory allocation
  - BytesToString(b []byte) string - Convert byte slice to string without memory allocation

Bytes Type:
The package also defines a Bytes type with the following methods:

  Construction:
  - FromString(s string) Bytes - Create Bytes from string (zero-copy, see warning below)
  - FromBytes(b []byte) Bytes - Create Bytes from byte slice
  - FromHexString(s string) (Bytes, error) - Create Bytes from hexadecimal string
  - FromBase64String(s string) (Bytes, error) - Create Bytes from Base64 string
  - FromUint64(v uint64, endian binary.ByteOrder) Bytes - Create Bytes from uint64
  - FromInt64(v int64, endian binary.ByteOrder) Bytes - Create Bytes from int64
  - FromUint32(v uint32, endian binary.ByteOrder) Bytes - Create Bytes from uint32
  - FromInt32(v int32, endian binary.ByteOrder) Bytes - Create Bytes from int32
  - FromUint16(v uint16, endian binary.ByteOrder) Bytes - Create Bytes from uint16
  - FromInt16(v int16, endian binary.ByteOrder) Bytes - Create Bytes from int16

  Encoding:
  - (r Bytes) Hex() string - Convert to lowercase hexadecimal string
  - (r Bytes) UpperHex() string - Convert to uppercase hexadecimal string
  - (r Bytes) Base64() string - Convert to Base64 string

  Parsing:
  - (r Bytes) ToUint64(endian binary.ByteOrder) (uint64, error) - Parse bytes as uint64
  - (r Bytes) ToInt64(endian binary.ByteOrder) (int64, error) - Parse bytes as int64
  - (r Bytes) ToUint32(endian binary.ByteOrder) (uint32, error) - Parse bytes as uint32
  - (r Bytes) ToInt32(endian binary.ByteOrder) (int32, error) - Parse bytes as int32
  - (r Bytes) ToUint16(endian binary.ByteOrder) (uint16, error) - Parse bytes as uint16
  - (r Bytes) ToInt16(endian binary.ByteOrder) (int16, error) - Parse bytes as int16

  Comparison & Search:
  - Equal(a, b []byte) bool - Check if two byte slices are equal
  - Compare(a, b []byte) int - Compare two byte slices, returns -1, 0, or 1
  - Contains(b, subslice []byte) bool - Check if b contains subslice
  - Index(s, sep []byte) int - Find first occurrence of sep in s, -1 if not found
  - HasPrefix(s, prefix []byte) bool - Check if s starts with prefix
  - HasSuffix(s, suffix []byte) bool - Check if s ends with suffix

  Serialization:
  - (r Bytes) MarshalJSON() ([]byte, error) - JSON serialization using Base64 encoding
  - (r *Bytes) UnmarshalJSON(data []byte) error - JSON deserialization using Base64 decoding
  - (r Bytes) MarshalText() ([]byte, error) - Text serialization using Hex encoding
  - (r *Bytes) UnmarshalText(data []byte) error - Text deserialization using Hex decoding

  Data Access:
  - (r Bytes) Bytes() []byte - Get underlying byte slice (reference, see warning)
  - (r Bytes) String() string - Get safe printable string with escape sequences
  - (r Bytes) Clone() Bytes - Create a deep copy of the bytes
  - (r Bytes) Len() int - Get the length of the bytes
  - (r Bytes) Open() io.Reader - Create an io.Reader from the data

Reader Type:
The package also includes a Reader type for stream-based reading of Bytes with support for:
  - Position control (Seek, Position, Remaining)
  - Typed reading with endianness support (ReadUint16, ReadUint32, ReadUint64)
  - Convenient methods (Peek, ReadBytes, ReadString)
  - Standard io.Reader interface implementation
  - Protocol processing scenarios

  Methods:
  - NewReader(data Bytes) *Reader - Create a new Reader instance
  - (r *Reader) Position() int - Get current read position
  - (r *Reader) Remaining() int - Get remaining readable bytes
  - (r *Reader) Seek(offset int64, whence int) (int64, error) - Set read position
  - (r *Reader) ReadByte() (byte, error) - Read a single byte
  - (r *Reader) ReadUint16(order binary.ByteOrder) (uint16, error) - Read uint16 with byte order
  - (r *Reader) ReadUint32(order binary.ByteOrder) (uint32, error) - Read uint32 with byte order
  - (r *Reader) ReadUint64(order binary.ByteOrder) (uint64, error) - Read uint64 with byte order
  - (r *Reader) ReadBytes(n int) (Bytes, error) - Read n bytes
  - (r *Reader) ReadString(n int) (string, error) - Read n bytes as string
  - (r *Reader) Peek(n int) (Bytes, error) - Preview next n bytes without moving position
  - (r *Reader) Reset() - Reset position to beginning
  - (r *Reader) Read(p []byte) (n int, err error) - Implement io.Reader interface

Endian Conversion:

Note: For new code, prefer using the standard library's binary.BigEndian
or binary.LittleEndian directly. The custom BigEndian/LittleEndian types
in this package are retained for backward compatibility.

The package provides endian-aware conversion utilities compatible with encoding/binary:
  - BigEndian.BytesToUint64(b []byte) (uint64, error) - Convert big-endian bytes (≤8 bytes) to uint64
  - LittleEndian.BytesToUint64(b []byte) (uint64, error) - Convert little-endian bytes (≤8 bytes) to uint64
  - BigEndian and LittleEndian implement binary.ByteOrder interface

Examples:

	// Zero-copy string to bytes conversion (use with caution!)
	s := "hello world"
	b := bytex.StringToBytes(s)  // Shares memory with s, do not modify b!
	fmt.Printf("%s\n", string(b))    // "hello world"

	// Zero-copy bytes to string conversion (use with caution!)
	b := []byte("hello world")
	s := bytex.BytesToString(b)  // Shares memory with b, ensure b is not modified!
	fmt.Printf("%s\n", s)            // "hello world"

	// Using Bytes type
	result := bytex.FromString("hello")
	hexStr := result.Hex()           // "68656c6c6f"
	base64Str := result.Base64()     // "aGVsbG8="

  // Comparison and search
  data := bytex.FromString("hello world")
  Contains(data, []byte("world"))   // true
  HasPrefix(data, []byte("hello"))  // true

	// Numeric round-trip
	b := bytex.FromUint64(0x1234567890ABCDEF, binary.BigEndian)
	v, _ := b.ToUint64(binary.BigEndian)  // 0x1234567890ABCDEF

	// JSON serialization (Base64)
	data, _ := json.Marshal(result)  // "aGVsbG8="

	// Text serialization (Hex)
	text, _ := result.MarshalText()  // "68656c6c6f"

	// Using Reader for protocol parsing
	data = bytex.FromBytes([]byte{0x00, 0x00, 0x00, 0x05, 'h', 'e', 'l', 'l', 'o'})
	reader := bytex.NewReader(data)
	length, _ := reader.ReadUint32(binary.BigEndian)  // 5
	message, _ := reader.ReadString(int(length))      // "hello"

Security Considerations:
Because the conversion functions use unsafe operations, they bypass Go's type safety. This can lead
to undefined behavior if the resulting data is modified inappropriately. Only use these functions
when performance is critical and you fully understand the implications of sharing memory between
strings and byte slices.
*/
package bytex
