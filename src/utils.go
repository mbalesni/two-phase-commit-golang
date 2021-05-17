package src

func MapCopy(originalMap map[string]*Process) map[string]*Process {
	newMap := make(map[string]*Process)

	for k, v := range originalMap {
		newMap[k] = v
	}

	return newMap
}
