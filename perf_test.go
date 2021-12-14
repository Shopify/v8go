package v8go_test

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
	v8 "rogchap.com/v8go"
)

var fib = `function fibonacci(num) {
  if (num <= 1) return 1;

  return fibonacci(num - 1) + fibonacci(num - 2);
}`

func BenchmarkProfiler(b *testing.B) {
	vm := v8go.NewIsolate()
	profiler := v8go.NewCPUProfiler(vm)
	ctx := v8go.NewContext(vm)
	ctx.RunScript(fib, "script.js")
	input, _ := v8.NewValue(vm, int64(500))
	val, _ := ctx.Global().Get("fibonacci")
	fn, _ := val.AsFunction()

	profiler.StartProfiling("outside-profile")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		profileName := fmt.Sprintf("test-%d", n)

		profiler.StartProfiling(profileName)

		fn.Call(ctx.Global(), input)

		profile := profiler.StopProfiling(profileName)
		// start := time.Now()
		printTree(profile.GetTopDownRoot())
		// fmt.Println("time printing tree ", time.Since(start))
	}

	profiler.StopProfiling("outside-profile")
	// profile.Delete()
	ctx.Close()
	profiler.Dispose()
	vm.Dispose()
}

func TestProfilerPerf(t *testing.T) {
	fmt.Println("starting test")
	vm := v8go.NewIsolate()
	profiler := v8go.NewCPUProfiler(vm)
	ctx := v8go.NewContext(vm)
	ctx.RunScript(profileScript, "script.js")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()

	rounds := 1000
	var total time.Duration
	for n := 0; n < rounds; n++ {
		profileName := fmt.Sprintf("test-%d", n)

		profiler.StartProfiling(profileName)

		for j := 0; j < 10; j++ {
			fn.Call(ctx.Global())
		}

		start := time.Now()
		profile := profiler.StopProfiling(profileName)
		printTree(profile.GetTopDownRoot())
		total += time.Since(start)
	}

	// profile.Delete()
	ctx.Close()
	profiler.Dispose()
	vm.Dispose()

	fmt.Printf("average duration to generate profile %s\n", time.Duration(total.Nanoseconds()/int64(rounds)))
}

// Note that to focus on exercising each function, we skip printing
// but printing is useful to just double-check the tree looks correct.
func printTree(node *v8go.CPUProfileNode) {
	// fmt.Printf("%s %s:%d:%d\n", node.GetFunctionName(), node.GetScriptResourceName(), node.GetLineNumber(), node.GetColumnNumber())
	node.GetFunctionName()
	node.GetLineNumber()
	node.GetColumnNumber()
	node.GetScriptResourceName()
	count := node.GetChildrenCount()
	for i := 0; i < count; i++ {
		printTree(node.GetChild(i))
	}
}
