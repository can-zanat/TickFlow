// subject_test.go
package observer

import (
	"testing"

	"go.uber.org/mock/gomock"
)

func TestSubjectAttachDetachNotify(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockObserver1 := NewMockObserver(ctrl)
	mockObserver2 := NewMockObserver(ctrl)

	subject := NewSubject()

	subject.Attach(mockObserver1)
	subject.Attach(mockObserver2)

	data := map[string]interface{}{"key": "value"}

	mockObserver1.EXPECT().Update(data).Times(1)
	mockObserver2.EXPECT().Update(data).Times(1)

	subject.Notify(data)

	subject.Detach(mockObserver1)

	newData := map[string]interface{}{"anotherKey": "anotherValue"}

	mockObserver2.EXPECT().Update(newData).Times(1)

	subject.Notify(newData)
}
