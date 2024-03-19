package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outArr := make([]Out, 0)
	outFlagArr := make([]bool, 0)
	resultArr := make([]interface{}, 0)
	for task := range in {
		out := make(Bi, 1)
		out <- task
		in = out
		for _, s := range stages {
			in = s(in)
		}
		outArr = append(outArr, in)
		outFlagArr = append(outFlagArr, false)
		resultArr = append(resultArr, nil)
	}
	result := make(Bi, len(resultArr))
	defer close(result)
	for {
		select {
		case _, ok := <-done:
			if !ok {
				return result
			}
		default:
		}
		for i := 0; i < len(outArr); i++ {
			select {
			case x := <-outArr[i]:
				outFlagArr[i] = true
				resultArr[i] = x
			default:
			}
		}
		quit := true
		for _, f := range outFlagArr {
			quit = quit && f
		}
		if quit {
			break
		}
	}
	for _, v := range resultArr {
		result <- v
	}
	return result
}
