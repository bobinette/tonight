package tonight

func isDone(task Task) bool {
	for i := len(task.Log) - 1; i >= 0; i-- {
		log := task.Log[i]

		if log.Completion == 100 {
			return true
		}
	}

	return false
}

func isWontDo(task Task) bool {
	for i := len(task.Log) - 1; i >= 0; i-- {
		log := task.Log[i]

		switch log.Type {
		case LogTypeWontDo:
			return true
		}
	}

	return false
}

func isPaused(task Task) bool {
	for i := len(task.Log) - 1; i >= 0; i-- {
		log := task.Log[i]

		switch log.Type {
		case LogTypePause:
			return true
		}
	}

	return false
}

func isStarted(task Task) bool {
	for i := len(task.Log) - 1; i >= 0; i-- {
		log := task.Log[i]
		// Change to a Done type
		if log.Completion == 100 {
			return false
		}

		switch log.Type {
		case LogTypePause, LogTypeWontDo:
			return false
		case LogTypeStart:
			return true
		}
	}

	return false
}

func not(f func(Task) bool) func(Task) bool {
	return func(task Task) bool {
		return !f(task)
	}
}

func isTransitionAllowed(task Task, transition LogType) bool {
	stateMachine := map[LogType][]func(Task) bool{
		LogTypeStart: {
			not(isDone),
			not(isWontDo),
			not(isStarted),
		},
		LogTypePause: {
			isStarted,
		},
		LogTypeWontDo: {
			not(isDone),
			not(isWontDo),
		},
		LogTypeProgress: {
			not(isDone),
			not(isWontDo),
		},
		LogTypeComment:  {},
		LogTypePostpone: {},
	}

	for _, check := range stateMachine[transition] {
		if !check(task) {
			return false
		}
	}

	return true
}
