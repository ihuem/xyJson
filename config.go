package xyJson

import (
	"time"
)

// Config 全局配置
// Config represents global configuration
type Config struct {
	// 解析器配置
	Parser ParserConfig `json:"parser"`

	// 序列化器配置
	Serializer SerializerConfig `json:"serializer"`

	// 对象池配置
	ObjectPool ObjectPoolConfig `json:"object_pool"`

	// 性能监控配置
	Performance PerformanceConfig `json:"performance"`

	// JSONPath配置
	JSONPath JSONPathConfig `json:"jsonpath"`
}

// ParserConfig 解析器配置
// ParserConfig represents parser configuration
type ParserConfig struct {
	// 最大嵌套深度
	MaxNestingDepth int `json:"max_nesting_depth"`

	// 缓冲区大小
	BufferSize int `json:"buffer_size"`

	// 是否允许注释
	AllowComments bool `json:"allow_comments"`

	// 是否允许尾随逗号
	AllowTrailingCommas bool `json:"allow_trailing_commas"`

	// 是否严格模式
	StrictMode bool `json:"strict_mode"`
}

// SerializerConfig 序列化器配置
// SerializerConfig represents serializer configuration
type SerializerConfig struct {
	// 默认缩进
	DefaultIndent string `json:"default_indent"`

	// 是否转义HTML
	EscapeHTML bool `json:"escape_html"`

	// 是否转义Unicode字符为\u格式
	EscapeUnicode bool `json:"escape_unicode"`

	// 是否排序键
	SortKeys bool `json:"sort_keys"`

	// 最大序列化深度
	MaxDepth int `json:"max_depth"`

	// 是否紧凑模式
	CompactMode bool `json:"compact_mode"`
}

// ObjectPoolConfig 对象池配置
// ObjectPoolConfig represents object pool configuration
type ObjectPoolConfig struct {
	// 是否启用对象池
	Enabled bool `json:"enabled"`

	// 最大池大小
	MaxPoolSize int `json:"max_pool_size"`

	// 清理间隔
	CleanupInterval time.Duration `json:"cleanup_interval"`

	// 最大空闲时间
	MaxIdleTime time.Duration `json:"max_idle_time"`

	// 预分配大小
	PreallocateSize int `json:"preallocate_size"`
}

// PerformanceConfig 性能监控配置
// PerformanceConfig represents performance monitoring configuration
type PerformanceConfig struct {
	// 是否启用监控
	Enabled bool `json:"enabled"`

	// 最大快照数量
	MaxSnapshots int `json:"max_snapshots"`

	// 快照间隔
	SnapshotInterval time.Duration `json:"snapshot_interval"`

	// 是否启用内存分析
	EnableMemoryProfiling bool `json:"enable_memory_profiling"`

	// 是否启用详细统计
	DetailedStats bool `json:"detailed_stats"`
}

// JSONPathConfig JSONPath配置
// JSONPathConfig represents JSONPath configuration
type JSONPathConfig struct {
	// 最大查询深度
	MaxQueryDepth int `json:"max_query_depth"`

	// 是否启用缓存
	EnableCache bool `json:"enable_cache"`

	// 缓存大小
	CacheSize int `json:"cache_size"`

	// 缓存TTL
	CacheTTL time.Duration `json:"cache_ttl"`
}

// DefaultConfig 返回默认配置
// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Parser: ParserConfig{
			MaxNestingDepth:     MaxNestingDepth,
			BufferSize:          DefaultParserBufferSize,
			AllowComments:       false,
			AllowTrailingCommas: false,
			StrictMode:          true,
		},
		Serializer: SerializerConfig{
			DefaultIndent: DefaultIndent,
			EscapeHTML:    true,
			EscapeUnicode: false,
			SortKeys:      false,
			MaxDepth:      MaxNestingDepth,
			CompactMode:   false,
		},
		ObjectPool: ObjectPoolConfig{
			Enabled:         true,
			MaxPoolSize:     1000,
			CleanupInterval: 5 * time.Minute,
			MaxIdleTime:     10 * time.Minute,
			PreallocateSize: 10,
		},
		Performance: PerformanceConfig{
			Enabled:               false,
			MaxSnapshots:          DefaultMaxSnapshots,
			SnapshotInterval:      DefaultSnapshotInterval,
			EnableMemoryProfiling: false,
			DetailedStats:         false,
		},
		JSONPath: JSONPathConfig{
			MaxQueryDepth: MaxNestingDepth,
			EnableCache:   true,
			CacheSize:     100,
			CacheTTL:      30 * time.Minute,
		},
	}
}

// ProductionConfig 返回生产环境配置
// ProductionConfig returns production environment configuration
func ProductionConfig() *Config {
	config := DefaultConfig()

	// 生产环境优化
	config.Parser.StrictMode = true
	config.Serializer.CompactMode = true
	config.Serializer.SortKeys = true
	config.ObjectPool.Enabled = true
	config.ObjectPool.MaxPoolSize = 2000
	config.ObjectPool.PreallocateSize = 50
	config.Performance.Enabled = true
	config.Performance.DetailedStats = false
	config.JSONPath.EnableCache = true
	config.JSONPath.CacheSize = 500

	return config
}

// DevelopmentConfig 返回开发环境配置
// DevelopmentConfig returns development environment configuration
func DevelopmentConfig() *Config {
	config := DefaultConfig()

	// 开发环境优化
	config.Parser.AllowComments = true
	config.Parser.AllowTrailingCommas = true
	config.Parser.StrictMode = false
	config.Serializer.CompactMode = false
	config.Serializer.DefaultIndent = "  "
	config.Performance.Enabled = true
	config.Performance.DetailedStats = true
	config.Performance.EnableMemoryProfiling = true

	return config
}

// HighPerformanceConfig 返回高性能配置
// HighPerformanceConfig returns high performance configuration
func HighPerformanceConfig() *Config {
	config := DefaultConfig()

	// 高性能优化
	config.Parser.BufferSize = DefaultParserBufferSize * 2
	config.Serializer.CompactMode = true
	config.Serializer.EscapeHTML = false
	config.ObjectPool.Enabled = true
	config.ObjectPool.MaxPoolSize = 5000
	config.ObjectPool.PreallocateSize = 100
	config.ObjectPool.CleanupInterval = 1 * time.Minute
	config.Performance.Enabled = false // 关闭监控以获得最佳性能
	config.JSONPath.EnableCache = true
	config.JSONPath.CacheSize = 1000

	return config
}

// 全局配置实例
// Global configuration instance
var globalConfig *Config = DefaultConfig()

// SetGlobalConfig 设置全局配置
// SetGlobalConfig sets the global configuration
func SetGlobalConfig(config *Config) {
	if config == nil {
		return
	}
	globalConfig = config

	// 应用配置到各个组件
	applyConfigToComponents(config)
}

// GetGlobalConfig 获取全局配置
// GetGlobalConfig gets the global configuration
func GetGlobalConfig() *Config {
	return globalConfig
}

// applyConfigToComponents 将配置应用到各个组件
// applyConfigToComponents applies configuration to components
func applyConfigToComponents(config *Config) {
	// 应用性能监控配置
	if config.Performance.Enabled {
		EnablePerformanceMonitoring()
	} else {
		DisablePerformanceMonitoring()
	}

	// 应用内存分析配置
	if config.Performance.EnableMemoryProfiling {
		StartMemoryProfiling()
	} else {
		StopMemoryProfiling()
	}
}

// LoadConfigFromEnvironment 从环境变量加载配置
// LoadConfigFromEnvironment loads configuration from environment variables
func LoadConfigFromEnvironment() *Config {
	config := DefaultConfig()

	// 这里可以添加从环境变量读取配置的逻辑
	// 例如：
	// if env := os.Getenv("XYJSON_PARSER_MAX_DEPTH"); env != "" {
	//     if depth, err := strconv.Atoi(env); err == nil {
	//         config.Parser.MaxNestingDepth = depth
	//     }
	// }

	return config
}

// ValidateConfig 验证配置的有效性
// ValidateConfig validates the configuration
func ValidateConfig(config *Config) error {
	if config == nil {
		return NewInvalidOperationError("config validation", "config cannot be nil")
	}

	// 验证解析器配置
	if config.Parser.MaxNestingDepth <= 0 {
		return NewInvalidOperationError("config validation", "parser max nesting depth must be positive")
	}
	if config.Parser.BufferSize <= 0 {
		return NewInvalidOperationError("config validation", "parser buffer size must be positive")
	}

	// 验证序列化器配置
	if config.Serializer.MaxDepth <= 0 {
		return NewInvalidOperationError("config validation", "serializer max depth must be positive")
	}

	// 验证对象池配置
	if config.ObjectPool.MaxPoolSize < 0 {
		return NewInvalidOperationError("config validation", "object pool max size cannot be negative")
	}
	if config.ObjectPool.CleanupInterval < 0 {
		return NewInvalidOperationError("config validation", "object pool cleanup interval cannot be negative")
	}

	// 验证性能配置
	if config.Performance.MaxSnapshots < 0 {
		return NewInvalidOperationError("config validation", "performance max snapshots cannot be negative")
	}
	if config.Performance.SnapshotInterval < 0 {
		return NewInvalidOperationError("config validation", "performance snapshot interval cannot be negative")
	}

	// 验证JSONPath配置
	if config.JSONPath.MaxQueryDepth <= 0 {
		return NewInvalidOperationError("config validation", "jsonpath max query depth must be positive")
	}
	if config.JSONPath.CacheSize < 0 {
		return NewInvalidOperationError("config validation", "jsonpath cache size cannot be negative")
	}

	return nil
}

// OptimizeForUseCase 根据使用场景优化配置
// OptimizeForUseCase optimizes configuration for specific use cases
func OptimizeForUseCase(useCase string) *Config {
	switch useCase {
	case "web-api":
		return ProductionConfig()
	case "data-processing":
		return HighPerformanceConfig()
	case "development":
		return DevelopmentConfig()
	case "embedded":
		config := DefaultConfig()
		config.ObjectPool.MaxPoolSize = 100
		config.Performance.Enabled = false
		config.JSONPath.EnableCache = false
		return config
	default:
		return DefaultConfig()
	}
}
