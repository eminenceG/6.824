package mapreduce

import "fmt"
import "sync"

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//
	var done sync.WaitGroup
	for i, file := range mapFiles {
		done.Add(1)
		go func(file string, taskNumber int) {
			// Tells a worker to excute a task by sending a Worker.DoTask RPC to the 
			// worker.

			args := DoTaskArgs{
				jobName,
				file,
				phase,
				taskNumber, 
				n_other,
			}

			finished := false 
			worker := ""
			for !finished {
	   	  worker = <- registerChan
			  finished = call(worker, "Worker.DoTask", args, nil)
			}
			done.Done()
			registerChan <- worker
			// re-register the worker after the task is finished
		}(file, i)
	}
	// waits until all tasks have completed  
	done.Wait()
	fmt.Printf("Schedule: %v done\n", phase)
}
