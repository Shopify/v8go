package v8go_test

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestProfilerPerf(t *testing.T) {
	fmt.Println("starting test")
	// rounds := 1000
	// for n := 0; n < rounds; n++ {
	vm, _ := v8go.NewIsolate()
	ctx, _ := v8go.NewContext(vm)
	profiler := v8go.NewCpuProfiler(vm)

	profileName := fmt.Sprintf("test-%d", 0)

	profiler.StartProfiling(profileName)
	ctx.RunScript(profileScript, "script.js")
	val, _ := ctx.Global().Get("start")
	fn, _ := val.AsFunction()
	fn.Call(ctx.Global())

	start := time.Now()
	profile := profiler.StopProfiling(profileName, "")
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
func printTree(node *v8go.CpuProfileNode) {
	fmt.Printf("%s %s:%d:%d\n", node.GetFunctionName(), node.GetScriptResourceName(), node.GetLineNumber(), node.GetColumnNumber())
	node.GetFunctionName()
	node.GetLineNumber()
	node.GetColumnNumber()
	node.GetScriptResourceName()
	count := node.GetChildrenCount()
	for i := 0; i < count; i++ {
		printTree(node.GetChild(i))
	}
}
