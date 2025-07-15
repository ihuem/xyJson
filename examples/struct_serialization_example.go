package main

import (
	"fmt"
	"log"
	"time"

	xyJson "github.com/ihuem/xyJson"
)

// 基础示例结构体
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// 复杂结构体示例

type Users struct {
	ID       int               `json:"id"`
	Name     string            `json:"name"`
	Email    string            `json:"email"`
	Active   bool              `json:"active"`
	Balance  float64           `json:"balance"`
	Tags     []string          `json:"tags"`
	Scores   []int             `json:"scores"`
	Metadata map[string]string `json:"metadata"`
}

// 嵌套结构体示例

/* type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip"`
} */

type UserWithAddress struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Address Address `json:"address"`
}

// 指针字段示例
type UserWithPointers struct {
	ID      *int     `json:"id"`
	Name    *string  `json:"name"`
	Active  *bool    `json:"active"`
	Balance *float64 `json:"balance"`
}

// JSON标签示例
type UserWithCustomTags struct {
	ID       int    `json:"user_id"`    // 自定义字段名
	Name     string `json:"full_name"`  // 自定义字段名
	Email    string `json:"-"`          // 忽略字段
	Password string `json:",omitempty"` // 空值忽略
	Active   bool   `json:"is_active"`  // 自定义字段名
}

// 时间字段示例
type Event struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main8() {
	fmt.Println("=== xyJson 序列化到结构体示例 ===")
	fmt.Println()

	// 1. 基础用法示例
	basicExample()
	fmt.Println()

	// 2. 复杂结构体示例
	complexExample()
	fmt.Println()

	// 3. 嵌套结构体示例
	nestedExample()
	fmt.Println()

	// 4. 指针字段示例
	pointerExample()
	fmt.Println()

	// 5. JSON标签示例
	customTagsExample()
	fmt.Println()

	// 6. 时间字段示例
	timeExample()
	fmt.Println()

	// 7. 直接解析示例
	directUnmarshalExample()
	fmt.Println()

	// 8. 错误处理示例
	errorHandlingExample()
	fmt.Println()

	// 9. Must版本示例
	mustVersionExample()
}

// 基础用法示例
func basicExample() {
	fmt.Println("1. 基础用法示例:")

	// JSON数据
	jsonData := `{"name":"Alice","age":25}`
	fmt.Printf("JSON数据: %s\n", jsonData)

	// 解析JSON
	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// 序列化到结构体
	var person Person
	err = xyJson.SerializeToStruct(value, &person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("结果: Name=%s, Age=%d\n", person.Name, person.Age)
}

// 复杂结构体示例
func complexExample() {
	fmt.Println("2. 复杂结构体示例:")

	jsonData := `{
		"id": 123,
		"name": "John Doe",
		"email": "john@example.com",
		"active": true,
		"balance": 1234.56,
		"tags": ["admin", "user", "premium"],
		"scores": [85, 92, 78, 96],
		"metadata": {
			"department": "engineering",
			"role": "senior",
			"location": "remote"
		}
	}`

	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var user Users
	err = xyJson.SerializeToStruct(value, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用户ID: %d\n", user.ID)
	fmt.Printf("姓名: %s\n", user.Name)
	fmt.Printf("邮箱: %s\n", user.Email)
	fmt.Printf("激活状态: %t\n", user.Active)
	fmt.Printf("余额: %.2f\n", user.Balance)
	fmt.Printf("标签: %v\n", user.Tags)
	fmt.Printf("分数: %v\n", user.Scores)
	fmt.Printf("元数据: %v\n", user.Metadata)
}

// 嵌套结构体示例
func nestedExample() {
	fmt.Println("3. 嵌套结构体示例:")

	jsonData := `{
		"id": 456,
		"name": "Jane Smith",
		"address": {
			"street": "123 Main St",
			"city": "New York",
			"zip": "10001"
		}
	}`

	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var user UserWithAddress
	err = xyJson.SerializeToStruct(value, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用户: %s (ID: %d)\n", user.Name, user.ID)
	fmt.Printf("地址: %s, %s %s\n", user.Address.Street, user.Address.City, user.Address.Zip)
}

// 指针字段示例
func pointerExample() {
	fmt.Println("4. 指针字段示例:")

	jsonData := `{
		"id": 789,
		"name": "Bob Wilson",
		"active": false,
		"balance": 567.89
	}`

	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var user UserWithPointers
	err = xyJson.SerializeToStruct(value, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用户ID: %v\n", *user.ID)
	fmt.Printf("姓名: %v\n", *user.Name)
	fmt.Printf("激活状态: %v\n", *user.Active)
	fmt.Printf("余额: %v\n", *user.Balance)
}

// JSON标签示例
func customTagsExample() {
	fmt.Println("5. JSON标签示例:")

	jsonData := `{
		"user_id": 999,
		"full_name": "Custom User",
		"email": "custom@example.com",
		"password": "secret123",
		"is_active": true
	}`

	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var user UserWithCustomTags
	err = xyJson.SerializeToStruct(value, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("用户ID (user_id -> ID): %d\n", user.ID)
	fmt.Printf("姓名 (full_name -> Name): %s\n", user.Name)
	fmt.Printf("邮箱 (忽略字段): '%s' (应该为空)\n", user.Email)
	fmt.Printf("密码: %s\n", user.Password)
	fmt.Printf("激活状态 (is_active -> Active): %t\n", user.Active)
}

// 时间字段示例
func timeExample() {
	fmt.Println("6. 时间字段示例:")

	jsonData := `{
		"id": 111,
		"name": "Time Event",
		"created_at": "2023-01-15T10:30:00Z",
		"updated_at": "2023-12-25T15:45:30Z"
	}`

	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	var event Event
	err = xyJson.SerializeToStruct(value, &event)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("事件: %s (ID: %d)\n", event.Name, event.ID)
	fmt.Printf("创建时间: %s\n", event.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("更新时间: %s\n", event.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// 直接解析示例
func directUnmarshalExample() {
	fmt.Println("7. 直接解析示例:")

	// 使用UnmarshalToStruct直接从JSON字节数组解析
	jsonData := []byte(`{"name":"Direct User","age":35}`)

	var person Person
	err := xyJson.UnmarshalToStruct(jsonData, &person)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("直接解析结果: Name=%s, Age=%d\n", person.Name, person.Age)

	// 使用UnmarshalStringToStruct直接从JSON字符串解析
	jsonString := `{"name":"String User","age":40}`

	var person2 Person
	err = xyJson.UnmarshalStringToStruct(jsonString, &person2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("字符串解析结果: Name=%s, Age=%d\n", person2.Name, person2.Age)
}

// 错误处理示例
func errorHandlingExample() {
	fmt.Println("8. 错误处理示例:")

	jsonData := `{"name":"Error Test","age":"invalid_age"}`
	value, err := xyJson.ParseString(jsonData)
	if err != nil {
		fmt.Printf("解析错误: %v\n", err)
		return
	}

	var person Person
	err = xyJson.SerializeToStruct(value, &person)
	if err != nil {
		fmt.Printf("序列化错误: %v\n", err)
		fmt.Println("这是预期的错误，因为age字段不是有效的整数")
	} else {
		fmt.Printf("意外成功: %+v\n", person)
	}

	// 测试其他错误情况
	fmt.Println("\n测试其他错误情况:")

	// nil target
	err = xyJson.SerializeToStruct(value, nil)
	if err != nil {
		fmt.Printf("nil target错误: %v\n", err)
	}

	// 非指针target
	var person2 Person
	err = xyJson.SerializeToStruct(value, person2)
	if err != nil {
		fmt.Printf("非指针target错误: %v\n", err)
	}

	// 指向非结构体的指针
	var str string
	err = xyJson.SerializeToStruct(value, &str)
	if err != nil {
		fmt.Printf("非结构体target错误: %v\n", err)
	}
}

// Must版本示例
func mustVersionExample() {
	fmt.Println("9. Must版本示例:")

	// Must版本在成功时正常工作
	jsonData := `{"name":"Must User","age":50}`
	value := xyJson.MustParseString(jsonData)

	var person Person
	xyJson.MustSerializeToStruct(value, &person)

	fmt.Printf("Must版本结果: Name=%s, Age=%d\n", person.Name, person.Age)

	// 直接使用MustUnmarshalToStruct
	jsonBytes := []byte(`{"name":"Must Direct","age":60}`)
	var person2 Person
	xyJson.MustUnmarshalToStruct(jsonBytes, &person2)

	fmt.Printf("Must直接解析结果: Name=%s, Age=%d\n", person2.Name, person2.Age)

	fmt.Println("\n注意: Must版本在遇到错误时会panic，请谨慎使用")
}
