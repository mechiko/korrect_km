package templates

type Semaphore interface {
	Acquire()
	Release()
}
type semaphore struct {
	semC chan struct{}
}

func NewSemaphore(maxConcurrency int) Semaphore {
	if maxConcurrency <= 0 {
		panic("maxConcurrency must be positive")
	}
	return &semaphore{
		semC: make(chan struct{}, maxConcurrency),
	}
}
func (s *semaphore) Acquire() {
	s.semC <- struct{}{}
}
func (s *semaphore) Release() {
	<-s.semC
}
