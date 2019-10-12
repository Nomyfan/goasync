package goasync

type Task struct {
	done   chan bool
	todo   func()
	result interface{}
}

func NewTask(m func()) *Task {
	done := make(chan bool)
	return &Task{done: done, todo: func() {
		m()
		done <- true
	}}
}

func NewResultTask(m func() interface{}) *Task {
	done := make(chan bool)
	task := Task{done: done}
	task.todo = func() {
		task.result = m()
		done <- true
	}
	return &task
}

func StartNew(m func()) *Task {
	task := NewTask(m)
	task.InvokeAsync()
	return task
}

func StartNewResult(m func() interface{}) *Task {
	task := NewResultTask(m)
	task.InvokeAsync()
	return task
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

func (t *Task) ContinueWith(next func(t *Task)) *Task {

	return StartNew(func() {
		t.Await()
		next(t)
	})
}

func (t *Task) ContinueWithResult(next func(t *Task) interface{}) *Task {

	return StartNewResult(func() interface{} {
		t.Await()
		return next(t)
	})
}
