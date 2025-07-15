package xyJson

// ValueType JSON值类型枚举
// ValueType represents the type of a JSON value
type ValueType int

const (
	// NullValueType 空值类型
	// NullValueType represents a null value
	NullValueType ValueType = iota
	// StringValueType 字符串类型
	// StringValueType represents a string value
	StringValueType
	// NumberValueType 数字类型
	// NumberValueType represents a numeric value
	NumberValueType
	// BoolValueType 布尔类型
	// BoolValueType represents a boolean value
	BoolValueType
	// ObjectValueType 对象类型
	// ObjectValueType represents an object value
	ObjectValueType
	// ArrayValueType 数组类型
	// ArrayValueType represents an array value
	ArrayValueType
)

// String 返回值类型的字符串表示
// String returns the string representation of the value type
func (vt ValueType) String() string {
	switch vt {
	case NullValueType:
		return "null"
	case StringValueType:
		return "string"
	case NumberValueType:
		return "number"
	case BoolValueType:
		return "boolean"
	case ObjectValueType:
		return "object"
	case ArrayValueType:
		return "array"
	default:
		return "unknown"
	}
}

// 默认容量常量
// Default capacity constants
const (
	// DefaultMapCapacity 默认Map容量
	// DefaultMapCapacity is the default capacity for maps
	DefaultMapCapacity = 16
	// DefaultArrayCapacity 默认数组容量
	// DefaultArrayCapacity is the default capacity for arrays
	DefaultArrayCapacity = 8
	// DefaultParserBufferSize 默认解析器缓冲区大小
	// DefaultParserBufferSize is the default buffer size for parsers
	DefaultParserBufferSize = 4096
	// MaxNestingDepth 最大嵌套深度
	// MaxNestingDepth is the maximum allowed nesting depth
	MaxNestingDepth = 1000
)

// 路径段类型枚举
// Path segment type enumeration
type SegmentType int

const (
	// PropertySegmentType 属性段类型
	// PropertySegmentType represents a property-based path segment
	PropertySegmentType SegmentType = iota
	// IndexSegmentType 索引段类型
	// IndexSegmentType represents an index-based path segment
	IndexSegmentType
	// FilterSegmentType 过滤段类型
	// FilterSegmentType represents a filter-based path segment
	FilterSegmentType
	// WildcardSegmentType 通配符段类型
	// WildcardSegmentType represents a wildcard path segment
	WildcardSegmentType
)

// 序列化选项常量
// Serialization option constants
const (
	// DefaultIndent 默认缩进字符
	// DefaultIndent is the default indentation string
	DefaultIndent = "  "
	// CompactMode 紧凑模式标识
	// CompactMode indicates compact serialization
	CompactMode = true
	// PrettyMode 美化模式标识
	// PrettyMode indicates pretty-printed serialization
	PrettyMode = false
)

// 性能监控常量
// Performance monitoring constants
const (
	// DefaultMaxSnapshots 默认最大快照数量
	// DefaultMaxSnapshots is the default maximum number of snapshots
	DefaultMaxSnapshots = 100
	// DefaultSnapshotInterval 默认快照间隔（毫秒）
	// DefaultSnapshotInterval is the default snapshot interval in milliseconds
	DefaultSnapshotInterval = 1000
)
