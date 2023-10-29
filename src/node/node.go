package node

type Runner interface {
	// Run method is responsible for running the coordinator or worker.
	Run()
}

type Coordinator struct {
}

func (c *Coordinator) Run() {

}

type Worker struct {
}

func (w *Worker) Run() {

}
