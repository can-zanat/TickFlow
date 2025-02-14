package observer

// Observer interface: create Update method for observer pattern.
type Observer interface {
	Update(data map[string]interface{})
}
