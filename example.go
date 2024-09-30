package example

import (
	"context"
	"encoding/json"
	"fmt"
	"go-hutool/async"
	"go-hutool/list"
	_map "go-hutool/map"
	"log"
	"os"
	"time"
)

func testAsync() {
	logger := log.New(os.Stdout, "自定义日志：", log.LstdFlags)
	ctx := context.Background()

	executor := async.NewAsyncExecutor(5,
		async.WithLogger(logger),
		async.WithContext(ctx),
	)

	// 执行单个任务
	executor.Execute(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second)
		return "任务1完成", nil
	})

	// 执行带超时的任务
	result := executor.ExecuteWithTimeout(func(ctx context.Context) (interface{}, error) {
		select {
		case <-time.After(2 * time.Second):
			return "任务2完成", nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}, 3*time.Second)
	fmt.Printf("任务2结果：%v，错误：%v\n", result.Value, result.Err)

	// 执行多个任务
	tasks := []async.Task{
		func(ctx context.Context) (interface{}, error) { return "任务3", nil },
		func(ctx context.Context) (interface{}, error) { return "任务4", nil },
		func(ctx context.Context) (interface{}, error) { panic("任务5发生panic") },
	}

	results := executor.ExecuteAll(tasks...)
	for i, result := range results {
		fmt.Printf("任务%d结果：%v，错误：%v\n", i+3, result.Value, result.Err)
	}

	// 执行带超时的多个任务
	timeoutResults := executor.ExecuteAllWithTimeout(2*time.Second, tasks...)
	fmt.Printf("超时结果：%v\n", timeoutResults)

	// 优雅关闭
	if err := executor.Shutdown(5 * time.Second); err != nil {
		logger.Printf("关闭过程中出错：%v", err)
	}
}

func testList() {
	// 创建整数列表
	intList := list.New(1, 2, 3, 4, 5)
	fmt.Println("原始列表:", intList)

	// 添加元素
	intList.Add(6)
	fmt.Println("添加元素后:", intList)

	// 获取元素
	if val, err := intList.Get(2); err == nil {
		fmt.Println("索引2的元素:", val)
	}

	// 移除元素
	intList.Remove(3)
	fmt.Println("移除索引3的元素后:", intList)

	// 过滤偶数
	evenList := intList.Filter(func(i int) bool { return i%2 == 0 })
	fmt.Println("偶数列表:", evenList)

	// 将每个元素乘2
	doubledList := intList.Map(func(i int) int { return i * 2 })
	fmt.Println("每个元素乘2后:", doubledList)

	// 求和
	sum := intList.Reduce(func(a, b int) int { return a + b }, 0)
	fmt.Println("列表元素之和:", sum)

	// 排序
	intList.Sort(func(a, b int) bool { return a < b })
	fmt.Println("排序后:", intList)

	// 反转
	intList.Reverse()
	fmt.Println("反转后:", intList)

	//遍历
	intList.ForEach(func(i int) {
		println(i)
	})
}

func testMap() {
	// 创建字符串到整数的映射
	m := _map.New[string, int]()

	// 添加键值对
	m.Put("one", 1)
	m.Put("two", 2)
	m.Put("three", 3)

	// 获取值
	if val, exists := m.Get("two"); exists {
		fmt.Println("Value of 'two':", val)
	}

	// 移除键值对
	m.Remove("three")

	// 遍历映射
	m.ForEach(func(k string, v int) {
		fmt.Printf("%s: %d\n", k, v)
	})

	// 过滤
	evenMap := m.Filter(func(k string, v int) bool {
		return v%2 == 0
	})
	fmt.Println("Even numbers:", evenMap)

	// 映射值
	doubledMap := m.Map(func(k string, v int) int {
		return v * 2
	})
	fmt.Println("Doubled values:", doubledMap)

	// 归约
	sum := m.Reduce(func(acc int, k string, v int) int {
		return acc + v
	}, 0)
	fmt.Println("Sum of values:", sum)

	// JSON 序列化
	jsonStr, _ := m.ToJSON()
	fmt.Println("JSON representation:", jsonStr)

	// 从 JSON 创建新映射
	newMap, _ := _map.FromJSON[string, int](jsonStr)
	fmt.Println("Map from JSON:", newMap)

	// 合并映射
	other := _map.New[string, int]()
	other.Put("four", 4)
	m.Merge(other, func(v1, v2 int) int {
		return v1 // 在冲突时保留原值
	})
	fmt.Println("Merged map:", m)

	// 获取或计算
	value := m.GetOrCompute("five", func() int {
		return 5
	})
	fmt.Println("Computed value for 'five':", value)

	// 更新值
	m.Update("one", func(v int) int {
		return v + 10
	})
	fmt.Println("Updated map:", m)
}

func testJson() {
	ju := json.New()

	// 基本操作
	data := map[string]interface{}{
		"name": "John Doe",
		"age":  30,
		"city": "New York",
	}

	jsonStr, err := ju.ToJSON(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("JSON string:", jsonStr)

	// 格式化
	prettyJSON, _ := ju.PrettyPrint(data)
	fmt.Println("Pretty JSON:")
	fmt.Println(prettyJSON)

	// 验证
	isValid := ju.IsValidJSON(jsonStr)
	fmt.Println("Is valid JSON:", isValid)

	// 路径查询
	value, _ := ju.GetValueByPath(jsonStr, "name")
	fmt.Println("Name:", value)

	// 合并
	json1 := `{"a": 1, "b": {"c": 2}}`
	json2 := `{"b": {"d": 3}, "e": 4}`
	mergedJSON, _ := ju.MergeJSON(json1, json2)
	fmt.Println("Merged JSON:", mergedJSON)

	// 比较
	diff, _ := ju.Diff(json1, json2)
	fmt.Println("Diff:", diff)

	// 文件操作
	err = ju.ToJSONFile(data, "output.json")
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	var readData map[string]interface{}
	err = ju.FromJSONFile("output.json", &readData)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println("Data read from file:", readData)

	// 流式处理
	file, _ := os.Open("large.json")
	defer file.Close()

	err = ju.StreamingDecode(file, func(t json.Token) error {
		fmt.Printf("Token: %v\n", t)
		return nil
	})
	if err != nil {
		fmt.Println("Error in streaming decode:", err)
	}

}
