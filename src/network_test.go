package src

// func GetProcessIdxById(id int, processes []*Process) int {
// 	for i, process := range processes {
// 		if process.Id == id {
// 			return i
// 		}
// 	}
// 	return -1
// }
/*
func TestSpawn(t *testing.T) {
	//verbose := true

	var processes []*Process

	// processes, err := Parse("../processes_file.txt")

	// if err != nil {

	// 	t.Fail()

	// }

	processes = append(processes, NewProcess("P1", true))
	processes = append(processes, NewProcess("P2", true))
	processes = append(processes, NewProcess("P3", true))
	processes = append(processes, NewProcess("P4", true))

	network := SpawnNetwork(&processes)
	network.Coordinator = (processes)[3]

	for _, process := range processes {

		assert.Equal(t, 3, len(process.OtherProcesses))

	}

	network.Coordinator.InitCommit(5)

	assert.Equal(t, 0, len(network.Coordinator.History), "no values should be commited on start")

	network.Cycle()
	network.Cycle()
	network.Cycle()

	for _, process := range processes {
		assert.Equal(t, 1, len(process.History), "a value should be commited after the 3rd network cycle")
		assert.Equal(t, 5, process.History[0], "the correct value should be commited to all participants")

		fmt.Println(process.History)

	}

}

// TODO: add "block" noticing via timeout,
// so that network can continue to operate
// even when a node is unreachable

// func TestRecoveryFromTimeFailure(t *testing.T) {

// 	var processes []*Process

// 	processes = append(processes, NewProcess("P1", true))
// 	processes = append(processes, NewProcess("P2", true))
// 	processes = append(processes, NewProcess("P3", true))
// 	processes = append(processes, NewProcess("P4", true))

// 	network := SpawnNetwork(&processes)
// 	network.Coordinator = (processes)[3]

// 	network.Processes["P2"].TimeFailure = true

// 	network.Coordinator.InitCommit(5)

// 	network.Cycle()
// 	network.Cycle()
// 	network.Cycle()

// 	for _, process := range processes {
// 		fmt.Println(process.Name)
// 		if process.Name != "P2" {
// 			assert.Equal(t, 1, len(process.History), "the correct value should be commited to non-failed participants")
// 		} else {
// 			assert.Equal(t, 0, len(process.History), "currently failed process shouldn't have its history updated")
// 		}

// 		fmt.Println(process.History)

// 	}

// }


 */