package core

func batchHasTest(batch, projectID, branch, target, browser string) bool {
	res, err := store.GetResultsByBatch(batch)

	if err != nil {
		return false
	}

	for i := 0; i < len(res); i++ {
		if res[i].ProjectID == projectID && res[i].Branch == branch && res[i].Target == target && res[i].Browser == browser {
			return true
		}
	}
	return false
}

func batchIsOld(batch string) bool {
	//TODO
	return false
}

func batchHasDifferentBranch(batch, branch string) bool {
	res, err := store.GetResultsByBatch(batch)

	if err != nil {
		return false
	}

	for i := 0; i < len(res); i++ {
		if res[i].Branch != branch {
			return true
		}
	}

	return false
}
