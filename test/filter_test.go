package test

import (
	"testing"

	xyJson "github.com/ihuem/xyJson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFilter 测试Filter函数的基本功能
// TestFilter tests the basic functionality of the Filter function
func TestFilter(t *testing.T) {
	// 准备测试数据
	data := `{
		"employees": [
			{"name": "Alice", "salary": 120000, "department": "Engineering"},
			{"name": "Bob", "salary": 80000, "department": "Marketing"},
			{"name": "Charlie", "salary": 150000, "department": "Engineering"},
			{"name": "David", "salary": 90000, "department": "Sales"},
			{"name": "Eve", "salary": 110000, "department": "Engineering"}
		]
	}`

	root, err := xyJson.ParseString(data)
	require.NoError(t, err)

	t.Run("filter_high_earners", func(t *testing.T) {
		// 过滤高薪员工（薪资 > 100000）
		highEarners, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
			salary, err := xyJson.Get(emp, "$.salary")
			if err != nil {
				return false
			}
			if scalarValue, ok := salary.(xyJson.IScalarValue); ok {
				num, _ := scalarValue.Float64()
				return num > 100000
			}
			return false
		})

		assert.NoError(t, err)
		assert.Len(t, highEarners, 3) // Alice, Charlie, Eve

		// 验证结果
		expectedNames := []string{"Alice", "Charlie", "Eve"}
		for i, emp := range highEarners {
			name, err := xyJson.GetString(emp, "$.name")
			assert.NoError(t, err)
			assert.Contains(t, expectedNames, name)
			_ = i // 避免未使用变量警告
		}
	})

	t.Run("filter_by_department", func(t *testing.T) {
		// 过滤工程部门员工
		engineers, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
			dept, err := xyJson.GetString(emp, "$.department")
			if err != nil {
				return false
			}
			return dept == "Engineering"
		})

		assert.NoError(t, err)
		assert.Len(t, engineers, 3) // Alice, Charlie, Eve

		// 验证所有结果都是工程部门
		for _, emp := range engineers {
			dept, err := xyJson.GetString(emp, "$.department")
			assert.NoError(t, err)
			assert.Equal(t, "Engineering", dept)
		}
	})

	t.Run("filter_empty_result", func(t *testing.T) {
		// 过滤不存在的条件
		result, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
			salary, err := xyJson.Get(emp, "$.salary")
			if err != nil {
				return false
			}
			if scalarValue, ok := salary.(xyJson.IScalarValue); ok {
				num, _ := scalarValue.Float64()
				return num > 200000 // 没有员工薪资超过200000
			}
			return false
		})

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("filter_all_match", func(t *testing.T) {
		// 所有员工都匹配的条件
		result, err := xyJson.Filter(root, "$.employees[*]", func(emp xyJson.IValue) bool {
			salary, err := xyJson.Get(emp, "$.salary")
			if err != nil {
				return false
			}
			if scalarValue, ok := salary.(xyJson.IScalarValue); ok {
				num, _ := scalarValue.Float64()
				return num > 0 // 所有员工薪资都大于0
			}
			return false
		})

		assert.NoError(t, err)
		assert.Len(t, result, 5) // 所有5个员工
	})
}

// TestFilterErrorCases 测试Filter函数的错误情况
// TestFilterErrorCases tests error cases of the Filter function
func TestFilterErrorCases(t *testing.T) {
	data := `{"items": [1, 2, 3]}`
	root, err := xyJson.ParseString(data)
	require.NoError(t, err)

	t.Run("nil_predicate", func(t *testing.T) {
		// 测试nil谓词函数
		result, err := xyJson.Filter(root, "$.items[*]", nil)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "predicate cannot be nil")
	})

	t.Run("invalid_path", func(t *testing.T) {
		// 测试无效路径
		result, err := xyJson.Filter(root, "$.nonexistent[*]", func(item xyJson.IValue) bool {
			return true
		})
		assert.NoError(t, err) // GetAll对于不存在的路径返回空数组，不是错误
		assert.Len(t, result, 0)
	})

	t.Run("malformed_path", func(t *testing.T) {
		// 测试格式错误的路径
		result, err := xyJson.Filter(root, "$.[invalid", func(item xyJson.IValue) bool {
			return true
		})
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestMustFilter 测试MustFilter函数
// TestMustFilter tests the MustFilter function
func TestMustFilter(t *testing.T) {
	data := `{"numbers": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}`
	root, err := xyJson.ParseString(data)
	require.NoError(t, err)

	t.Run("must_filter_success", func(t *testing.T) {
		// 过滤偶数
		evenNumbers := xyJson.MustFilter(root, "$.numbers[*]", func(num xyJson.IValue) bool {
			if scalarValue, ok := num.(xyJson.IScalarValue); ok {
				value, _ := scalarValue.Float64()
				return int(value)%2 == 0
			}
			return false
		})

		assert.Len(t, evenNumbers, 5) // 2, 4, 6, 8, 10

		// 验证所有结果都是偶数
		for _, num := range evenNumbers {
			if scalarValue, ok := num.(xyJson.IScalarValue); ok {
				value, _ := scalarValue.Float64()
				assert.Equal(t, 0, int(value)%2)
			}
		}
	})

	t.Run("must_filter_with_nil_predicate", func(t *testing.T) {
		// 测试nil谓词函数，应该返回空数组而不是panic
		result := xyJson.MustFilter(root, "$.numbers[*]", nil)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})

	t.Run("must_filter_with_invalid_path", func(t *testing.T) {
		// 测试无效路径，应该返回空数组而不是panic
		result := xyJson.MustFilter(root, "$.[invalid", func(item xyJson.IValue) bool {
			return true
		})
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
	})
}

// TestFilterComplexData 测试复杂数据结构的过滤
// TestFilterComplexData tests filtering on complex data structures
func TestFilterComplexData(t *testing.T) {
	data := `{
		"products": [
			{
				"id": 1,
				"name": "Laptop",
				"price": 999.99,
				"category": "Electronics",
				"inStock": true,
				"tags": ["computer", "portable"]
			},
			{
				"id": 2,
				"name": "Book",
				"price": 29.99,
				"category": "Education",
				"inStock": false,
				"tags": ["learning", "paper"]
			},
			{
				"id": 3,
				"name": "Phone",
				"price": 699.99,
				"category": "Electronics",
				"inStock": true,
				"tags": ["mobile", "communication"]
			}
		]
	}`

	root, err := xyJson.ParseString(data)
	require.NoError(t, err)

	t.Run("filter_expensive_electronics", func(t *testing.T) {
		// 过滤价格超过500的电子产品
		expensiveElectronics, err := xyJson.Filter(root, "$.products[*]", func(product xyJson.IValue) bool {
			category, err := xyJson.GetString(product, "$.category")
			if err != nil || category != "Electronics" {
				return false
			}

			price, err := xyJson.GetFloat64(product, "$.price")
			if err != nil {
				return false
			}

			return price > 500
		})

		assert.NoError(t, err)
		assert.Len(t, expensiveElectronics, 2) // Laptop and Phone

		// 验证结果
		for _, product := range expensiveElectronics {
			category, err := xyJson.GetString(product, "$.category")
			assert.NoError(t, err)
			assert.Equal(t, "Electronics", category)

			price, err := xyJson.GetFloat64(product, "$.price")
			assert.NoError(t, err)
			assert.Greater(t, price, 500.0)
		}
	})

	t.Run("filter_in_stock_products", func(t *testing.T) {
		// 过滤有库存的产品
		inStockProducts, err := xyJson.Filter(root, "$.products[*]", func(product xyJson.IValue) bool {
			inStock, err := xyJson.GetBool(product, "$.inStock")
			if err != nil {
				return false
			}
			return inStock
		})

		assert.NoError(t, err)
		assert.Len(t, inStockProducts, 2) // Laptop and Phone

		// 验证所有结果都有库存
		for _, product := range inStockProducts {
			inStock, err := xyJson.GetBool(product, "$.inStock")
			assert.NoError(t, err)
			assert.True(t, inStock)
		}
	})
}
