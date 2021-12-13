package v8go_test

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

var total time.Duration

func BenchmarkProfiler(b *testing.B) {
	b.ReportAllocs()

	vm := v8go.NewIsolate()
	profiler := v8go.NewCPUProfiler(vm)

	ctx := v8go.NewContext(vm)
	_, _ = ctx.RunScript(profileScript, "script.js")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()

	for n := 0; n < b.N; n++ {
		profileName := fmt.Sprintf("test-%d", n)
		profiler.StartProfiling(profileName)

		for j := 0; j < 10; j++ {
			_, _ = fn.Call(ctx.Global())
		}

		start := time.Now()
		profile := profiler.StopProfiling(profileName)
		printTree(profile.GetTopDownRoot())
		duration := time.Since(start)
		total += duration

		profile.Delete()
	}

	ctx.Close()
	profiler.Dispose()
	vm.Dispose()

	fmt.Printf("average duration %dus\n", total.Microseconds()/int64(b.N))
}

func TestProfilerPerf(t *testing.T) {
	fmt.Println("starting test")
	// rounds := 1000
	// for n := 0; n < rounds; n++ {
	vm := v8go.NewIsolate()
	ctx := v8go.NewContext(vm)
	profiler := v8go.NewCPUProfiler(vm)

	profileName := fmt.Sprintf("test-%d", 0)

	profiler.StartProfiling(profileName)
	ctx.RunScript(profileScript, "script.js")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()
	fn.Call(ctx.Global())

	start := time.Now()
	profile := profiler.StopProfiling(profileName)
	printTree(profile.GetTopDownRoot())
	fmt.Printf(" duration to generate profile %s\n", time.Since(start))

	profile.Delete()
	profiler.Dispose()
	ctx.Close()
	vm.Dispose()
	// }

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
