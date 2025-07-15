package xyJson

import (
	"os"
	"strconv"
	"time"
)

// Config 全局配置
type Config struct {
	Parser      ParserConfig      `json:"parser"`
	Serializer  SerializerConfig  `json:"serializer"`
	ObjectPool  ObjectPoolConfig  `json:"object_pool"`
	Performance PerformanceConfig `json:"performance"`
	JSONPath    JSONPathConfig    `json:"json_path"`
}

// ParserConfig 解析器配置
type ParserConfig struct {
	MaxDepth        int  `json:"max_depth"`         // 最大嵌套深度
	AllowComments   bool `json:"allow_comments"`    // 是否允许注释
	AllowTrailing   bool `json:"allow_trailing"`    // 是否允许尾随逗号
	StrictMode      bool `json:"strict_mode"`       // 严格模式
	ValidateUTF8    bool `json:"validate_utf8"`     // 验证UTF-8编码
	BufferSize      int  `json:"buffer_size"`       // 缓冲区大小
	EnableStreaming bool `json:"enable_streaming"`  // 启用流式解析
}

// SerializerConfig 序列化器配置
type SerializerConfig struct {
	Indent          string `json:"indent"`            // 缩进字符
	EscapeHTML      bool   `json:"escape_html"`       // 转义HTML
	SortKeys        bool   `json:"sort_keys"`         // 排序键
	CompactOutput   bool   `json:"compact_output"`    // 紧凑输出
	BufferSize      int    `json:"buffer_size"`       // 缓冲区大小
	EnableStreaming bool   `json:"enable_streaming"`  // 启用流式序列化
}

// ObjectPoolConfig 对象池配置
type ObjectPoolConfig struct {
	InitialSize     int           `json:"initial_size"`     // 初始大小
	MaxSize         int           `json:"max_size"`         // 最大大小
	CleanupInterval time.Duration `json:"cleanup_interval"` // 清理间隔
	EnableStats     bool          `json:"enable_stats"`     // 启用统计
	AutoResize      bool          `json:"auto_resize"`      // 自动调整大小
	MaxIdleTime     time.Duration `json:"max_idle_time"`    // 最大空闲时间
}

// PerformanceConfig 性能配置
type PerformanceConfig struct {
	EnableMonitoring bool          `json:"enable_monitoring"` // 启用性能监控
	SampleRate       float64       `json:"sample_rate"`       // 采样率
	MetricsInterval  time.Duration `json:"metrics_interval"`  // 指标收集间隔
	EnableMemStats   bool          `json:"enable_mem_stats"`  // 启用内存统计
	EnableCPUProfile bool          `json:"enable_cpu_profile"` // 启用CPU性能分析
	LogSlowOps       bool          `json:"log_slow_ops"`      // 记录慢操作
	SlowOpThreshold  time.Duration `json:"slow_op_threshold"` // 慢操作阈值
}

// JSONPathConfig JSONPath配置
type JSONPathConfig struct {
	CacheSize       int           `json:"cache_size"`       // 缓存大小
	CacheTTL        time.Duration `json:"cache_ttl"`        // 缓存TTL
	MaxComplexity   int           `json:"max_complexity"`   // 最大复杂度
	EnableOptimizer bool          `json:"enable_optimizer"` // 启用优化器
	ParallelQuery   bool          `json:"parallel_query"`   // 并行查询
	MaxWorkers      int           `json:"max_workers"`      // 最大工作协程数
}

// 全局配置实例
var globalConfig *Config

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Parser: ParserConfig{
			MaxDepth:        64,
			AllowComments:   false,
			AllowTrailing:   false,
			StrictMode:      true,
			ValidateUTF8:    true,
			BufferSize:      4096,
			EnableStreaming: false,
		},
		Serializer: SerializerConfig{
			Indent:          "  ",
			EscapeHTML:      true,
			SortKeys:        false,
			CompactOutput:   false,
			BufferSize:      4096,
			EnableStreaming: false,
		},
		ObjectPool: ObjectPoolConfig{
			InitialSize:     10,
			MaxSize:         1000,
			CleanupInterval: 5 * time.Minute,
			EnableStats:     true,
			AutoResize:      true,
			MaxIdleTime:     10 * time.Minute,
		},
		Performance: PerformanceConfig{
			EnableMonitoring: false,
			SampleRate:       1.0,
			MetricsInterval:  time.Minute,
			EnableMemStats:   false,
			EnableCPUProfile: false,
			LogSlowOps:       false,
			SlowOpThreshold:  100 * time.Millisecond,
		},
		JSONPath: JSONPathConfig{
			CacheSize:       100,
			CacheTTL:        time.Hour,
			MaxComplexity:   1000,
			EnableOptimizer: true,
			ParallelQuery:   false,
			MaxWorkers:      4,
		},
	}
}

// ProductionConfig 返回生产环境配置
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.Parser.StrictMode = true
	config.Parser.ValidateUTF8 = true
	config.Serializer.EscapeHTML = true
	config.ObjectPool.MaxSize = 10000
	config.Performance.EnableMonitoring = true
	config.Performance.LogSlowOps = true
	config.JSONPath.EnableOptimizer = true
	config.JSONPath.ParallelQuery = true
	return config
}

// DevelopmentConfig 返回开发环境配置
func DevelopmentConfig() *Config {
	config := DefaultConfig()
	config.Parser.AllowComments = true
	config.Parser.AllowTrailing = true
	config.Parser.StrictMode = false
	config.Serializer.SortKeys = true
	config.Performance.EnableMonitoring = true
	config.Performance.EnableMemStats = true
	config.Performance.LogSlowOps = true
	config.Performance.SlowOpThreshold = 50 * time.Millisecond
	return config
}

// HighPerformanceConfig 返回高性能配置
func HighPerformanceConfig() *Config {
	config := DefaultConfig()
	config.Parser.BufferSize = 8192
	config.Parser.EnableStreaming = true
	config.Serializer.BufferSize = 8192
	config.Serializer.CompactOutput = true
	config.Serializer.EnableStreaming = true
	config.ObjectPool.InitialSize = 100
	config.ObjectPool.MaxSize = 50000
	config.ObjectPool.CleanupInterval = time.Minute
	config.Performance.EnableMonitoring = true
	config.Performance.SampleRate = 0.1
	config.JSONPath.CacheSize = 1000
	config.JSONPath.EnableOptimizer = true
	config.JSONPath.ParallelQuery = true
	config.JSONPath.MaxWorkers = 8
	return config
}

// SetGlobalConfig 设置全局配置
func SetGlobalConfig(config *Config) {
	globalConfig = config
}

// GetGlobalConfig 获取全局配置
func GetGlobalConfig() *Config {
	if globalConfig == nil {
		globalConfig = DefaultConfig()
	}
	return globalConfig
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	// 解析器配置
	if val := os.Getenv("XYJSON_MAX_DEPTH"); val != "" {
		if depth, err := strconv.Atoi(val); err == nil {
			config.Parser.MaxDepth = depth
		}
	}
	if val := os.Getenv("XYJSON_ALLOW_COMMENTS"); val != "" {
		config.Parser.AllowComments = val == "true"
	}
	if val := os.Getenv("XYJSON_STRICT_MODE"); val != "" {
		config.Parser.StrictMode = val == "true"
	}

	// 对象池配置
	if val := os.Getenv("XYJSON_POOL_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.ObjectPool.MaxSize = size
		}
	}

	// 性能配置
	if val := os.Getenv("XYJSON_ENABLE_MONITORING"); val != "" {
		config.Performance.EnableMonitoring = val == "true"
	}
	if val := os.Getenv("XYJSON_SAMPLE_RATE"); val != "" {
		if rate, err := strconv.ParseFloat(val, 64); err == nil {
			config.Performance.SampleRate = rate
		}
	}

	// JSONPath配置
	if val := os.Getenv("XYJSON_CACHE_SIZE"); val != "" {
		if size, err := strconv.Atoi(val); err == nil {
			config.JSONPath.CacheSize = size
		}
	}

	return config
}

// ValidateConfig 验证配置
func ValidateConfig(config *Config) error {
	if config.Parser.MaxDepth <= 0 {
		return NewConfigError("parser.max_depth must be positive")
	}
	if config.Parser.BufferSize <= 0 {
		return NewConfigError("parser.buffer_size must be positive")
	}
	if config.ObjectPool.MaxSize < config.ObjectPool.InitialSize {
		return NewConfigError("object_pool.max_size must be >= initial_size")
	}
	if config.Performance.SampleRate < 0 || config.Performance.SampleRate > 1 {
		return NewConfigError("performance.sample_rate must be between 0 and 1")
	}
	if config.JSONPath.CacheSize < 0 {
		return NewConfigError("json_path.cache_size must be non-negative")
	}
	return nil
}

// OptimizeForUseCase 根据使用场景优化配置
func OptimizeForUseCase(useCase string) *Config {
	switch useCase {
	case "web_api":
		config := DefaultConfig()
		config.Parser.BufferSize = 2048
		config.Serializer.CompactOutput = true
		config.ObjectPool.MaxSize = 5000
		config.Performance.EnableMonitoring = true
		return config

	case "data_processing":
		config := HighPerformanceConfig()
		config.Parser.EnableStreaming = true
		config.Serializer.EnableStreaming = true
		config.JSONPath.ParallelQuery = true
		return config

	case "embedded":
		config := DefaultConfig()
		config.ObjectPool.MaxSize = 100
		config.Performance.EnableMonitoring = false
		config.JSONPath.CacheSize = 10
		return config

	case "testing":
		return DevelopmentConfig()

	default:
		return DefaultConfig()
	}
}

// ConfigError 配置错误
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}

// NewConfigError 创建配置错误
func NewConfigError(message string) *ConfigError {
	return &ConfigError{Message: message}
}
