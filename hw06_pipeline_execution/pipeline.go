package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func exec(in In, done In) Out {
	out := make(Bi)
	go func() {
		for {
			select {
			case <-in:
				v := <-in
				out <- v
			case <-done:
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(In)
	for _, s := range stages {
		out = exec(in, done)
		in = s(out)
	}
	return out
}
