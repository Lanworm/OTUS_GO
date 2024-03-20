package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func exec(done In, intStream In, stage Stage) {
	resultStream := make(Bi)
	go func() {
		defer close(resultStream)
		for range intStream {
			select {
			case <-done:
				return
			case resultStream <- stage(intStream):
			}
		}
	}()
	intStream = resultStream
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := stages[0](in)
	for _, s := range stages {
		exec(done, out, s)
	}
	return out
}
