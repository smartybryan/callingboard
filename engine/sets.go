package engine

func SetDifference(mainSet, subtractSet []MemberName) (names []MemberName) {
	for _, name := range mainSet {
		if inSet(subtractSet, name) {
			continue
		}
		names = append(names, name)
	}
	return names
}

func inSet(set []MemberName, name MemberName) bool {
	for _, value := range set {
		if value ==  name {
			return true
		}
	}
	return false
}
