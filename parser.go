package xyJson

import (
	"bytes"
	"strconv"
	"strings"
	"unicode/utf8"
)

// parser JSON解析器实现
// parser implements the JSON parser
type parser struct {
	factory  IValueFactory
	maxDepth int
	data     []byte
	pos      int
	line     int
	column   int
	depth    int
	lastChar rune
	lastSize int
}

// NewParser 创建新的JSON解析器
// NewParser creates a new JSON parser
func NewParser() IParser {
	return &parser{
		factory:  NewValueFactory(),
		maxDepth: DefaultMaxDepth,
		line:     1,
		column:   1,
	}
}

// NewParserWithFactory 使用指定工厂创建JSON解析器
// NewParserWithFactory creates a JSON parser with the specified factory
func NewParserWithFactory(factory IValueFactory) IParser {
	if factory == nil {
		factory = NewValueFactory()
	}
	return &parser{
		factory:  factory,
		maxDepth: DefaultMaxDepth,
		line:     1,
		column:   1,
	}
}

// Parse 解析JSON字节数组
// Parse parses JSON byte array
func (p *parser) Parse(data []byte) (IValue, error) {
	if len(data) == 0 {
		return nil, NewInvalidJSONError("empty input", nil)
	}

	p.reset(data)
	p.skipWhitespace()

	if p.pos >= len(p.data) {
		return nil, NewInvalidJSONError("unexpected end of input", nil)
	}

	value, err := p.parseValue()
	if err != nil {
		return nil, err
	}

	// 检查是否还有多余的字符
	p.skipWhitespace()
	if p.pos < len(p.data) {
		return nil, NewInvalidJSONError("unexpected character after JSON", nil)
	}

	return value, nil
}

// ParseString 解析JSON字符串
// ParseString parses JSON string
func (p *parser) ParseString(data string) (IValue, error) {
	return p.Parse([]byte(data))
}

// SetMaxDepth 设置最大解析深度
// SetMaxDepth sets the maximum parsing depth
func (p *parser) SetMaxDepth(depth int) {
	if depth > 0 {
		p.maxDepth = depth
	}
}

// GetMaxDepth 获取最大解析深度
// GetMaxDepth gets the maximum parsing depth
func (p *parser) GetMaxDepth() int {
	return p.maxDepth
}

// reset 重置解析器状态
// reset resets the parser state
func (p *parser) reset(data []byte) {
	p.data = data
	p.pos = 0
	p.line = 1
	p.column = 1
	p.depth = 0
	p.lastChar = 0
	p.lastSize = 0
}

// parseValue 解析JSON值
// parseValue parses a JSON value
func (p *parser) parseValue() (IValue, error) {
	p.skipWhitespace()

	if p.pos >= len(p.data) {
		return nil, NewInvalidJSONError("unexpected end of input", nil)
	}

	ch := p.data[p.pos]
	switch ch {
	case '"':
		return p.parseString()
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case 't', 'f':
		return p.parseBool()
	case 'n':
		return p.parseNull()
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return p.parseNumber()
	default:
		return nil, NewInvalidJSONError("unexpected character: "+string(ch), nil)
	}
}

// parseString 解析字符串
// parseString parses a string
func (p *parser) parseString() (IValue, error) {
	if p.data[p.pos] != '"' {
		return nil, NewInvalidJSONError("expected '\"'", nil)
	}

	p.advance() // 跳过开始的引号
	start := p.pos
	var buf []byte
	hasEscape := false

	for p.pos < len(p.data) {
		ch := p.data[p.pos]
		if ch == '"' {
			// 字符串结束
			var str string
			if hasEscape {
				if buf == nil {
					buf = make([]byte, 0, p.pos-start)
					buf = append(buf, p.data[start:p.pos]...)
				}
				var err error
				str, err = p.unescapeString(string(buf))
				if err != nil {
					return nil, err
				}
			} else {
				str = string(p.data[start:p.pos])
			}
			p.advance() // 跳过结束的引号
			return p.factory.CreateString(str), nil
		}

		if ch == '\\' {
			hasEscape = true
			if buf == nil {
				buf = make([]byte, 0, len(p.data)-start)
				buf = append(buf, p.data[start:p.pos]...)
			}
			// 将转义序列原样添加到缓冲区，稍后由unescapeString处理
			buf = append(buf, ch) // 添加反斜杠
			p.advance() // 跳过反斜杠
			if p.pos >= len(p.data) {
				return nil, NewInvalidJSONError("unexpected end of input in string escape", nil)
			}
			escapeChar := p.data[p.pos]
			buf = append(buf, escapeChar) // 添加转义字符
			
			// 验证转义字符的有效性
			switch escapeChar {
			case '"', '\\', '/', 'b', 'f', 'n', 'r', 't':
				// 有效的转义字符
			case 'u':
				// Unicode转义，需要验证后续4个十六进制字符
				if p.pos+4 >= len(p.data) {
					return nil, NewInvalidJSONError("incomplete unicode escape", nil)
				}
				p.advance() // 跳过'u'
				for i := 0; i < 4; i++ {
					if p.pos >= len(p.data) {
						return nil, NewInvalidJSONError("incomplete unicode escape", nil)
					}
					ch := p.data[p.pos]
					if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
						return nil, NewInvalidJSONError("invalid unicode escape", nil)
					}
					buf = append(buf, ch)
					p.advance()
				}
				continue
			default:
				return nil, NewInvalidJSONError("invalid escape character: \\"+string(escapeChar), nil)
			}
			p.advance()
			continue
		}

		if ch < 0x20 {
			return nil, NewInvalidJSONError("invalid character in string", nil)
		}

		if hasEscape && buf != nil {
			buf = append(buf, ch)
		}
		p.advance()
	}

	return nil, NewInvalidJSONError("unterminated string", nil)
}

// parseUnicodeEscape 解析Unicode转义序列
// parseUnicodeEscape parses Unicode escape sequence
func (p *parser) parseUnicodeEscape() ([]byte, error) {
	// 当前位置应该在'u'字符
	p.advance() // 跳过'u'

	if p.pos+4 > len(p.data) {
		return nil, NewInvalidJSONError("incomplete unicode escape", nil)
	}

	hexStr := string(p.data[p.pos : p.pos+4])
	codePoint, err := strconv.ParseUint(hexStr, 16, 16)
	if err != nil {
		return nil, NewInvalidJSONError("invalid unicode escape: \\u"+hexStr, nil)
	}

	p.pos += 4
	p.column += 4

	// 处理UTF-16代理对
	if 0xD800 <= codePoint && codePoint <= 0xDBFF {
		// 高代理，需要低代理
		if p.pos+6 > len(p.data) || p.data[p.pos] != '\\' || p.data[p.pos+1] != 'u' {
			return nil, NewInvalidJSONError("incomplete surrogate pair", nil)
		}

		p.pos += 2 // 跳过\u
		p.column += 2

		lowHexStr := string(p.data[p.pos : p.pos+4])
		lowCodePoint, err := strconv.ParseUint(lowHexStr, 16, 16)
		if err != nil {
			return nil, NewInvalidJSONError("invalid unicode escape: \\u"+lowHexStr, nil)
		}

		if !(0xDC00 <= lowCodePoint && lowCodePoint <= 0xDFFF) {
			return nil, NewInvalidJSONError("invalid low surrogate", nil)
		}

		p.pos += 4
		p.column += 4

		// 组合代理对
		fullCodePoint := 0x10000 + (codePoint-0xD800)<<10 + (lowCodePoint - 0xDC00)
		buf := make([]byte, 4)
		n := utf8.EncodeRune(buf, rune(fullCodePoint))
		return buf[:n], nil
	}

	// 普通Unicode字符
	buf := make([]byte, 4)
	n := utf8.EncodeRune(buf, rune(codePoint))
	return buf[:n], nil
}

// parseObject 解析对象
// parseObject parses an object
func (p *parser) parseObject() (IValue, error) {
	if p.data[p.pos] != '{' {
		return nil, NewInvalidJSONError("expected '{'", nil)
	}

	p.depth++
	defer func() { p.depth-- }()

	if p.depth > p.maxDepth {
		return nil, NewInvalidJSONError("maximum depth exceeded", nil)
	}

	p.advance() // 跳过 '{'
	p.skipWhitespace()

	obj := p.factory.CreateObject()

	// 空对象
	if p.pos < len(p.data) && p.data[p.pos] == '}' {
		p.advance()
		return obj, nil
	}

	for {
		// 解析键
		p.skipWhitespace()
		if p.pos >= len(p.data) {
			return nil, NewInvalidJSONError("unexpected end of input in object", nil)
		}

		if p.data[p.pos] != '"' {
			return nil, NewInvalidJSONError("expected string key", nil)
		}

		keyValue, err := p.parseString()
		if err != nil {
			return nil, err
		}
		key := keyValue.String()

		// 解析冒号
		p.skipWhitespace()
		if p.pos >= len(p.data) || p.data[p.pos] != ':' {
			return nil, NewInvalidJSONError("expected ':'", nil)
		}
		p.advance() // 跳过 ':'

		// 检查重复键
		if obj.Has(key) {
			return nil, NewInvalidJSONError("duplicate key: "+key, nil)
		}

		// 解析值
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		if err := obj.Set(key, value); err != nil {
			return nil, err
		}

		// 检查下一个字符
		p.skipWhitespace()
		if p.pos >= len(p.data) {
			return nil, NewInvalidJSONError("unexpected end of input in object", nil)
		}

		ch := p.data[p.pos]
		if ch == '}' {
			p.advance()
			break
		} else if ch == ',' {
			p.advance()
			// 继续下一个键值对
		} else {
			return nil, NewInvalidJSONError("expected ',' or '}'", nil)
		}
	}

	return obj, nil
}

// parseArray 解析数组
// parseArray parses an array
func (p *parser) parseArray() (IValue, error) {
	if p.data[p.pos] != '[' {
		return nil, NewInvalidJSONError("expected '['", nil)
	}

	p.depth++
	defer func() { p.depth-- }()

	if p.depth > p.maxDepth {
		return nil, NewInvalidJSONError("maximum depth exceeded", nil)
	}

	p.advance() // 跳过 '['
	p.skipWhitespace()

	arr := p.factory.CreateArray()

	// 空数组
	if p.pos < len(p.data) && p.data[p.pos] == ']' {
		p.advance()
		return arr, nil
	}

	for {
		// 解析值
		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		if err := arr.Append(value); err != nil {
			return nil, err
		}

		// 检查下一个字符
		p.skipWhitespace()
		if p.pos >= len(p.data) {
			return nil, NewInvalidJSONError("unexpected end of input in array", nil)
		}

		ch := p.data[p.pos]
		if ch == ']' {
			p.advance()
			break
		} else if ch == ',' {
			p.advance()
			// 继续下一个元素
		} else {
			return nil, NewInvalidJSONError("expected ',' or ']'", nil)
		}
	}

	return arr, nil
}

// parseBool 解析布尔值
// parseBool parses a boolean value
func (p *parser) parseBool() (IValue, error) {
	if p.pos+4 <= len(p.data) && string(p.data[p.pos:p.pos+4]) == "true" {
		p.pos += 4
		p.column += 4
		return p.factory.CreateBool(true), nil
	}

	if p.pos+5 <= len(p.data) && string(p.data[p.pos:p.pos+5]) == "false" {
		p.pos += 5
		p.column += 5
		return p.factory.CreateBool(false), nil
	}

	return nil, NewInvalidJSONError("invalid boolean value", nil)
}

// parseNull 解析null值
// parseNull parses a null value
func (p *parser) parseNull() (IValue, error) {
	if p.pos+4 <= len(p.data) && string(p.data[p.pos:p.pos+4]) == "null" {
		p.pos += 4
		p.column += 4
		return p.factory.CreateNull(), nil
	}

	return nil, NewInvalidJSONError("invalid null value", nil)
}

// parseNumber 解析数字
// parseNumber parses a number
func (p *parser) parseNumber() (IValue, error) {
	start := p.pos

	// 处理负号
	if p.pos < len(p.data) && p.data[p.pos] == '-' {
		p.advance()
	}

	// 处理整数部分
	if p.pos >= len(p.data) {
		return nil, NewInvalidJSONError("incomplete number", nil)
	}

	if p.data[p.pos] == '0' {
		p.advance()
	} else if p.data[p.pos] >= '1' && p.data[p.pos] <= '9' {
		p.advance()
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.advance()
		}
	} else {
		return nil, NewInvalidJSONError("invalid number", nil)
	}

	isFloat := false

	// 处理小数部分
	if p.pos < len(p.data) && p.data[p.pos] == '.' {
		isFloat = true
		p.advance()
		if p.pos >= len(p.data) || p.data[p.pos] < '0' || p.data[p.pos] > '9' {
			return nil, NewInvalidJSONError("invalid number: missing digits after decimal point", nil)
		}
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.advance()
		}
	}

	// 处理指数部分
	if p.pos < len(p.data) && (p.data[p.pos] == 'e' || p.data[p.pos] == 'E') {
		isFloat = true
		p.advance()
		if p.pos < len(p.data) && (p.data[p.pos] == '+' || p.data[p.pos] == '-') {
			p.advance()
		}
		if p.pos >= len(p.data) || p.data[p.pos] < '0' || p.data[p.pos] > '9' {
			return nil, NewInvalidJSONError("invalid number: missing digits in exponent", nil)
		}
		for p.pos < len(p.data) && p.data[p.pos] >= '0' && p.data[p.pos] <= '9' {
			p.advance()
		}
	}

	numStr := string(p.data[start:p.pos])

	if isFloat {
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return nil, NewInvalidJSONError("invalid number: "+numStr, nil)
		}
		return p.factory.CreateNumber(val), nil
	} else {
		val, err := strconv.ParseInt(numStr, 10, 64)
		if err != nil {
			return nil, NewInvalidJSONError("invalid number: "+numStr, nil)
		}
		return p.factory.CreateNumber(val), nil
	}
}

// skipWhitespace 跳过空白字符
// skipWhitespace skips whitespace characters
func (p *parser) skipWhitespace() {
	for p.pos < len(p.data) {
		ch := p.data[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n' {
			p.advance()
		} else {
			break
		}
	}
}

// advance 推进位置并更新行列信息
// advance advances position and updates line/column information
func (p *parser) advance() {
	if p.pos < len(p.data) {
		ch := p.data[p.pos]
		if ch == '\n' {
			p.line++
			p.column = 1
		} else {
			p.column++
		}
		p.pos++
	}
}

// unescapeString 反转义字符串
// unescapeString unescapes a string
func (p *parser) unescapeString(s string) (string, error) {
	if !strings.Contains(s, "\\") {
		return s, nil
	}

	var buf bytes.Buffer
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '"', '\\', '/':
				buf.WriteByte(s[i+1])
			case 'b':
				buf.WriteByte('\b')
			case 'f':
				buf.WriteByte('\f')
			case 'n':
				buf.WriteByte('\n')
			case 'r':
				buf.WriteByte('\r')
			case 't':
				buf.WriteByte('\t')
			case 'u':
				if i+5 < len(s) {
					hexStr := s[i+2 : i+6]
					codePoint, err := strconv.ParseUint(hexStr, 16, 16)
					if err != nil {
						return "", NewInvalidJSONError("invalid unicode escape: \\u"+hexStr, nil)
					}
					buf.WriteRune(rune(codePoint))
					i += 4 // 额外跳过4个字符
				} else {
					return "", NewInvalidJSONError("incomplete unicode escape", nil)
				}
			default:
				return "", NewInvalidJSONError("invalid escape character", nil)
			}
			i++ // 跳过转义字符
		} else {
			buf.WriteByte(s[i])
		}
	}
	return buf.String(), nil
}
