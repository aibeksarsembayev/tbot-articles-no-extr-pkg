package articledigest

type DigestSender interface {
	Send()
	SendByWeekday()
}

// TODO: add sender as function to be added into any channel or chat
