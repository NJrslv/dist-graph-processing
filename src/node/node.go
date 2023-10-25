package node

type Node interface {
	/*
		Run starts the coordinator or worker component.
		This method is responsible for running the coordinator or worker.
	*/
	Run()
}
