package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	xyJson "github/ihuem/xyJson" // 导入xyJson包
)

// User 用户结构体
type User struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Age     int       `json:"age"`
	Active  bool      `json:"active"`
	Created time.Time `json:"created"`
	Profile Profile   `json:"profile"`
	Tags    []string  `json:"tags"`
}

// Profile 用户档案
type Profile struct {
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
	Location string `json:"location"`
	Website  string `json:"website"`
}

// 模拟数据库
var users = []User{
	{
		ID:      1,
		Name:    "张三",
		Email:   "zhangsan@example.com",
		Age:     28,
		Active:  true,
		Created: time.Now().Add(-time.Hour * 24 * 30),
		Profile: Profile{
			Avatar:   "https://example.com/avatar1.jpg",
			Bio:      "全栈开发工程师",
			Location: "北京",
			Website:  "https://zhangsan.dev",
		},
		Tags: []string{"golang", "react", "docker"},
	},
	{
		ID:      2,
		Name:    "李四",
		Email:   "lisi@example.com",
		Age:     32,
		Active:  true,
		Created: time.Now().Add(-time.Hour * 24 * 60),
		Profile: Profile{
			Avatar:   "https://example.com/avatar2.jpg",
			Bio:      "DevOps工程师",
			Location: "上海",
			Website:  "https://lisi.blog",
		},
		Tags: []string{"kubernetes", "terraform", "aws"},
	},
	{
		ID:      3,
		Name:    "王五",
		Email:   "wangwu@example.com",
		Age:     25,
		Active:  false,
		Created: time.Now().Add(-time.Hour * 24 * 90),
		Profile: Profile{
			Avatar:   "https://example.com/avatar3.jpg",
			Bio:      "前端开发工程师",
			Location: "深圳",
			Website:  "https://wangwu.io",
		},
		Tags: []string{"vue", "typescript", "webpack"},
	},
}

// APIResponse 统一API响应格式
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta 元数据
type Meta struct {
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
}

func main() {
	// 配置xyJson
	setupXyJson()

	// 设置路由
	http.HandleFunc("/users", handleUsers)
	http.HandleFunc("/users/", handleUserByID)
	http.HandleFunc("/users/search", handleUserSearch)
	http.HandleFunc("/users/stats", handleUserStats)
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/metrics", handleMetrics)

	log.Println("服务器启动在 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// setupXyJson 配置xyJson
func setupXyJson() {
	// 使用生产环境配置
	config := xyJson.ProductionConfig()
	xyJson.SetGlobalConfig(config)

	// 启用性能监控
	xyJson.EnablePerformanceMonitoring()

	// 设置对象池
	pool := xyJson.NewObjectPool()
	xyJson.SetDefaultPool(pool)

	log.Println("xyJson配置完成")
}

// handleUsers 处理用户列表请求
func handleUsers(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	switch r.Method {
	case "GET":
		handleGetUsers(w, r)
	case "POST":
		handleCreateUser(w, r)
	default:
		responseError(w, 405, "方法不允许")
	}

	// 记录请求时间
	log.Printf("%s %s - %v", r.Method, r.URL.Path, time.Since(start))
}

// handleGetUsers 获取用户列表
func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	// 解析查询参数
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 分页逻辑
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(users) {
		end = len(users)
	}
	if start > len(users) {
		start = len(users)
	}

	paginatedUsers := users[start:end]

	// 使用xyJson构建响应
	response := buildUsersResponse(paginatedUsers, page, pageSize, len(users))
	responseJSON(w, 200, response)
}

// handleCreateUser 创建用户
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	defer r.Body.Close()

	// 使用xyJson解析请求
	value, err := xyJson.Parse(body)
	if err != nil {
		responseError(w, 400, fmt.Sprintf("JSON解析错误: %v", err))
		return
	}

	// 验证必填字段
	if !xyJson.Exists(value, "$.name") || !xyJson.Exists(value, "$.email") {
		responseError(w, 400, "缺少必填字段: name, email")
		return
	}

	// 提取用户数据
	name := xyJson.MustGet(value, "$.name").String()
	email := xyJson.MustGet(value, "$.email").String()
	age := 0
	if xyJson.Exists(value, "$.age") {
		age = xyJson.MustToInt(xyJson.MustGet(value, "$.age"))
	}

	// 创建新用户
	newUser := User{
		ID:      len(users) + 1,
		Name:    name,
		Email:   email,
		Age:     age,
		Active:  true,
		Created: time.Now(),
		Profile: Profile{},
		Tags:    []string{},
	}

	// 添加到"数据库"
	users = append(users, newUser)

	// 返回创建的用户
	userJSON := buildUserJSON(newUser)
	responseJSON(w, 201, APIResponse{
		Code:    201,
		Message: "用户创建成功",
		Data:    userJSON,
	})
}

// handleUserByID 根据ID处理用户
func handleUserByID(w http.ResponseWriter, r *http.Request) {
	// 提取用户ID
	path := r.URL.Path
	idStr := path[len("/users/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		responseError(w, 400, "无效的用户ID")
		return
	}

	// 查找用户
	var user *User
	for i := range users {
		if users[i].ID == id {
			user = &users[i]
			break
		}
	}

	if user == nil {
		responseError(w, 404, "用户不存在")
		return
	}

	switch r.Method {
	case "GET":
		userJSON := buildUserJSON(*user)
		responseJSON(w, 200, APIResponse{
			Code:    200,
			Message: "获取成功",
			Data:    userJSON,
		})
	case "PUT":
		handleUpdateUser(w, r, user)
	case "DELETE":
		handleDeleteUser(w, r, id)
	default:
		responseError(w, 405, "方法不允许")
	}
}

// handleUserSearch 用户搜索
func handleUserSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		responseError(w, 400, "缺少搜索参数")
		return
	}

	// 构建所有用户的JSON
	allUsersJSON := buildUsersJSON(users)

	// 使用JSONPath搜索
	var results []xyJson.IValue
	var err error

	// 根据不同条件搜索
	if age, parseErr := strconv.Atoi(query); parseErr == nil {
		// 按年龄搜索
		results, err = xyJson.GetAll(allUsersJSON, fmt.Sprintf("$[?(@.age == %d)]", age))
	} else {
		// 按名字或邮箱搜索
		nameResults, _ := xyJson.GetAll(allUsersJSON, fmt.Sprintf("$[?(@.name =~ /%s/i)]", query))
		emailResults, _ := xyJson.GetAll(allUsersJSON, fmt.Sprintf("$[?(@.email =~ /%s/i)]", query))
		results = append(nameResults, emailResults...)
	}

	if err != nil {
		responseError(w, 500, fmt.Sprintf("搜索错误: %v", err))
		return
	}

	responseJSON(w, 200, APIResponse{
		Code:    200,
		Message: fmt.Sprintf("找到 %d 个结果", len(results)),
		Data:    results,
		Meta: &Meta{
			Total:     len(results),
			Timestamp: time.Now(),
		},
	})
}

// handleUserStats 用户统计
func handleUserStats(w http.ResponseWriter, r *http.Request) {
	allUsersJSON := buildUsersJSON(users)

	// 使用JSONPath进行统计
	activeUsers, _ := xyJson.GetAll(allUsersJSON, "$[?(@.active == true)]")
	inactiveUsers, _ := xyJson.GetAll(allUsersJSON, "$[?(@.active == false)]")
	allAges, _ := xyJson.GetAll(allUsersJSON, "$[*].age")

	// 计算平均年龄
	totalAge := 0
	for _, ageValue := range allAges {
		totalAge += xyJson.MustToInt(ageValue)
	}
	avgAge := float64(totalAge) / float64(len(allAges))

	// 构建统计响应
	stats := xyJson.NewBuilder().
		SetInt("total_users", len(users)).
		SetInt("active_users", len(activeUsers)).
		SetInt("inactive_users", len(inactiveUsers)).
		SetFloat64("average_age", avgAge).
		SetTime("generated_at", time.Now()).
		MustBuild()

	responseJSON(w, 200, APIResponse{
		Code:    200,
		Message: "统计信息",
		Data:    stats,
	})
}

// handleHealth 健康检查
func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := xyJson.NewBuilder().
		SetString("status", "healthy").
		SetTime("timestamp", time.Now()).
		SetString("version", xyJson.GetVersion()).
		MustBuild()

	responseJSON(w, 200, health)
}

// handleMetrics 性能指标
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	stats := xyJson.GetPerformanceStats()
	poolStats := xyJson.GetDefaultPool().GetStats()

	metrics := xyJson.NewBuilder().
		BeginObject("performance").
		SetInt("parse_count", int(stats.ParseCount)).
		SetInt("serialize_count", int(stats.SerializeCount)).
		SetString("avg_parse_time", stats.AvgParseTime.String()).
		SetString("avg_serialize_time", stats.AvgSerializeTime.String()).
		End().
		BeginObject("memory_pool").
		SetInt("total_allocated", int(poolStats.TotalAllocated)).
		SetInt("total_reused", int(poolStats.TotalReused)).
		SetInt("current_in_use", int(poolStats.CurrentInUse)).
		SetFloat64("hit_rate", poolStats.PoolHitRate).
		End().
		SetTime("timestamp", time.Now()).
		MustBuild()

	responseJSON(w, 200, metrics)
}

// 辅助函数

// buildUsersResponse 构建用户列表响应
func buildUsersResponse(users []User, page, pageSize, total int) APIResponse {
	usersJSON := buildUsersJSON(users)
	return APIResponse{
		Code:    200,
		Message: "获取成功",
		Data:    usersJSON,
		Meta: &Meta{
			Total:     total,
			Page:      page,
			PageSize:  pageSize,
			Timestamp: time.Now(),
		},
	}
}

// buildUsersJSON 构建用户JSON数组
func buildUsersJSON(users []User) xyJson.IArray {
	arr := xyJson.CreateArray()
	for _, user := range users {
		userJSON := buildUserJSON(user)
		arr.Append(userJSON)
	}
	return arr
}

// buildUserJSON 构建单个用户JSON
func buildUserJSON(user User) xyJson.IObject {
	profile := xyJson.NewBuilder().
		SetString("avatar", user.Profile.Avatar).
		SetString("bio", user.Profile.Bio).
		SetString("location", user.Profile.Location).
		SetString("website", user.Profile.Website).
		MustBuild()

	tags := xyJson.CreateArray()
	for _, tag := range user.Tags {
		tags.Append(xyJson.CreateString(tag))
	}

	return xyJson.NewBuilder().
		SetInt("id", user.ID).
		SetString("name", user.Name).
		SetString("email", user.Email).
		SetInt("age", user.Age).
		SetBool("active", user.Active).
		SetTime("created", user.Created).
		SetValue("profile", profile).
		SetValue("tags", tags).
		MustBuild().(xyJson.IObject)
}

// responseJSON 发送JSON响应
func responseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	var jsonStr string
	var err error

	if value, ok := data.(xyJson.IValue); ok {
		// 使用xyJson序列化
		jsonStr, err = xyJson.SerializeToString(value)
	} else {
		// 使用xyJson构建器
		value, createErr := xyJson.CreateFromRaw(data)
		if createErr != nil {
			http.Error(w, "序列化错误", 500)
			return
		}
		jsonStr, err = xyJson.Pretty(value)
	}

	if err != nil {
		http.Error(w, "序列化错误", 500)
		return
	}

	w.Write([]byte(jsonStr))
}

// responseError 发送错误响应
func responseError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := APIResponse{
		Code:    statusCode,
		Message: message,
	}
	responseJSON(w, statusCode, errorResponse)
}

// handleUpdateUser 更新用户
func handleUpdateUser(w http.ResponseWriter, r *http.Request, user *User) {
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	defer r.Body.Close()

	value, err := xyJson.Parse(body)
	if err != nil {
		responseError(w, 400, fmt.Sprintf("JSON解析错误: %v", err))
		return
	}

	// 更新字段
	if xyJson.Exists(value, "$.name") {
		user.Name = xyJson.MustGet(value, "$.name").String()
	}
	if xyJson.Exists(value, "$.email") {
		user.Email = xyJson.MustGet(value, "$.email").String()
	}
	if xyJson.Exists(value, "$.age") {
		user.Age = xyJson.MustToInt(xyJson.MustGet(value, "$.age"))
	}
	if xyJson.Exists(value, "$.active") {
		user.Active = xyJson.MustToBool(xyJson.MustGet(value, "$.active"))
	}

	userJSON := buildUserJSON(*user)
	responseJSON(w, 200, APIResponse{
		Code:    200,
		Message: "更新成功",
		Data:    userJSON,
	})
}

// handleDeleteUser 删除用户
func handleDeleteUser(w http.ResponseWriter, r *http.Request, userID int) {
	// 从"数据库"中删除
	for i, user := range users {
		if user.ID == userID {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}

	responseJSON(w, 200, APIResponse{
		Code:    200,
		Message: "删除成功",
	})
}
