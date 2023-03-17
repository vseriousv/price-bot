package telegram

type BlackList map[int64]bool

func getBlackList() BlackList {
	var bl BlackList = make(map[int64]bool)
	return bl
}
