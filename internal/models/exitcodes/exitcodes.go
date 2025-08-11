package exitcodes

const (
	ExitFatal       = 1  // Error happened
	ExitOK          = 0  // normal shutdown
	ExitUpdateReady = 70 // server staged new binary and requests swap+restart
)
