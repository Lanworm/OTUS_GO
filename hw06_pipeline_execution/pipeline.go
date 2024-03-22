package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func exec(done In, in In, out Out) {
	go func() {
		for {
			select {
			case <-in:
				in = out
			case <-done:
				return
			}
		}
	}()
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Out)
	for _, s := range stages {
		out = s(in)
		exec(done, out, in)
	}
	return out
}
