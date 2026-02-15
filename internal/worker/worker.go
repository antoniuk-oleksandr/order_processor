package worker

// Worker represents a worker that processes tasks from a channel.
type Worker interface {
	// Run starts the worker and processes tasks until the channel is closed.
	Run()
}

type worker struct {
	tasks chan Task
}

// NewWorker creates a new Worker that processes tasks from the given channel.
func NewWorker(tasks chan Task) Worker {
	return &worker{
		tasks: tasks,
	}
}

func (w *worker) Run() {
	for task := range w.tasks {
		task.Process()
	}
}
