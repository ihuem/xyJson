package xyJson

import (
	"reflect"
	"strconv"
	"unsafe"
)

// reflectKindToValueType 将reflect.Kind转换为ValueType
func reflectKindToValueType(kind reflect.Kind) ValueType {
	switch kind {
	case reflect.String:
		return StringValueType
	case reflect.Bool:
		return BoolValueType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return NumberValueType
	case reflect.Slice, reflect.Array:
		return ArrayValueType
	case reflect.Map, reflect.Struct:
		return ObjectValueType
	case reflect.Ptr, reflect.Interface:
		return NullValueType // 可能为null
	default:
		return NullValueType
	}
}

// 自定义JSON解析器常量
// Custom JSON parser constants
const (
	// 字符常量
	// Character constants
	CharQuote     = '"'
	CharBackslash = '\\'
	CharSlash     = '/'
	CharBackspace = '\b'
	CharFormfeed  = '\f'
	CharNewline   = '\n'
	CharReturn    = '\r'
	CharTab       = '\t'
	CharSpace     = ' '
	CharComma     = ','
	CharColon     = ':'
	CharLeftBrace = '{'
	CharRightBrace = '}'
	CharLeftBracket = '['
	CharRightBracket = ']'
	
	// 缓冲区大小常量
	// Buffer size constants
	InitialBufferSize = 256
	MaxBufferSize     = 1024 * 1024 // 1MB
)

// ICustomParser 自定义解析器接口
// ICustomParser defines the interface for custom parser
type ICustomParser interface {
	// UnmarshalDirect 直接解析JSON到结构体，避免中间IValue表示
	// UnmarshalDirect parses JSON directly to struct, avoiding intermediate IValue representation
	UnmarshalDirect(data []byte, target interface{}) error
	
	// UnmarshalDirectString 直接解析JSON字符串到结构体
	// UnmarshalDirectString parses JSON string directly to struct
	UnmarshalDirectString(data string, target interface{}) error
}

// customParser 自定义JSON解析器实现
// customParser implements custom JSON parser
type customParser struct {
	// 解析状态
	// Parsing state
	data   []byte
	pos    int
	length int
	
	// 缓存的反射信息
	// Cached reflection info
	structInfoCache map[reflect.Type]*customStructInfo
}

// customStructInfo 自定义结构体信息
// customStructInfo holds custom struct information
type customStructInfo struct {
	Fields map[string]*customFieldInfo
}

// customFieldInfo 自定义字段信息
// customFieldInfo holds custom field information
type customFieldInfo struct {
	Index    int
	Name     string
	Type     reflect.Type
	Kind     reflect.Kind
	IsPtr    bool
	Offset   uintptr
	Setter   func(ptr unsafe.Pointer, value interface{}) error
}

// NewCustomParser 创建新的自定义解析器
// NewCustomParser creates a new custom parser
func NewCustomParser() ICustomParser {
	return &customParser{
		structInfoCache: make(map[reflect.Type]*customStructInfo),
	}
}

// UnmarshalDirect 直接解析JSON到结构体
// UnmarshalDirect parses JSON directly to struct
func (cp *customParser) UnmarshalDirect(data []byte, target interface{}) error {
	if target == nil {
		return NewNullPointerError("target cannot be nil")
	}
	
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr {
		return NewInvalidJSONError("target must be a pointer", nil)
	}
	
	rv = rv.Elem()
	if !rv.CanSet() {
		return NewInvalidJSONError("target must be settable", nil)
	}
	
	cp.reset(data)
	return cp.parseValueDirect(rv)
}

// UnmarshalDirectString 直接解析JSON字符串到结构体
// UnmarshalDirectString parses JSON string directly to struct
func (cp *customParser) UnmarshalDirectString(data string, target interface{}) error {
	return cp.UnmarshalDirect([]byte(data), target)
}

// reset 重置解析器状态
// reset resets parser state
func (cp *customParser) reset(data []byte) {
	cp.data = data
	cp.pos = 0
	cp.length = len(data)
}

// parseValueDirect 直接解析值到reflect.Value
// parseValueDirect parses value directly to reflect.Value
func (cp *customParser) parseValueDirect(rv reflect.Value) error {
	cp.skipWhitespace()
	
	if cp.pos >= cp.length {
		return NewInvalidJSONError("unexpected end of input", nil)
	}
	
	ch := cp.data[cp.pos]
	switch ch {
	case CharQuote:
		return cp.parseStringDirect(rv)
	case CharLeftBrace:
		return cp.parseObjectDirect(rv)
	case CharLeftBracket:
		return cp.parseArrayDirect(rv)
	case 't', 'f':
		return cp.parseBoolDirect(rv)
	case 'n':
		return cp.parseNullDirect(rv)
	default:
		if (ch >= '0' && ch <= '9') || ch == '-' {
			return cp.parseNumberDirect(rv)
		}
		return NewInvalidJSONError("unexpected character", nil)
	}
}

// parseStringDirect 直接解析字符串
// parseStringDirect parses string directly
func (cp *customParser) parseStringDirect(rv reflect.Value) error {
	if cp.data[cp.pos] != CharQuote {
		return NewInvalidJSONError("expected quote", nil)
	}
	
	cp.pos++ // 跳过开始引号
	start := cp.pos
	
	// 快速扫描到结束引号
	for cp.pos < cp.length {
		ch := cp.data[cp.pos]
		if ch == CharQuote {
			// 找到结束引号
			str := string(cp.data[start:cp.pos])
			cp.pos++ // 跳过结束引号
			
			if rv.Kind() == reflect.String {
				rv.SetString(str)
				return nil
			}
			return NewTypeMismatchError(StringValueType, reflectKindToValueType(rv.Type().Kind()), "")
		}
		if ch == CharBackslash {
			// 处理转义字符
			return cp.parseStringWithEscapeDirect(rv, start)
		}
		cp.pos++
	}
	
	return NewInvalidJSONError("unterminated string", nil)
}

// parseStringWithEscapeDirect 解析包含转义字符的字符串
// parseStringWithEscapeDirect parses string with escape characters
func (cp *customParser) parseStringWithEscapeDirect(rv reflect.Value, start int) error {
	// 回退到转义字符位置
	buf := make([]byte, 0, cp.length-start)
	buf = append(buf, cp.data[start:cp.pos]...)
	
	for cp.pos < cp.length {
		ch := cp.data[cp.pos]
		if ch == CharQuote {
			str := string(buf)
			cp.pos++
			if rv.Kind() == reflect.String {
				rv.SetString(str)
				return nil
			}
			return NewTypeMismatchError(StringValueType, reflectKindToValueType(rv.Type().Kind()), "")
		}
		if ch == CharBackslash {
			cp.pos++
			if cp.pos >= cp.length {
				return NewInvalidJSONError("unexpected end in escape", nil)
			}
			escapeChar := cp.data[cp.pos]
			switch escapeChar {
			case CharQuote:
				buf = append(buf, CharQuote)
			case CharBackslash:
				buf = append(buf, CharBackslash)
			case CharSlash:
				buf = append(buf, CharSlash)
			case 'b':
				buf = append(buf, CharBackspace)
			case 'f':
				buf = append(buf, CharFormfeed)
			case 'n':
				buf = append(buf, CharNewline)
			case 'r':
				buf = append(buf, CharReturn)
			case 't':
				buf = append(buf, CharTab)
			default:
				return NewInvalidJSONError("invalid escape character", nil)
			}
		} else {
			buf = append(buf, ch)
		}
		cp.pos++
	}
	
	return NewInvalidJSONError("unterminated string", nil)
}

// parseNumberDirect 直接解析数字
// parseNumberDirect parses number directly
func (cp *customParser) parseNumberDirect(rv reflect.Value) error {
	start := cp.pos
	
	// 扫描数字
	if cp.data[cp.pos] == '-' {
		cp.pos++
	}
	
	if cp.pos >= cp.length || (cp.data[cp.pos] < '0' || cp.data[cp.pos] > '9') {
		return NewInvalidJSONError("invalid number", nil)
	}
	
	// 扫描整数部分
	for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
		cp.pos++
	}
	
	// 检查小数点
	hasDecimal := false
	if cp.pos < cp.length && cp.data[cp.pos] == '.' {
		hasDecimal = true
		cp.pos++
		for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
			cp.pos++
		}
	}
	
	// 检查指数
	if cp.pos < cp.length && (cp.data[cp.pos] == 'e' || cp.data[cp.pos] == 'E') {
		hasDecimal = true
		cp.pos++
		if cp.pos < cp.length && (cp.data[cp.pos] == '+' || cp.data[cp.pos] == '-') {
			cp.pos++
		}
		for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
			cp.pos++
		}
	}
	
	numStr := string(cp.data[start:cp.pos])
	
	// 根据目标类型解析
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if hasDecimal {
			return NewTypeMismatchError(NumberValueType, reflectKindToValueType(rv.Kind()), "decimal number for integer field")
		}
		val, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return NewInvalidJSONError("invalid integer", err)
		}
		rv.SetInt(val)
		return nil
		
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if hasDecimal {
			return NewTypeMismatchError(NumberValueType, reflectKindToValueType(rv.Kind()), "decimal number for unsigned integer field")
		}
		val, err := strconv.ParseUint(numStr, 10, 64)
		if err != nil {
			return NewInvalidJSONError("invalid unsigned integer", err)
		}
		rv.SetUint(val)
		return nil
		
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return NewInvalidJSONError("invalid float", err)
		}
		rv.SetFloat(val)
		return nil
		
	default:
		return NewTypeMismatchError(NumberValueType, reflectKindToValueType(rv.Kind()), "")
	}
}

// parseBoolDirect 直接解析布尔值
// parseBoolDirect parses boolean directly
func (cp *customParser) parseBoolDirect(rv reflect.Value) error {
	if rv.Kind() != reflect.Bool {
		return NewTypeMismatchError(BoolValueType, reflectKindToValueType(rv.Kind()), "")
	}
	
	if cp.pos+4 <= cp.length && string(cp.data[cp.pos:cp.pos+4]) == "true" {
		cp.pos += 4
		rv.SetBool(true)
		return nil
	}
	
	if cp.pos+5 <= cp.length && string(cp.data[cp.pos:cp.pos+5]) == "false" {
		cp.pos += 5
		rv.SetBool(false)
		return nil
	}
	
	return NewInvalidJSONError("invalid boolean value", nil)
}

// parseNullDirect 直接解析null值
// parseNullDirect parses null directly
func (cp *customParser) parseNullDirect(rv reflect.Value) error {
	if cp.pos+4 <= cp.length && string(cp.data[cp.pos:cp.pos+4]) == "null" {
		cp.pos += 4
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}
	return NewInvalidJSONError("invalid null value", nil)
}

// parseObjectDirect 直接解析对象
// parseObjectDirect parses object directly
func (cp *customParser) parseObjectDirect(rv reflect.Value) error {
	if rv.Kind() != reflect.Struct {
		return NewTypeMismatchError(ObjectValueType, reflectKindToValueType(rv.Kind()), "")
	}
	
	if cp.data[cp.pos] != CharLeftBrace {
		return NewInvalidJSONError("expected '{'", nil)
	}
	
	cp.pos++ // 跳过 '{'
	cp.skipWhitespace()
	
	// 检查空对象
	if cp.pos < cp.length && cp.data[cp.pos] == CharRightBrace {
		cp.pos++
		return nil
	}
	
	// 获取结构体信息
	structInfo := cp.getCustomStructInfo(rv.Type())
	
	for {
		// 解析键
		cp.skipWhitespace()
		if cp.pos >= cp.length || cp.data[cp.pos] != CharQuote {
			return NewInvalidJSONError("expected string key", nil)
		}
		
		keyStart := cp.pos + 1
		cp.pos++
		for cp.pos < cp.length && cp.data[cp.pos] != CharQuote {
			if cp.data[cp.pos] == CharBackslash {
				cp.pos++ // 跳过转义字符
			}
			cp.pos++
		}
		
		if cp.pos >= cp.length {
			return NewInvalidJSONError("unterminated string key", nil)
		}
		
		key := string(cp.data[keyStart:cp.pos])
		cp.pos++ // 跳过结束引号
		
		// 跳过冒号
		cp.skipWhitespace()
		if cp.pos >= cp.length || cp.data[cp.pos] != CharColon {
			return NewInvalidJSONError("expected ':'", nil)
		}
		cp.pos++
		
		// 解析值
		if fieldInfo, exists := structInfo.Fields[key]; exists {
			fieldValue := rv.Field(fieldInfo.Index)
			if err := cp.parseValueDirect(fieldValue); err != nil {
				return err
			}
		} else {
			// 跳过未知字段
			if err := cp.skipValue(); err != nil {
				return err
			}
		}
		
		// 检查是否结束
		cp.skipWhitespace()
		if cp.pos >= cp.length {
			return NewInvalidJSONError("unexpected end of object", nil)
		}
		
		if cp.data[cp.pos] == CharRightBrace {
			cp.pos++
			break
		}
		
		if cp.data[cp.pos] != CharComma {
			return NewInvalidJSONError("expected ',' or '}'", nil)
		}
		cp.pos++
	}
	
	return nil
}

// parseArrayDirect 直接解析数组
// parseArrayDirect parses array directly
func (cp *customParser) parseArrayDirect(rv reflect.Value) error {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return NewTypeMismatchError(ArrayValueType, reflectKindToValueType(rv.Kind()), "")
	}
	
	if cp.data[cp.pos] != CharLeftBracket {
		return NewInvalidJSONError("expected '['", nil)
	}
	
	cp.pos++ // 跳过 '['
	cp.skipWhitespace()
	
	// 检查空数组
	if cp.pos < cp.length && cp.data[cp.pos] == CharRightBracket {
		cp.pos++
		if rv.Kind() == reflect.Slice {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}
		return nil
	}
	
	var elements []reflect.Value
	elementType := rv.Type().Elem()
	
	for {
		// 创建新元素
		element := reflect.New(elementType).Elem()
		if err := cp.parseValueDirect(element); err != nil {
			return err
		}
		elements = append(elements, element)
		
		// 检查是否结束
		cp.skipWhitespace()
		if cp.pos >= cp.length {
			return NewInvalidJSONError("unexpected end of array", nil)
		}
		
		if cp.data[cp.pos] == CharRightBracket {
			cp.pos++
			break
		}
		
		if cp.data[cp.pos] != CharComma {
			return NewInvalidJSONError("expected ',' or ']'", nil)
		}
		cp.pos++
		cp.skipWhitespace()
	}
	
	// 设置数组/切片值
	if rv.Kind() == reflect.Slice {
		slice := reflect.MakeSlice(rv.Type(), len(elements), len(elements))
		for i, elem := range elements {
			slice.Index(i).Set(elem)
		}
		rv.Set(slice)
	} else {
		// 数组
		for i, elem := range elements {
			if i >= rv.Len() {
				break
			}
			rv.Index(i).Set(elem)
		}
	}
	
	return nil
}

// skipValue 跳过一个JSON值
// skipValue skips a JSON value
func (cp *customParser) skipValue() error {
	cp.skipWhitespace()
	
	if cp.pos >= cp.length {
		return NewInvalidJSONError("unexpected end of input", nil)
	}
	
	ch := cp.data[cp.pos]
	switch ch {
	case CharQuote:
		return cp.skipString()
	case CharLeftBrace:
		return cp.skipObject()
	case CharLeftBracket:
		return cp.skipArray()
	case 't':
		if cp.pos+4 <= cp.length && string(cp.data[cp.pos:cp.pos+4]) == "true" {
			cp.pos += 4
			return nil
		}
		return NewInvalidJSONError("invalid boolean", nil)
	case 'f':
		if cp.pos+5 <= cp.length && string(cp.data[cp.pos:cp.pos+5]) == "false" {
			cp.pos += 5
			return nil
		}
		return NewInvalidJSONError("invalid boolean", nil)
	case 'n':
		if cp.pos+4 <= cp.length && string(cp.data[cp.pos:cp.pos+4]) == "null" {
			cp.pos += 4
			return nil
		}
		return NewInvalidJSONError("invalid null", nil)
	default:
		if (ch >= '0' && ch <= '9') || ch == '-' {
			return cp.skipNumber()
		}
		return NewInvalidJSONError("unexpected character", nil)
	}
}

// skipString 跳过字符串
// skipString skips a string
func (cp *customParser) skipString() error {
	if cp.data[cp.pos] != CharQuote {
		return NewInvalidJSONError("expected quote", nil)
	}
	
	cp.pos++
	for cp.pos < cp.length {
		ch := cp.data[cp.pos]
		if ch == CharQuote {
			cp.pos++
			return nil
		}
		if ch == CharBackslash {
			cp.pos += 2 // 跳过转义字符
		} else {
			cp.pos++
		}
	}
	return NewInvalidJSONError("unterminated string", nil)
}

// skipObject 跳过对象
// skipObject skips an object
func (cp *customParser) skipObject() error {
	if cp.data[cp.pos] != CharLeftBrace {
		return NewInvalidJSONError("expected '{'", nil)
	}
	
	cp.pos++
	cp.skipWhitespace()
	
	if cp.pos < cp.length && cp.data[cp.pos] == CharRightBrace {
		cp.pos++
		return nil
	}
	
	for {
		// 跳过键
		if err := cp.skipString(); err != nil {
			return err
		}
		
		// 跳过冒号
		cp.skipWhitespace()
		if cp.pos >= cp.length || cp.data[cp.pos] != CharColon {
			return NewInvalidJSONError("expected ':'", nil)
		}
		cp.pos++
		
		// 跳过值
		if err := cp.skipValue(); err != nil {
			return err
		}
		
		// 检查结束
		cp.skipWhitespace()
		if cp.pos >= cp.length {
			return NewInvalidJSONError("unexpected end of object", nil)
		}
		
		if cp.data[cp.pos] == CharRightBrace {
			cp.pos++
			return nil
		}
		
		if cp.data[cp.pos] != CharComma {
			return NewInvalidJSONError("expected ',' or '}'", nil)
		}
		cp.pos++
		cp.skipWhitespace()
	}
}

// skipArray 跳过数组
// skipArray skips an array
func (cp *customParser) skipArray() error {
	if cp.data[cp.pos] != CharLeftBracket {
		return NewInvalidJSONError("expected '['", nil)
	}
	
	cp.pos++
	cp.skipWhitespace()
	
	if cp.pos < cp.length && cp.data[cp.pos] == CharRightBracket {
		cp.pos++
		return nil
	}
	
	for {
		// 跳过值
		if err := cp.skipValue(); err != nil {
			return err
		}
		
		// 检查结束
		cp.skipWhitespace()
		if cp.pos >= cp.length {
			return NewInvalidJSONError("unexpected end of array", nil)
		}
		
		if cp.data[cp.pos] == CharRightBracket {
			cp.pos++
			return nil
		}
		
		if cp.data[cp.pos] != CharComma {
			return NewInvalidJSONError("expected ',' or ']'", nil)
		}
		cp.pos++
		cp.skipWhitespace()
	}
}

// skipNumber 跳过数字
// skipNumber skips a number
func (cp *customParser) skipNumber() error {
	if cp.data[cp.pos] == '-' {
		cp.pos++
	}
	
	if cp.pos >= cp.length || (cp.data[cp.pos] < '0' || cp.data[cp.pos] > '9') {
		return NewInvalidJSONError("invalid number", nil)
	}
	
	// 跳过整数部分
	for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
		cp.pos++
	}
	
	// 跳过小数部分
	if cp.pos < cp.length && cp.data[cp.pos] == '.' {
		cp.pos++
		for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
			cp.pos++
		}
	}
	
	// 跳过指数部分
	if cp.pos < cp.length && (cp.data[cp.pos] == 'e' || cp.data[cp.pos] == 'E') {
		cp.pos++
		if cp.pos < cp.length && (cp.data[cp.pos] == '+' || cp.data[cp.pos] == '-') {
			cp.pos++
		}
		for cp.pos < cp.length && cp.data[cp.pos] >= '0' && cp.data[cp.pos] <= '9' {
			cp.pos++
		}
	}
	
	return nil
}

// skipWhitespace 跳过空白字符
// skipWhitespace skips whitespace characters
func (cp *customParser) skipWhitespace() {
	for cp.pos < cp.length {
		ch := cp.data[cp.pos]
		if ch == CharSpace || ch == CharTab || ch == CharNewline || ch == CharReturn {
			cp.pos++
		} else {
			break
		}
	}
}

// getCustomStructInfo 获取自定义结构体信息
// getCustomStructInfo gets custom struct info
func (cp *customParser) getCustomStructInfo(t reflect.Type) *customStructInfo {
	if info, exists := cp.structInfoCache[t]; exists {
		return info
	}
	
	info := &customStructInfo{
		Fields: make(map[string]*customFieldInfo),
	}
	
	// 遍历所有字段
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		
		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}
		
		// 解析JSON标签
		tag := parseJSONTag(field.Tag.Get("json"))
		
		// 跳过标记为忽略的字段
		if tag.Skip {
			continue
		}
		
		// 确定字段名
		fieldName := field.Name
		if tag.Name != "" {
			fieldName = tag.Name
		}
		
		// 检查是否为指针类型
		isPtr := field.Type.Kind() == reflect.Ptr
		fieldType := field.Type
		if isPtr {
			fieldType = field.Type.Elem()
		}
		
		info.Fields[fieldName] = &customFieldInfo{
			Index:  i,
			Name:   fieldName,
			Type:   fieldType,
			Kind:   fieldType.Kind(),
			IsPtr:  isPtr,
			Offset: field.Offset,
		}
	}
	
	// 缓存结构体信息
	if len(cp.structInfoCache) < StructCacheSize {
		cp.structInfoCache[t] = info
	}
	
	return info
}