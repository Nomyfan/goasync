package goasync

type Any = interface{}

type Awaitable interface {
	Await()
}

type VoidTask interface {
	ContinueWithVoidThenVoid(func()) VoidTask
	ContinueWithVoidThenAny(func() Any) ResultTask
	Await()
}

type ResultTask interface {
	ContinueWithAnyThenVoid(func(Any)) VoidTask
	ContinueWithAnyThenAny(func(Any) Any) ResultTask
	Await()
	Result() Any
}

type Task struct {
	done   chan bool
	todo   func()
	result Any
}

func newVoidTask(m func()) *Task {
	done := make(chan bool)
	return &Task{done: done, todo: func() {
		m()
		done <- true
	}}
}

func NewVoidTask(m func()) VoidTask {
	return newVoidTask(m)
}

func newResultTask(m func() Any) *Task {
	task := Task{done: make(chan bool)}
	task.todo = func() {
		task.result = m()
		task.done <- true
	}
	return &task
}

func NewResultTask(m func() Any) ResultTask {
	return newResultTask(m)
}

func StartNewVoidTask(m func()) VoidTask {
	task := newVoidTask(m)
	task.InvokeAsync()
	return task
}

func StartNewResultTask(m func() Any) ResultTask {
	task := newResultTask(m)
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

func (t *Task) Result() Any {
	t.Await()
	return t.result
}

func (t *Task) ContinueWithVoidThenVoid(next func()) VoidTask {
	return StartNewVoidTask(func() {
		t.Await()
		next()
	})
}

func (t *Task) ContinueWithVoidThenAny(next func() Any) ResultTask {
	return StartNewResultTask(func() Any {
		t.Await()
		return next()
	})
}

func (t *Task) ContinueWithAnyThenVoid(next func(any Any)) VoidTask {
	return StartNewVoidTask(func() {
		next(t.Result())
	})
}

func (t *Task) ContinueWithAnyThenAny(next func(Any) Any) ResultTask {
	return StartNewResultTask(func() Any {
		return next(t.Result())
	})
}