package async

type Task struct {
	done   chan bool
	todo   func()
	result interface{}
}

func NewTask(m func()) *Task {
	done := make(chan bool)
	task := Task{}
	task.done = done
	task.todo = func() {
		m()
		done <- true
	}

	return &task
}

func NewResultTask(m func() interface{}) *Task {
	done := make(chan bool)
	task := Task{}
	task.done = done
	task.todo = func() {
		task.result = m()
		done <- true
	}
	return &task
}

func (t *Task) InvokeAsync() {
	go t.todo()
}

func (t *Task) Await() {
	if _, ok := <-t.done; ok {
		close(t.done)
	}
}

func (t *Task) GetResult() interface{} {
	t.Await()
	return t.result
}
