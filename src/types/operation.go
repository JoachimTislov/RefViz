package types

type Operation []struct {
	Condition bool
	Action    func() error
	Msg       string
}
