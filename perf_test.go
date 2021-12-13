package v8go_test

import (
	"fmt"
	"testing"
	"time"

	"rogchap.com/v8go"
)

func TestProfilerPerf(t *testing.T) {
	var durations []time.Duration
	for n := 0; n < 1000; n++ {
		vm, _ := v8go.NewIsolate()
		ctx, _ := v8go.NewContext(vm)
		profiler := v8go.NewCpuProfiler(vm)

		profileName := fmt.Sprintf("test-%d", n)

		profiler.StartProfiling(profileName)
		for i := 0; i < 10; i++ {
			_, err := ctx.RunScript(profileScript, "script.js")
			if err != nil {
				panic(err)
			}
			val, err := ctx.Global().Get("start")
			if err != nil {
				panic(err)
			}
			fn, err := val.AsFunction()
			if err != nil {
				panic(err)
			}
			_, err = fn.Call(ctx.Global())
			if err != nil {
				panic(err)
			}
		}

		start := time.Now()
		profile := profiler.StopProfiling(profileName, "")
		printTree("", profile.GetTopDownRoot())
		durations = append(durations, time.Since(start))

		profiler.Dispose()
		ctx.Close()
		vm.Dispose()
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	fmt.Printf("average duration %s\n", time.Duration(total.Nanoseconds()/int64(len(durations))))
}

// Note that to focus on exercising each function, we skip printing
// but printing is useful to just double-check the tree looks correct.
func printTree(nest string, node *v8go.CpuProfileNode) {
	// fmt.Printf("%s%s %d:%d\n", nest, node.GetFunctionName(), node.GetLineNumber(), node.GetColumnNumber())
	node.GetFunctionName()
	node.GetLineNumber()
	node.GetColumnNumber()
	count := node.GetChildrenCount()
	if count == 0 {
		return
	}
	// nest = fmt.Sprintf("%s  ", nest)
	for i := 0; i < count; i++ {
		printTree(nest, node.GetChild(i))
	}
}
