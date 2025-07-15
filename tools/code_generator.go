package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// CodeGenerator 代码生成器
type CodeGenerator struct {
	templates map[string]*template.Template
	outputDir string
}

// StructField 结构体字段
type StructField struct {
	Name     string
	Type     string
	JSONTag  string
	Comment  string
	Required bool
}

// StructInfo 结构体信息
type StructInfo struct {
	Name         string
	Package      string
	Fields       []StructField
	Comment      string
	GenerateAPI  bool
	GenerateTest bool
}

// APIEndpoint API端点信息
type APIEndpoint struct {
	Method      string
	Path        string
	Handler     string
	Description string
	Request     string
	Response    string
}

// ProjectInfo 项目信息
type ProjectInfo struct {
	Name        string
	Module      string
	Description string
	Author      string
	Version     string
	License     string
}

func main1() {
	generator := NewCodeGenerator("./generated")

	// 示例：生成用户管理相关代码
	userStruct := StructInfo{
		Name:         "User",
		Package:      "models",
		Comment:      "User 用户信息",
		GenerateAPI:  true,
		GenerateTest: true,
		Fields: []StructField{
			{Name: "ID", Type: "int64", JSONTag: "id", Comment: "用户ID", Required: true},
			{Name: "Username", Type: "string", JSONTag: "username", Comment: "用户名", Required: true},
			{Name: "Email", Type: "string", JSONTag: "email", Comment: "邮箱", Required: true},
			{Name: "Password", Type: "string", JSONTag: "password,omitempty", Comment: "密码", Required: false},
			{Name: "Age", Type: "int", JSONTag: "age", Comment: "年龄", Required: false},
			{Name: "Active", Type: "bool", JSONTag: "active", Comment: "是否激活", Required: false},
			{Name: "CreatedAt", Type: "time.Time", JSONTag: "created_at", Comment: "创建时间", Required: false},
			{Name: "UpdatedAt", Type: "time.Time", JSONTag: "updated_at", Comment: "更新时间", Required: false},
		},
	}

	projectInfo := ProjectInfo{
		Name:        "user-api",
		Module:      "github.com/example/user-api",
		Description: "用户管理API服务",
		Author:      "开发者",
		Version:     "v1.0.0",
		License:     "MIT",
	}

	// 生成代码
	if err := generator.GenerateStruct(userStruct); err != nil {
		log.Fatalf("生成结构体失败: %v", err)
	}

	if err := generator.GenerateAPI(userStruct); err != nil {
		log.Fatalf("生成API失败: %v", err)
	}

	if err := generator.GenerateTests(userStruct); err != nil {
		log.Fatalf("生成测试失败: %v", err)
	}

	if err := generator.GenerateProject(projectInfo); err != nil {
		log.Fatalf("生成项目文件失败: %v", err)
	}

	log.Println("代码生成完成！")
}

// NewCodeGenerator 创建代码生成器
func NewCodeGenerator(outputDir string) *CodeGenerator {
	generator := &CodeGenerator{
		templates: make(map[string]*template.Template),
		outputDir: outputDir,
	}

	generator.loadTemplates()
	return generator
}

// loadTemplates 加载模板
func (cg *CodeGenerator) loadTemplates() {
	// 结构体模板
	cg.templates["struct"] = template.Must(template.New("struct").Parse(structTemplate))

	// API处理器模板
	cg.templates["api"] = template.Must(template.New("api").Parse(apiTemplate))

	// 测试模板
	cg.templates["test"] = template.Must(template.New("test").Parse(testTemplate))

	// 项目文件模板
	cg.templates["main"] = template.Must(template.New("main").Parse(mainTemplate))
	cg.templates["gomod"] = template.Must(template.New("gomod").Parse(goModTemplate))
	cg.templates["dockerfile"] = template.Must(template.New("dockerfile").Parse(dockerfileTemplate))
	cg.templates["readme"] = template.Must(template.New("readme").Parse(readmeTemplate))
	cg.templates["config"] = template.Must(template.New("config").Parse(configTemplate))
}

// GenerateStruct 生成结构体
func (cg *CodeGenerator) GenerateStruct(info StructInfo) error {
	filename := fmt.Sprintf("%s.go", strings.ToLower(info.Name))
	filepath := filepath.Join(cg.outputDir, "models", filename)

	return cg.generateFile("struct", filepath, info)
}

// GenerateAPI 生成API处理器
func (cg *CodeGenerator) GenerateAPI(info StructInfo) error {
	if !info.GenerateAPI {
		return nil
	}

	filename := fmt.Sprintf("%s_handler.go", strings.ToLower(info.Name))
	filepath := filepath.Join(cg.outputDir, "handlers", filename)

	return cg.generateFile("api", filepath, info)
}

// GenerateTests 生成测试文件
func (cg *CodeGenerator) GenerateTests(info StructInfo) error {
	if !info.GenerateTest {
		return nil
	}

	filename := fmt.Sprintf("%s_test.go", strings.ToLower(info.Name))
	filepath := filepath.Join(cg.outputDir, "tests", filename)

	return cg.generateFile("test", filepath, info)
}

// GenerateProject 生成项目文件
func (cg *CodeGenerator) GenerateProject(info ProjectInfo) error {
	// 生成main.go
	if err := cg.generateFile("main", filepath.Join(cg.outputDir, "main.go"), info); err != nil {
		return err
	}

	// 生成go.mod
	if err := cg.generateFile("gomod", filepath.Join(cg.outputDir, "go.mod"), info); err != nil {
		return err
	}

	// 生成Dockerfile
	if err := cg.generateFile("dockerfile", filepath.Join(cg.outputDir, "Dockerfile"), info); err != nil {
		return err
	}

	// 生成README.md
	if err := cg.generateFile("readme", filepath.Join(cg.outputDir, "README.md"), info); err != nil {
		return err
	}

	// 生成配置文件
	if err := cg.generateFile("config", filepath.Join(cg.outputDir, "config", "config.go"), info); err != nil {
		return err
	}

	return nil
}

// generateFile 生成文件
func (cg *CodeGenerator) generateFile(templateName, filepath string, data interface{}) error {
	// 创建目录
	dir := filepath[:strings.LastIndex(filepath, string(os.PathSeparator))]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成内容
	templ, exists := cg.templates[templateName]
	if !exists {
		return fmt.Errorf("模板不存在: %s", templateName)
	}

	var content strings.Builder
	if err := templ.Execute(&content, data); err != nil {
		return fmt.Errorf("执行模板失败: %v", err)
	}

	// 格式化Go代码
	var finalContent []byte
	if strings.HasSuffix(filepath, ".go") {
		formatted, err := format.Source([]byte(content.String()))
		if err != nil {
			log.Printf("格式化代码失败: %v", err)
			finalContent = []byte(content.String())
		} else {
			finalContent = formatted
		}
	} else {
		finalContent = []byte(content.String())
	}

	// 写入文件
	if err := os.WriteFile(filepath, finalContent, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	log.Printf("生成文件: %s", filepath)
	return nil
}

// 模板定义

const structTemplate = `package {{.Package}}

import (
	"time"
	"github.com/your-org/xyJson"
)

// {{.Comment}}
type {{.Name}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSONTag}}\"`" + ` // {{.Comment}}
{{end}}}

// New{{.Name}} 创建新的{{.Name}}
func New{{.Name}}() *{{.Name}} {
	return &{{.Name}}{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ToJSON 转换为JSON
func ({{.Name | ToLower}} *{{.Name}}) ToJSON() (xyJson.IValue, error) {
	return xyJson.CreateFromRaw({{.Name | ToLower}})
}

// FromJSON 从JSON创建
func ({{.Name | ToLower}} *{{.Name}}) FromJSON(value xyJson.IValue) error {
	if !value.IsObject() {
		return fmt.Errorf("期望对象类型")
	}

	obj := value.(xyJson.IObject)
{{range .Fields}}{{if .Required}}	if !obj.Has("{{.JSONTag}}") {
		return fmt.Errorf("缺少必填字段: {{.JSONTag}}")
	}
{{end}}{{end}}

{{range .Fields}}	if obj.Has("{{.JSONTag}}") {
		val, _ := obj.Get("{{.JSONTag}}")
{{if eq .Type "string"}}		{{$.Name | ToLower}}.{{.Name}} = val.String()
{{else if eq .Type "int"}}		{{$.Name | ToLower}}.{{.Name}} = xyJson.MustToInt(val)
{{else if eq .Type "int64"}}		{{$.Name | ToLower}}.{{.Name}} = int64(xyJson.MustToInt(val))
{{else if eq .Type "bool"}}		{{$.Name | ToLower}}.{{.Name}} = xyJson.MustToBool(val)
{{else if eq .Type "time.Time"}}		{{$.Name | ToLower}}.{{.Name}} = xyJson.MustToTime(val)
{{else}}		// TODO: 处理{{.Type}}类型
{{end}}	}
{{end}}
	return nil
}

// Validate 验证数据
func ({{.Name | ToLower}} *{{.Name}}) Validate() error {
{{range .Fields}}{{if .Required}}	{{if eq .Type "string"}}if {{$.Name | ToLower}}.{{.Name}} == "" {
		return fmt.Errorf("{{.Name}}不能为空")
	}
{{else if eq .Type "int" "int64"}}if {{$.Name | ToLower}}.{{.Name}} <= 0 {
		return fmt.Errorf("{{.Name}}必须大于0")
	}
{{end}}{{end}}{{end}}	return nil
}

// Update 更新字段
func ({{.Name | ToLower}} *{{.Name}}) Update(updates xyJson.IObject) error {
{{range .Fields}}{{if ne .Name "ID" "CreatedAt"}}	if updates.Has("{{.JSONTag}}") {
		val, _ := updates.Get("{{.JSONTag}}")
{{if eq .Type "string"}}		{{$.Name | ToLower}}.{{.Name}} = val.String()
{{else if eq .Type "int"}}		{{$.Name | ToLower}}.{{.Name}} = xyJson.MustToInt(val)
{{else if eq .Type "int64"}}		{{$.Name | ToLower}}.{{.Name}} = int64(xyJson.MustToInt(val))
{{else if eq .Type "bool"}}		{{$.Name | ToLower}}.{{.Name}} = xyJson.MustToBool(val)
{{end}}	}
{{end}}{{end}}	{{.Name | ToLower}}.UpdatedAt = time.Now()
	return nil
}
`

const apiTemplate = `package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"../models"
	"github.com/your-org/xyJson"
)

// {{.Name}}Handler {{.Name}}处理器
type {{.Name}}Handler struct {
	// 这里可以添加数据库连接等依赖
	data []models.{{.Name}} // 模拟数据存储
}

// New{{.Name}}Handler 创建{{.Name}}处理器
func New{{.Name}}Handler() *{{.Name}}Handler {
	return &{{.Name}}Handler{
		data: make([]models.{{.Name}}, 0),
	}
}

// RegisterRoutes 注册路由
func (h *{{.Name}}Handler) RegisterRoutes() {
	http.HandleFunc("/{{.Name | ToLower}}s", h.handle{{.Name}}s)
	http.HandleFunc("/{{.Name | ToLower}}s/", h.handle{{.Name}}ByID)
}

// handle{{.Name}}s 处理{{.Name}}列表请求
func (h *{{.Name}}Handler) handle{{.Name}}s(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get{{.Name}}s(w, r)
	case "POST":
		h.create{{.Name}}(w, r)
	default:
		h.responseError(w, 405, "方法不允许")
	}
}

// handle{{.Name}}ByID 根据ID处理{{.Name}}
func (h *{{.Name}}Handler) handle{{.Name}}ByID(w http.ResponseWriter, r *http.Request) {
	// 提取ID
	path := r.URL.Path
	idStr := path[len("/{{.Name | ToLower}}s/"):]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.responseError(w, 400, "无效的ID")
		return
	}

	switch r.Method {
	case "GET":
		h.get{{.Name}}ByID(w, r, id)
	case "PUT":
		h.update{{.Name}}(w, r, id)
	case "DELETE":
		h.delete{{.Name}}(w, r, id)
	default:
		h.responseError(w, 405, "方法不允许")
	}
}

// get{{.Name}}s 获取{{.Name}}列表
func (h *{{.Name}}Handler) get{{.Name}}s(w http.ResponseWriter, r *http.Request) {
	// 分页参数
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
	if end > len(h.data) {
		end = len(h.data)
	}
	if start > len(h.data) {
		start = len(h.data)
	}

	paginatedData := h.data[start:end]

	// 构建响应
	response := h.buildListResponse(paginatedData, page, pageSize, len(h.data))
	h.responseJSON(w, 200, response)
}

// create{{.Name}} 创建{{.Name}}
func (h *{{.Name}}Handler) create{{.Name}}(w http.ResponseWriter, r *http.Request) {
	// 读取请求体
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	defer r.Body.Close()

	// 解析JSON
	value, err := xyJson.Parse(body)
	if err != nil {
		h.responseError(w, 400, fmt.Sprintf("JSON解析错误: %v", err))
		return
	}

	// 创建{{.Name}}
	{{.Name | ToLower}} := models.New{{.Name}}()
	if err := {{.Name | ToLower}}.FromJSON(value); err != nil {
		h.responseError(w, 400, fmt.Sprintf("数据转换错误: %v", err))
		return
	}

	// 验证数据
	if err := {{.Name | ToLower}}.Validate(); err != nil {
		h.responseError(w, 400, fmt.Sprintf("数据验证失败: %v", err))
		return
	}

	// 设置ID
	{{.Name | ToLower}}.ID = int64(len(h.data) + 1)

	// 保存到"数据库"
	h.data = append(h.data, *{{.Name | ToLower}})

	// 返回创建的{{.Name}}
	{{.Name | ToLower}}JSON, _ := {{.Name | ToLower}}.ToJSON()
	response := xyJson.NewBuilder().
		SetInt("code", 201).
		SetString("message", "创建成功").
		SetValue("data", {{.Name | ToLower}}JSON).
		MustBuild()

	h.responseJSON(w, 201, response)
}

// get{{.Name}}ByID 根据ID获取{{.Name}}
func (h *{{.Name}}Handler) get{{.Name}}ByID(w http.ResponseWriter, r *http.Request, id int64) {
	// 查找{{.Name}}
	var found *models.{{.Name}}
	for i := range h.data {
		if h.data[i].ID == id {
			found = &h.data[i]
			break
		}
	}

	if found == nil {
		h.responseError(w, 404, "{{.Name}}不存在")
		return
	}

	{{.Name | ToLower}}JSON, _ := found.ToJSON()
	response := xyJson.NewBuilder().
		SetInt("code", 200).
		SetString("message", "获取成功").
		SetValue("data", {{.Name | ToLower}}JSON).
		MustBuild()

	h.responseJSON(w, 200, response)
}

// update{{.Name}} 更新{{.Name}}
func (h *{{.Name}}Handler) update{{.Name}}(w http.ResponseWriter, r *http.Request, id int64) {
	// 查找{{.Name}}
	var found *models.{{.Name}}
	for i := range h.data {
		if h.data[i].ID == id {
			found = &h.data[i]
			break
		}
	}

	if found == nil {
		h.responseError(w, 404, "{{.Name}}不存在")
		return
	}

	// 读取请求体
	body := make([]byte, r.ContentLength)
	r.Body.Read(body)
	defer r.Body.Close()

	// 解析JSON
	value, err := xyJson.Parse(body)
	if err != nil {
		h.responseError(w, 400, fmt.Sprintf("JSON解析错误: %v", err))
		return
	}

	// 更新{{.Name}}
	if err := found.Update(value.(xyJson.IObject)); err != nil {
		h.responseError(w, 400, fmt.Sprintf("更新失败: %v", err))
		return
	}

	{{.Name | ToLower}}JSON, _ := found.ToJSON()
	response := xyJson.NewBuilder().
		SetInt("code", 200).
		SetString("message", "更新成功").
		SetValue("data", {{.Name | ToLower}}JSON).
		MustBuild()

	h.responseJSON(w, 200, response)
}

// delete{{.Name}} 删除{{.Name}}
func (h *{{.Name}}Handler) delete{{.Name}}(w http.ResponseWriter, r *http.Request, id int64) {
	// 查找并删除{{.Name}}
	for i, item := range h.data {
		if item.ID == id {
			h.data = append(h.data[:i], h.data[i+1:]...)
			response := xyJson.NewBuilder().
				SetInt("code", 200).
				SetString("message", "删除成功").
				MustBuild()
			h.responseJSON(w, 200, response)
			return
		}
	}

	h.responseError(w, 404, "{{.Name}}不存在")
}

// buildListResponse 构建列表响应
func (h *{{.Name}}Handler) buildListResponse(data []models.{{.Name}}, page, pageSize, total int) xyJson.IValue {
	// 构建数据数组
	dataArray := xyJson.CreateArray()
	for _, item := range data {
		itemJSON, _ := item.ToJSON()
		dataArray.Append(itemJSON)
	}

	// 构建响应
	return xyJson.NewBuilder().
		SetInt("code", 200).
		SetString("message", "获取成功").
		SetValue("data", dataArray).
		BeginObject("meta").
			SetInt("total", total).
			SetInt("page", page).
			SetInt("page_size", pageSize).
			SetTime("timestamp", time.Now()).
		End().
		MustBuild()
}

// responseJSON 发送JSON响应
func (h *{{.Name}}Handler) responseJSON(w http.ResponseWriter, statusCode int, data xyJson.IValue) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	jsonStr, err := xyJson.Pretty(data)
	if err != nil {
		http.Error(w, "序列化错误", 500)
		return
	}

	w.Write([]byte(jsonStr))
}

// responseError 发送错误响应
func (h *{{.Name}}Handler) responseError(w http.ResponseWriter, statusCode int, message string) {
	errorResponse := xyJson.NewBuilder().
		SetInt("code", statusCode).
		SetString("message", message).
		SetTime("timestamp", time.Now()).
		MustBuild()

	h.responseJSON(w, statusCode, errorResponse)
}
`

const testTemplate = `package tests

import (
	"testing"
	"time"

	"../models"
	"github.com/your-org/xyJson"
)

func Test{{.Name}}_Creation(t *testing.T) {
	{{.Name | ToLower}} := models.New{{.Name}}()
	if {{.Name | ToLower}} == nil {
		t.Fatal("创建{{.Name}}失败")
	}

	// 验证默认值
	if {{.Name | ToLower}}.CreatedAt.IsZero() {
		t.Error("CreatedAt应该被设置")
	}
	if {{.Name | ToLower}}.UpdatedAt.IsZero() {
		t.Error("UpdatedAt应该被设置")
	}
}

func Test{{.Name}}_JSONConversion(t *testing.T) {
	{{.Name | ToLower}} := models.New{{.Name}}()
{{range .Fields}}{{if eq .Type "string"}}	{{$.Name | ToLower}}.{{.Name}} = "测试{{.Name}}"
{{else if eq .Type "int"}}	{{$.Name | ToLower}}.{{.Name}} = 123
{{else if eq .Type "int64"}}	{{$.Name | ToLower}}.{{.Name}} = 123
{{else if eq .Type "bool"}}	{{$.Name | ToLower}}.{{.Name}} = true
{{end}}{{end}}
	// 转换为JSON
	jsonValue, err := {{.Name | ToLower}}.ToJSON()
	if err != nil {
		t.Fatalf("转换为JSON失败: %v", err)
	}

	// 从JSON创建新实例
	new{{.Name}} := models.New{{.Name}}()
	if err := new{{.Name}}.FromJSON(jsonValue); err != nil {
		t.Fatalf("从JSON创建失败: %v", err)
	}

	// 验证数据
{{range .Fields}}{{if eq .Type "string"}}	if new{{$.Name}}.{{.Name}} != {{$.Name | ToLower}}.{{.Name}} {
		t.Errorf("{{.Name}}不匹配: 期望 %s, 得到 %s", {{$.Name | ToLower}}.{{.Name}}, new{{$.Name}}.{{.Name}})
	}
{{else if eq .Type "int" "int64"}}	if new{{$.Name}}.{{.Name}} != {{$.Name | ToLower}}.{{.Name}} {
		t.Errorf("{{.Name}}不匹配: 期望 %d, 得到 %d", {{$.Name | ToLower}}.{{.Name}}, new{{$.Name}}.{{.Name}})
	}
{{else if eq .Type "bool"}}	if new{{$.Name}}.{{.Name}} != {{$.Name | ToLower}}.{{.Name}} {
		t.Errorf("{{.Name}}不匹配: 期望 %t, 得到 %t", {{$.Name | ToLower}}.{{.Name}}, new{{$.Name}}.{{.Name}})
	}
{{end}}{{end}}
}

func Test{{.Name}}_Validation(t *testing.T) {
	tests := []struct {
		name    string
		{{.Name | ToLower}}    *models.{{.Name}}
		wantErr bool
	}{
		{
			name: "有效数据",
			{{.Name | ToLower}}: &models.{{.Name}}{
{{range .Fields}}{{if .Required}}{{if eq .Type "string"}}				{{.Name}}: "测试{{.Name}}",
{{else if eq .Type "int" "int64"}}				{{.Name}}: 1,
{{end}}{{end}}{{end}}			},
			wantErr: false,
		},
{{range .Fields}}{{if .Required}}		{
			name: "缺少{{.Name}}",
			{{$.Name | ToLower}}: &models.{{$.Name}}{
{{range $.Fields}}{{if and .Required (ne .Name $.Name)}}{{if eq .Type "string"}}				{{.Name}}: "测试{{.Name}}",
{{else if eq .Type "int" "int64"}}				{{.Name}}: 1,
{{end}}{{end}}{{end}}			},
			wantErr: true,
		},
{{end}}{{end}}	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.{{.Name | ToLower}}.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test{{.Name}}_Update(t *testing.T) {
	{{.Name | ToLower}} := models.New{{.Name}}()
{{range .Fields}}{{if eq .Type "string"}}	{{$.Name | ToLower}}.{{.Name}} = "原始{{.Name}}"
{{else if eq .Type "int"}}	{{$.Name | ToLower}}.{{.Name}} = 100
{{else if eq .Type "bool"}}	{{$.Name | ToLower}}.{{.Name}} = false
{{end}}{{end}}
	originalUpdatedAt := {{.Name | ToLower}}.UpdatedAt

	// 创建更新数据
	updates := xyJson.NewBuilder().
{{range .Fields}}{{if ne .Name "ID" "CreatedAt"}}{{if eq .Type "string"}}		SetString("{{.JSONTag}}", "更新的{{.Name}}").
{{else if eq .Type "int"}}		SetInt("{{.JSONTag}}", 200).
{{else if eq .Type "bool"}}		SetBool("{{.JSONTag}}", true).
{{end}}{{end}}{{end}}		MustBuild().(xyJson.IObject)

	// 执行更新
	time.Sleep(time.Millisecond) // 确保时间差异
	if err := {{.Name | ToLower}}.Update(updates); err != nil {
		t.Fatalf("更新失败: %v", err)
	}

	// 验证更新
{{range .Fields}}{{if ne .Name "ID" "CreatedAt"}}{{if eq .Type "string"}}	if {{$.Name | ToLower}}.{{.Name}} != "更新的{{.Name}}" {
		t.Errorf("{{.Name}}未更新")
	}
{{else if eq .Type "int"}}	if {{$.Name | ToLower}}.{{.Name}} != 200 {
		t.Errorf("{{.Name}}未更新")
	}
{{else if eq .Type "bool"}}	if {{$.Name | ToLower}}.{{.Name}} != true {
		t.Errorf("{{.Name}}未更新")
	}
{{end}}{{end}}{{end}}
	// 验证UpdatedAt被更新
	if !{{.Name | ToLower}}.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt应该被更新")
	}
}

func Benchmark{{.Name}}_ToJSON(b *testing.B) {
	{{.Name | ToLower}} := models.New{{.Name}}()
{{range .Fields}}{{if eq .Type "string"}}	{{$.Name | ToLower}}.{{.Name}} = "基准测试{{.Name}}"
{{else if eq .Type "int"}}	{{$.Name | ToLower}}.{{.Name}} = 999
{{else if eq .Type "bool"}}	{{$.Name | ToLower}}.{{.Name}} = true
{{end}}{{end}}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := {{.Name | ToLower}}.ToJSON()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func Benchmark{{.Name}}_FromJSON(b *testing.B) {
	{{.Name | ToLower}} := models.New{{.Name}}()
{{range .Fields}}{{if eq .Type "string"}}	{{$.Name | ToLower}}.{{.Name}} = "基准测试{{.Name}}"
{{else if eq .Type "int"}}	{{$.Name | ToLower}}.{{.Name}} = 999
{{else if eq .Type "bool"}}	{{$.Name | ToLower}}.{{.Name}} = true
{{end}}{{end}}
	jsonValue, _ := {{.Name | ToLower}}.ToJSON()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		new{{.Name}} := models.New{{.Name}}()
		err := new{{.Name}}.FromJSON(jsonValue)
		if err != nil {
			b.Fatal(err)
		}
	}
}
`

const mainTemplate = `package main

import (
	"log"
	"net/http"

	"./handlers"
	"./config"
	"github.com/your-org/xyJson"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 配置xyJson
	setupXyJson(cfg)

	// 创建处理器
	userHandler := handlers.NewUserHandler()

	// 注册路由
	userHandler.RegisterRoutes()
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/metrics", handleMetrics)

	log.Printf("服务器启动在端口 %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

// setupXyJson 配置xyJson
func setupXyJson(cfg *config.Config) {
	// 使用生产环境配置
	xyJsonConfig := xyJson.ProductionConfig()
	xyJson.SetGlobalConfig(xyJsonConfig)

	// 启用性能监控
	if cfg.EnableMonitoring {
		xyJson.EnablePerformanceMonitoring()
	}

	log.Println("xyJson配置完成")
}

// handleHealth 健康检查
func handleHealth(w http.ResponseWriter, r *http.Request) {
	health := xyJson.NewBuilder().
		SetString("status", "healthy").
		SetString("service", "{{.Name}}").
		SetString("version", "{{.Version}}").
		SetTime("timestamp", time.Now()).
		MustBuild()

	w.Header().Set("Content-Type", "application/json")
	jsonStr, _ := xyJson.SerializeToString(health)
	w.Write([]byte(jsonStr))
}

// handleMetrics 性能指标
func handleMetrics(w http.ResponseWriter, r *http.Request) {
	stats := xyJson.GetPerformanceStats()

	metrics := xyJson.NewBuilder().
		BeginObject("performance").
			SetInt("parse_count", int(stats.ParseCount)).
			SetInt("serialize_count", int(stats.SerializeCount)).
			SetString("avg_parse_time", stats.AvgParseTime.String()).
			SetString("avg_serialize_time", stats.AvgSerializeTime.String()).
		End().
		SetTime("timestamp", time.Now()).
		MustBuild()

	w.Header().Set("Content-Type", "application/json")
	jsonStr, _ := xyJson.Pretty(metrics)
	w.Write([]byte(jsonStr))
}
`

const goModTemplate = `module {{.Module}}

go 1.21

require (
	github.com/your-org/xyJson v1.0.0
)
`

const dockerfileTemplate = `FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o {{.Name}} .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/{{.Name}} .

EXPOSE 8080
CMD ["./{{.Name}}"]
`

const readmeTemplate = `# {{.Name}}

{{.Description}}

## 特性

- 基于 xyJson 的高性能 JSON 处理
- RESTful API 设计
- 完整的 CRUD 操作
- 性能监控和指标
- Docker 支持
- 完整的测试覆盖

## 快速开始

### 安装依赖

` + "```bash" + `
go mod download
` + "```" + `

### 运行服务

` + "```bash" + `
go run main.go
` + "```" + `

### 使用 Docker

` + "```bash" + `
docker build -t {{.Name}} .
docker run -p 8080:8080 {{.Name}}
` + "```" + `

## API 文档

### 用户管理

- ` + "`GET /users`" + ` - 获取用户列表
- ` + "`POST /users`" + ` - 创建用户
- ` + "`GET /users/{id}`" + ` - 获取用户详情
- ` + "`PUT /users/{id}`" + ` - 更新用户
- ` + "`DELETE /users/{id}`" + ` - 删除用户

### 系统接口

- ` + "`GET /health`" + ` - 健康检查
- ` + "`GET /metrics`" + ` - 性能指标

## 测试

` + "```bash" + `
go test ./...
` + "```" + `

## 许可证

{{.License}}

## 作者

{{.Author}}
`

const configTemplate = `package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Port             string
	DatabaseURL      string
	EnableMonitoring bool
	LogLevel         string
	JWTSecret        string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		DatabaseURL:      getEnv("DATABASE_URL", "sqlite://./app.db"),
		EnableMonitoring: getEnvBool("ENABLE_MONITORING", true),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
	}
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔型环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
`

// 模板函数
func init() {
	template.Must(template.New("").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(""))
}
