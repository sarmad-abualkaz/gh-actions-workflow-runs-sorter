package util

func ShouldComplete(previousRunStatus string) bool {

    if previousRunStatus != "completed" {
		return false
	}

	return true
}
