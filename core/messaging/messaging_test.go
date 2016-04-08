package messaging

import (
	"testing"
)


var (
	bus *EventBusInterface
	err error
)

func TestMain(m *testing.M) {
	bus = GetInstance()
	m.Run()
}

func TestOne(t *testing.T) {
	if err != nil {
		t.Fatalf("failed to Set MemCache: %v", err)
	}
	/*if err := mct.Set(testPackage); err != nil {
		t.Errorf("Error in Set function: %v", err)
	}*/
}