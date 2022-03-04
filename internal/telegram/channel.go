package telegram

type Channel struct {
	data chan string
	stop chan struct{}
}

type userChannel map[int64]*Channel

func makeChannel() *Channel {
	return &Channel{
		make(chan string, 10),
		make(chan struct{}),
	}
}
