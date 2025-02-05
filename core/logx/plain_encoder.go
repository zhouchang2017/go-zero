package logx

import (
	"encoding/base64"
	"encoding/json"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"io"
	"math"
	"sync"
	"time"
	"unicode/utf8"
)

const _hex = "0123456789abcdef"

var _plainPool = sync.Pool{New: func() interface{} {
	return &plainEncoder{}
}}

func getPlainEncoder() *plainEncoder {
	return _plainPool.Get().(*plainEncoder)
}

func putPlainEncoder(enc *plainEncoder) {
	if enc.reflectBuf != nil {
		enc.reflectBuf.Free()
	}
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.spaced = false
	enc.openNamespaces = 0
	enc.reflectBuf = nil
	enc.reflectEnc = nil
	_plainPool.Put(enc)
}

var (
	bufferpool = buffer.NewPool()
)

type plainEncoder struct {
	*zapcore.EncoderConfig
	buf            *buffer.Buffer
	spaced         bool // include spaces after colons and commas
	openNamespaces int

	// for encoding generic values by reflection
	reflectBuf *buffer.Buffer
	reflectEnc zapcore.ReflectedEncoder
}

// NewPlainEncoder creates a fast, low-allocation plain encoder. The encoder
func NewPlainEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	if cfg.ConsoleSeparator == "" {
		cfg.ConsoleSeparator = " "
	}
	return _newPlainEncoder(cfg, true)
}

func _newPlainEncoder(cfg zapcore.EncoderConfig, spaced bool) *plainEncoder {
	if cfg.SkipLineEnding {
		cfg.LineEnding = ""
	} else if cfg.LineEnding == "" {
		cfg.LineEnding = zapcore.DefaultLineEnding
	}

	// If no EncoderConfig.NewReflectedEncoder is provided by the user, then use default
	if cfg.NewReflectedEncoder == nil {
		cfg.NewReflectedEncoder = defaultReflectedEncoder
	}

	return &plainEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
		spaced:        spaced,
	}
}

func (enc *plainEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	defer enc.buf.AppendByte(']')
	enc.addKey(key)
	return enc.AppendArray(arr)
}

func (enc *plainEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *plainEncoder) AddBinary(key string, val []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(val))
}

func (enc *plainEncoder) AddByteString(key string, val []byte) {
	enc.addKey(key)
	enc.AppendByteString(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddBool(key string, val bool) {
	enc.addKey(key)
	enc.AppendBool(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddComplex128(key string, val complex128) {
	enc.addKey(key)
	enc.AppendComplex128(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddComplex64(key string, val complex64) {
	enc.addKey(key)
	enc.AppendComplex64(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddDuration(key string, val time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddFloat64(key string, val float64) {
	enc.addKey(key)
	enc.AppendFloat64(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddFloat32(key string, val float32) {
	enc.addKey(key)
	enc.AppendFloat32(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddInt64(key string, val int64) {
	enc.addKey(key)
	enc.AppendInt64(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) resetReflectBuf() {
	if enc.reflectBuf == nil {
		enc.reflectBuf = bufferpool.Get()
		enc.reflectEnc = enc.NewReflectedEncoder(enc.reflectBuf)
	} else {
		enc.reflectBuf.Reset()
	}
}

var nullLiteralBytes = []byte("null")

// Only invoke the standard JSON encoder if there is actually something to
// encode; otherwise write JSON null literal directly.
func (enc *plainEncoder) encodeReflected(obj interface{}) ([]byte, error) {
	if obj == nil {
		return nullLiteralBytes, nil
	}
	enc.resetReflectBuf()
	if err := enc.reflectEnc.Encode(obj); err != nil {
		return nil, err
	}
	enc.reflectBuf.TrimNewline()
	return enc.reflectBuf.Bytes(), nil
}

func (enc *plainEncoder) AddReflected(key string, obj interface{}) error {
	valueBytes, err := enc.encodeReflected(obj)
	if err != nil {
		return err
	}
	enc.addKey(key)
	_, err = enc.buf.Write(valueBytes)
	enc.buf.AppendByte(']')
	return err
}

func (enc *plainEncoder) OpenNamespace(key string) {
	enc.addKey(key)
	enc.buf.AppendByte('{')
	enc.openNamespaces++
}

func (enc *plainEncoder) AddString(key, val string) {
	enc.addKey(key)
	enc.AppendString(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddTime(key string, val time.Time) {
	enc.addKey(key)
	enc.AppendTime(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AddUint64(key string, val uint64) {
	enc.addKey(key)
	enc.AppendUint64(val)
	enc.buf.AppendByte(']')
}

func (enc *plainEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	err := arr.MarshalLogArray(enc)
	enc.buf.AppendByte(']')
	return err
}

func (enc *plainEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	// Close ONLY new openNamespaces that are created during
	// AppendObject().
	old := enc.openNamespaces
	enc.openNamespaces = 0
	enc.addElementSeparator()
	enc.buf.AppendByte('{')
	err := obj.MarshalLogObject(enc)
	enc.buf.AppendByte('}')
	enc.closeOpenNamespaces()
	enc.openNamespaces = old
	return err
}

func (enc *plainEncoder) AppendBool(val bool) {
	enc.addElementSeparator()
	enc.buf.AppendBool(val)
}

func (enc *plainEncoder) AppendByteString(val []byte) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.safeAddByteString(val)
	enc.buf.AppendByte('"')
}

// appendComplex appends the encoded form of the provided complex128 value.
// precision specifies the encoding precision for the real and imaginary
// components of the complex number.
func (enc *plainEncoder) appendComplex(val complex128, precision int) {
	enc.addElementSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(val)), float64(imag(val))
	enc.buf.AppendByte('"')
	// Because we're always in a quoted string, we can use strconv without
	// special-casing NaN and +/-Inf.
	enc.buf.AppendFloat(r, precision)
	// If imaginary part is less than 0, minus (-) sign is added by default
	// by AppendFloat.
	if i >= 0 {
		enc.buf.AppendByte('+')
	}
	enc.buf.AppendFloat(i, precision)
	enc.buf.AppendByte('i')
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AppendDuration(val time.Duration) {
	cur := enc.buf.Len()
	if e := enc.EncodeDuration; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeDuration is a no-op. Fall back to nanoseconds to keep
		// JSON valid.
		enc.AppendInt64(int64(val))
	}
}

func (enc *plainEncoder) AppendInt64(val int64) {
	enc.addElementSeparator()
	enc.buf.AppendInt(val)
}

func (enc *plainEncoder) AppendReflected(val interface{}) error {
	valueBytes, err := enc.encodeReflected(val)
	if err != nil {
		return err
	}
	enc.addElementSeparator()
	_, err = enc.buf.Write(valueBytes)
	return err
}

func (enc *plainEncoder) AppendString(val string) {
	enc.addElementSeparator()
	//enc.buf.AppendByte('"')
	enc.safeAddString(val)
	//enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AppendTimeLayout(time time.Time, layout string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('"')
	enc.buf.AppendTime(time, layout)
	enc.buf.AppendByte('"')
}

func (enc *plainEncoder) AppendTime(val time.Time) {
	cur := enc.buf.Len()
	if e := enc.EncodeTime; e != nil {
		e(val, enc)
	}
	if cur == enc.buf.Len() {
		// User-supplied EncodeTime is a no-op. Fall back to nanos since epoch to keep
		// output JSON valid.
		enc.AppendInt64(val.UnixNano())
	}
}

func (enc *plainEncoder) AppendUint64(val uint64) {
	enc.addElementSeparator()
	enc.buf.AppendUint(val)
}

func (enc *plainEncoder) AddInt(k string, v int)         { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt32(k string, v int32)     { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt16(k string, v int16)     { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddInt8(k string, v int8)       { enc.AddInt64(k, int64(v)) }
func (enc *plainEncoder) AddUint(k string, v uint)       { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint32(k string, v uint32)   { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint16(k string, v uint16)   { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUint8(k string, v uint8)     { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AddUintptr(k string, v uintptr) { enc.AddUint64(k, uint64(v)) }
func (enc *plainEncoder) AppendComplex64(v complex64)    { enc.appendComplex(complex128(v), 32) }
func (enc *plainEncoder) AppendComplex128(v complex128)  { enc.appendComplex(complex128(v), 64) }
func (enc *plainEncoder) AppendFloat64(v float64)        { enc.appendFloat(v, 64) }
func (enc *plainEncoder) AppendFloat32(v float32)        { enc.appendFloat(float64(v), 32) }
func (enc *plainEncoder) AppendInt(v int)                { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt32(v int32)            { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt16(v int16)            { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendInt8(v int8)              { enc.AppendInt64(int64(v)) }
func (enc *plainEncoder) AppendUint(v uint)              { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint32(v uint32)          { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint16(v uint16)          { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUint8(v uint8)            { enc.AppendUint64(uint64(v)) }
func (enc *plainEncoder) AppendUintptr(v uintptr)        { enc.AppendUint64(uint64(v)) }

func (enc *plainEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *plainEncoder) clone() *plainEncoder {
	clone := getPlainEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.spaced = enc.spaced
	clone.openNamespaces = enc.openNamespaces
	clone.buf = bufferpool.Get()
	return clone
}

func (enc *plainEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()
	//final.buf.AppendByte('{')

	if final.TimeKey != "" && final.EncodeTime != nil {
		final.buf.AppendByte('[')
		final.EncodeTime(ent.Time, final)
		final.buf.AppendByte(']')
		final.addElementSeparator()
	}

	if final.LevelKey != "" && final.EncodeLevel != nil {
		final.buf.AppendByte('[')
		final.EncodeLevel(ent.Level, final)
		final.buf.AppendByte(']')
		final.addElementSeparator()
	}

	if ent.Caller.Defined {
		if final.CallerKey != "" {
			final.buf.AppendByte('[')
			final.EncodeCaller(ent.Caller, final)
			final.buf.AppendByte(']')
			final.addElementSeparator()
		}
		if final.FunctionKey != "" {
			final.buf.AppendByte('[')
			final.AppendString(ent.Caller.Function)
			final.buf.AppendByte(']')
			final.addElementSeparator()
		}
	}

	if final.NameKey != "" {
		if ent.LoggerName == "" {
			ent.LoggerName = "-"
		}
		nameEncoder := enc.EncoderConfig.EncodeName

		// if no name encoder provided, fall back to FullNameEncoder for backwards
		// compatibility
		if nameEncoder == nil {
			nameEncoder = zapcore.FullNameEncoder
		}
		final.buf.AppendByte('[')
		nameEncoder(ent.LoggerName, final)
		final.buf.AppendByte(']')
		final.addElementSeparator()
	}

	if enc.buf.Len() > 0 {
		final.addElementSeparator()
		final.buf.Write(enc.buf.Bytes())
	}

	addFields(final, fields)

	if final.MessageKey != "" {
		//final.addKey(enc.MessageKey)
		final.AppendString(ent.Message)
	}

	final.closeOpenNamespaces()
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}
	//final.buf.AppendByte('}')
	final.buf.AppendString(final.LineEnding)

	ret := final.buf
	putPlainEncoder(final)
	return ret, nil
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func (enc *plainEncoder) truncate() {
	enc.buf.Reset()
}

func (enc *plainEncoder) closeOpenNamespaces() {
	for i := 0; i < enc.openNamespaces; i++ {
		enc.buf.AppendByte('}')
	}
	enc.openNamespaces = 0
}

func (enc *plainEncoder) addKey(key string) {
	enc.addElementSeparator()
	enc.buf.AppendByte('[')
	enc.safeAddString(key)
	//enc.buf.AppendByte('"')
	enc.buf.AppendByte(':')
	if enc.spaced {
		enc.buf.AppendByte(' ')
	}
}

func (enc *plainEncoder) addElementSeparator() {
	last := enc.buf.Len() - 1
	if last < 0 {
		return
	}
	switch enc.buf.Bytes()[last] {
	case '{', '[', ':', ',', ' ':
		return
	case ']':
		enc.buf.AppendByte(' ')
	default:
		//enc.buf.AppendByte(' ')
		if enc.spaced {
			enc.buf.AppendByte(' ')
		}
	}
}

func (enc *plainEncoder) appendFloat(val float64, bitSize int) {
	enc.addElementSeparator()
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`"NaN"`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`"+Inf"`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`"-Inf"`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *plainEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *plainEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *plainEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *plainEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func defaultReflectedEncoder(w io.Writer) zapcore.ReflectedEncoder {
	enc := json.NewEncoder(w)
	// For consistency with our custom JSON encoder.
	enc.SetEscapeHTML(false)
	return enc
}
