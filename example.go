package example

import (
	"context"
	"encoding/json"
	"fmt"
	asyncutil "go-hutool/async"
	"go-hutool/codec"
	"go-hutool/datetime"
	"go-hutool/desensitized"
	jsonutil "go-hutool/json"
	listutil "go-hutool/list"
	maputil "go-hutool/map"
	randutil "go-hutool/random"
	"log"
	"os"
	"time"
)

func testAsync() {
	logger := log.New(os.Stdout, "自定义日志：", log.LstdFlags)
	ctx := context.Background()

	executor := asyncutil.NewAsyncExecutor(5,
		asyncutil.WithLogger(logger),
		asyncutil.WithContext(ctx),
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
	tasks := []asyncutil.Task{
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
	intList := listutil.New(1, 2, 3, 4, 5)
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
	m := maputil.New[string, int]()

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
	newMap, _ := maputil.FromJSON[string, int](jsonStr)
	fmt.Println("Map from JSON:", newMap)

	// 合并映射
	other := maputil.New[string, int]()
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
	ju := jsonutil.New()

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

func testDateTime() {
	// 创建 DateTimeUtil 实例，使用本地时区
	dtu := datetime.New(time.Local)

	// 获取当前时间
	now := dtu.Now()
	fmt.Println("Current time:", now)

	// 解析时间字符串
	t, err := dtu.Parse("2006-01-02 15:04:05", "2023-05-01 10:30:00")
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println("Parsed time:", t)

	// 格式化时间
	formatted := dtu.Format(t, "2006-01-02 15:04:05")
	fmt.Println("Formatted time:", formatted)

	// 增加时间
	newTime := dtu.AddDuration(t, 24*time.Hour)
	fmt.Println("Time after adding 24 hours:", newTime)

	// 计算天数差
	days := dtu.DiffDays(t, newTime)
	fmt.Println("Days between:", days)

	// 判断闰年
	isLeap := dtu.IsLeapYear(2024)
	fmt.Println("Is 2024 a leap year?", isLeap)

	// 获取周几
	weekday := dtu.GetWeekday(t)
	fmt.Println("Weekday:", weekday)

	// 获取一年中的第几周
	week := dtu.GetWeekOfYear(t)
	fmt.Println("Week of year:", week)

	// 获取月份天数
	daysInMonth := dtu.GetDaysInMonth(2023, 5)
	fmt.Println("Days in May 2023:", daysInMonth)

	// 获取日期范围
	startOfDay := dtu.StartOfDay(t)
	endOfDay := dtu.EndOfDay(t)
	fmt.Println("Start of day:", startOfDay)
	fmt.Println("End of day:", endOfDay)

	// 添加工作日
	workday := dtu.AddWorkdays(t, 5)
	fmt.Println("After adding 5 workdays:", workday)

	// 计算年龄
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	age := dtu.Age(birthDate)
	fmt.Println("Age:", age)

	// 获取下一个周一
	nextMonday := dtu.NextOccurrence(t, time.Monday)
	fmt.Println("Next Monday:", nextMonday)

	// 格式化持续时间
	duration := 36*time.Hour + 15*time.Minute + 45*time.Second
	durationStr := dtu.DurationString(duration)
	fmt.Println("Formatted duration:", durationStr)

	// 解析持续时间
	parsedDuration, err := dtu.ParseDuration("2 weeks")
	if err != nil {
		fmt.Println("Error parsing duration:", err)
	} else {
		fmt.Println("Parsed duration:", parsedDuration)
	}
}

func testRand() {
	r := randutil.New()

	// 生成随机整数
	randInt, _ := r.Int(1, 100)
	fmt.Println("Random int:", randInt)

	// 生成随机浮点数
	randFloat, _ := r.Float64(0, 1)
	fmt.Println("Random float:", randFloat)

	// 生成随机字符串
	randStr, _ := r.String(10)
	fmt.Println("Random string:", randStr)

	// 生成随机数字字符串
	randDigits, _ := r.Digits(6)
	fmt.Println("Random digits:", randDigits)

	// 生成 UUID
	uuid, _ := r.UUID()
	fmt.Println("UUID:", uuid)

	// 随机打乱字符串
	shuffled, _ := r.ShuffleString("Hello, World!")
	fmt.Println("Shuffled string:", shuffled)

	// 随机选择元素
	fruits := []string{"apple", "banana", "orange", "grape"}
	choice, _ := r.Choice(fruits)
	fmt.Println("Random fruit:", choice)

	// 生成随机日期
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
	randomDate, _ := r.RandomDate(start, end)
	fmt.Println("Random date:", randomDate.Format("2006-01-02"))

	// 生成随机密码
	password, _ := r.Password(12)
	fmt.Println("Random password:", password)
}

func testCodec() {
	// Base64 编码和解码
	originalData := []byte("Hello, World!")
	encoded := codec.Base64Encode(originalData)
	fmt.Println("Base64 encoded:", encoded)
	decoded, _ := codec.Base64Decode(encoded)
	fmt.Println("Base64 decoded:", string(decoded))

	// Base62 编码和解码
	base62Encoded := codec.Base62Encode(originalData)
	fmt.Println("Base62 encoded:", base62Encoded)
	base62Decoded, _ := codec.Base62Decode(base62Encoded)
	fmt.Println("Base62 decoded:", string(base62Decoded))

	// MD5 哈希
	md5Hash := codec.MD5("Hello, World!")
	fmt.Println("MD5 hash:", md5Hash)

	// SHA256 哈希
	sha256Hash := codec.SHA256("Hello, World!")
	fmt.Println("SHA256 hash:", sha256Hash)

	// ROT13 编码/解码
	rot13Encoded := codec.ROT13("Hello, World!")
	fmt.Println("ROT13 encoded:", rot13Encoded)
	rot13Decoded := codec.ROT13(rot13Encoded)
	fmt.Println("ROT13 decoded:", rot13Decoded)

	// URL 编码和解码
	urlEncoded := codec.URLEncode("Hello, World! 你好，世界！")
	fmt.Println("URL encoded:", urlEncoded)
	urlDecoded, _ := codec.URLDecode(urlEncoded)
	fmt.Println("URL decoded:", urlDecoded)

	// XOR 加密和解密
	key := []byte("secret")
	xorEncrypted := codec.XOREncrypt(originalData, key)
	fmt.Println("XOR encrypted:", xorEncrypted)
	xorDecrypted := codec.XORDecrypt(xorEncrypted, key)
	fmt.Println("XOR decrypted:", string(xorDecrypted))

	// 凯撒密码加密和解密
	caesarEncrypted := codec.CaesarEncrypt("Hello, World!", 3)
	fmt.Println("Caesar encrypted:", caesarEncrypted)
	caesarDecrypted := codec.CaesarDecrypt(caesarEncrypted, 3)
	fmt.Println("Caesar decrypted:", caesarDecrypted)
}

func testDesensitizedUtil() {
	fmt.Println(desensitized.IDCardNum("51343620000320711X", 1, 2)) // 5***************1X
	fmt.Println(desensitized.MobilePhone("18049531999"))            // 180****1999
	fmt.Println(desensitized.Password("1234567890"))                // **********
	fmt.Println(desensitized.ChineseName("张三丰"))                    // 张**
	fmt.Println(desensitized.Email("zhangsan@example.com"))         // z*******n@example.com
	fmt.Println(desensitized.BankCard("6222021234567890123"))       // 622202******0123
	fmt.Println(desensitized.Address("北京市朝阳区长安街1号"))                // 北京市朝阳**********号
	fmt.Println(desensitized.LicensePlate("京A12345"))               // 京****5
	fmt.Println(desensitized.Landline("010-12345678"))              // 010-****5678
}
