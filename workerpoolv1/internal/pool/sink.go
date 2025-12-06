package pool

import "fmt"

func PrintSink(r Result) {
	if r.Err != "" {
		fmt.Printf("[ERROR] job=%03d cost=%v err=%s\n", r.JobID, r.Cost, r.Err)
		return
	}

	fmt.Printf("[OK] job=%03d cost=%v\n", r.JobID, r.Cost)
}
