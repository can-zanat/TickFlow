package observer

// Subject manage observer and notify them.
type Subject struct {
	observers []Observer
}

// NewSubject create new subject.
func NewSubject() *Subject {
	return &Subject{
		observers: make([]Observer, 0),
	}
}

// Attach put subject in an observer.
func (s *Subject) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// Detach removes an observer from the subject.
func (s *Subject) Detach(o Observer) {
	for i, observer := range s.observers {
		if observer == o {
			s.observers = append(s.observers[:i], s.observers[i+1:]...)
			break
		}
	}
}

// Notify calls update method of all the observers.
func (s *Subject) Notify(data map[string]interface{}) {
	for _, observer := range s.observers {
		observer.Update(data)
	}
}
