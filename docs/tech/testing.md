# Zenfeed 最新测试策略与风格

## 1. 引言

Zenfeed 的测试策略核心目标是：

*   **清晰性 (Clarity)**：测试本身应如文档般易于理解，清晰地表达被测功能的行为和预期。
*   **可信性 (Reliability)**：测试结果应准确反映代码的健康状况，确保每次提交的信心。
*   **可维护性 (Maintainability)**：测试代码应易于修改和扩展，以适应项目的持续演进。

本指南旨在详细介绍 Zenfeed 项目所遵循的测试理念、风格和具体实践。

## 2. 核心测试理念与风格

Zenfeed 的测试方法论深受行为驱动开发 (BDD) 的影响，并结合了表驱动测试等高效实践。

### 2.1 行为驱动开发

我们选择 BDD 作为核心的测试描述框架，主要基于以下原因（其理念也体现在 `pkg/test/test.go` 的 `Case` 结构设计中）：

*   **提升可读性 (Enhanced Readability)**：BDD 强调使用自然语言描述软件的行为。每个测试用例读起来都像一个用户故事或一个功能说明，这使得测试本身就成为了一种精确的"活文档"。
*   **关注行为 (Focus on Behavior)**：测试不再仅仅是验证代码片段的输入输出，而是从模块、组件或用户交互的层面描述其应有的行为。这有助于确保我们构建的功能符合预期。
*   **需求驱动 (Requirement-Driven)**：测试直接对应需求描述，而非实现细节。这种自顶向下的方法确保了测试的稳定性，即使内部实现重构，只要行为不变，测试依然有效。

BDD 通常使用 `Scenario`, `Given`, `When`, `Then` 的结构来组织测试：

*   **`Scenario` (场景)**：描述测试用例所针对的特性或功能点。
    *   例如：`"Query hot block with label filters"` (查询带标签过滤的热数据块)
*   **`Given` (给定)**：描述场景开始前的初始上下文或状态（**注意：这不是指方法的输入参数**）。
    *   例如：`"a hot block with indexed feeds"` (一个已索引了 Feed 的热数据块)
*   **`When` (当)**：描述触发场景的事件或操作（**这部分通常包含被测方法的输入参数**）。
    *   例如：`"querying with label filters"` (当使用标签过滤器进行查询时)
*   **`Then` (那么)**：描述场景结束后预期的结果或状态变化。
    *   例如：`"should return matching feeds"` (那么应该返回匹配的 Feed)

为了更好地在代码中实践 BDD，我们定义了 `pkg/test/test.go` 中的 `Case[GivenDetail, WhenDetail, ThenExpected]` 泛型结构。其中：

*   `GivenDetail`: 存储 `Given` 子句描述的初始状态的具体数据。
*   `WhenDetail`: 存储 `When` 子句描述的事件或方法调用的具体参数。
*   `ThenExpected`: 存储 `Then` 子句描述的预期结果。

这种结构化不仅增强了测试数据的类型安全，也使得测试用例的意图更加明确。对于需要模拟依赖项的组件，`GivenDetail` 通常会包含用于配置这些模拟行为的 `component.MockOption`，我们将在后续 Mocking 章节详细讨论。

### 2.2 表驱动测试

当一个功能或方法需要针对多种不同的输入组合、边界条件或状态进行测试时，表驱动测试是一种非常高效和整洁的组织方式。

*   **简洁性 (Conciseness)**：将所有测试用例的数据（输入、参数、预期输出）集中定义在一个表格（通常是切片）中，避免了为每个 case编写大量重复的测试逻辑。
*   **易扩展性 (Extensibility)**：添加新的测试场景变得非常简单，只需在表格中增加一条新记录即可。
*   **清晰性 (Clarity)**：所有相关的测试用例一目了然，便于快速理解被测功能的覆盖范围。

**实践约定**：
在 Zenfeed 中，**当存在多个测试用例时，必须使用表驱动测试**。

### 2.3 测试结构约定

为了保持项目范围内测试代码的一致性和可读性，我们约定在测试文件中遵循以下组织结构：

1.  **定义辅助类型 (Define Helper Types)**：在测试函数的开头部分，通常会为 `GivenDetail`, `WhenDetail`, `ThenExpected` 定义具体的结构体类型，以增强类型安全和表达力。
2.  **定义测试用例表 (Define Test Case Table)**：将所有测试用例集中定义在一个 `[]test.Case` 类型的切片中。
3.  **循环执行测试 (Loop Through Test Cases)**：使用 `for` 循环遍历测试用例表，并为每个用例运行 `t.Run(tt.Scenario, func(t *testing.T) { ... })`。
4.  **清晰的 G/W/T 逻辑块 (Clear G/W/T Blocks)**：在每个 `t.Run` 的匿名函数内部，根据需要组织代码块，以对应 `Given`（准备初始状态，通常基于 `tt.GivenDetail`），`When`（执行被测操作，通常使用 `tt.WhenDetail`），和 `Then`（断言结果，通常对比 `tt.ThenExpected`）。
5.  **描述性变量名 (Descriptive Variable Names)**：使用与 BDD 术语（如 `given`, `when`, `then`, `expected`, `actual`）相匹配或能清晰表达意图的变量名。

## 3. 依赖隔离：Mocking (Dependency Isolation: Mocking)

单元测试的核心原则之一是**隔离性 (Isolation)**，即被测试的代码单元（如一个函数或一个方法）应该与其依赖项隔离开来。Mocking (模拟) 是实现这种隔离的关键技术。

我们主要使用 `github.com/stretchr/testify/mock` 库来实现 Mocking。特别是对于实现了 `pkg/component/component.go` 中 `Component` 接口的组件，我们提供了一种标准的 Mocking 方式。


```go
type givenDetail struct {
    // Example of another initial state field for the component being tested
    initialProcessingPrefix string
    // MockOption to set up the behavior of dependencyA
    dependencyAMockSetup component.MockOption
    // ...
}

type whenDetail struct {
    processDataInput string
    // ...
}

type thenExpected struct {
    expectedOutput string
    expectedError error
    // ...
}

tests := []test.Case[givenDetail, whenDetail, thenExpected]{
    {
        Scenario: "Component processes data successfully with mocked dependency",
        Given:    "YourComponent with an initial prefix and dependencyA mocked to return 'related_data_value' for 'input_key'",
        When:     "ProcessData is called with 'input_key'",
        Then:     "Should return 'prefix:input_key:related_data_value' and no error",
        GivenDetail: givenDetail{
            initialProcessingPrefix: "prefix1",
            dependencyAMockSetup: func(m *mock.Mock) {
                // We expect DependencyA's FetchRelatedData to be called with "input_key"
                // and it should return "related_data_value" and no error.
                m.On("FetchRelatedData", "input_key").
                    Return("related_data_value", nil).
                    Once() // Expect it to be called exactly once.
            },
        },
        WhenDetail: whenDetail{
            processDataInput: "input_key",
        },
        ThenExpected: thenExpected{
            expectedOutput: "prefix1:input_key:related_data_value",
            expectedError:  nil,
        },
    },
    // ...更多测试用例...
}


// 在 for _, tt := range tests { t.Run(tt.Scenario, func(t *testing.T) { ... }) } 循环内部

// Given 阶段: Setup mocks and the component under test
var mockHelperForDepA *mock.Mock
defer func() { // 确保在每个子测试结束时断言
    if mockHelperForDepA != nil {
        mockHelperForDepA.AssertExpectations(t)
    }
}()

// 创建并配置 mockDependencyA
// dependency_a_pkg.NewFactory 应该是一个返回 DependencyA 接口和 error 的工厂函数
// 它接受 component.MockOption 来配置其内部的 mock.Mock 对象
mockDependencyA, err := dependency_a_pkg.NewFactory(
    component.MockOption(func(m *mock.Mock) {
        mockHelperForDepA = m // 保存 mock.Mock 实例以供 AssertExpectations 使用
        if tt.GivenDetail.dependencyAMockSetup != nil {
            // 应用测试用例中定义的 specific mock setup
            tt.GivenDetail.dependencyAMockSetup(m)
        }
    }),
).New("mocked_dep_a_instance", nil /* config for dep A */, dependency_a_pkg.Dependencies{})
Expect(err).NotTo(HaveOccurred())
Expect(mockDependencyA).NotTo(BeNil())

// 假设 YourComponent 的构造函数如下：
componentUnderTest := NewYourComponent(tt.GivenDetail.initialProcessingPrefix, mockDependencyA)

// When 阶段: Execute the action being tested
actualOutput, actualErr := componentUnderTest.ProcessData(context.Background(), tt.WhenDetail.processDataInput)

// Then 阶段: Assert the outcomes
if tt.ThenExpected.expectedError != nil {
    Expect(actualErr).To(HaveOccurred())
    Expect(actualErr.Error()).To(Equal(tt.ThenExpected.expectedError.Error()))
} else {
    Expect(actualErr).NotTo(HaveOccurred())
}
Expect(actualOutput).To(Equal(tt.ThenExpected.expectedOutput))
```