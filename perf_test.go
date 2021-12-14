package v8go_test

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

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
			_, err := fn.Call(ctx.Global())
			if err != nil {
				panic(err)
			}
		}

		start := time.Now()
		profile := profiler.StopProfiling(profileName)
		// fmt.Println("duration to StopProfiling ", time.Since(start))
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
