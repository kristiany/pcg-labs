package utils

// Oh this fudging language, you have to do everything yourself
// https://play.golang.org/p/tDdutH672-
type IntSet struct {
	set map[int]bool
}

func NewIntSet() *IntSet {
	return &IntSet{make(map[int]bool)}
}

func (set *IntSet) Add(i int) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found	//False if it existed already
}

func (set *IntSet) Contains(i int) bool {
	_, found := set.set[i]
	return found	//true if it existed already
}
