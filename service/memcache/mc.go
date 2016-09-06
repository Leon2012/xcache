package memcache

const (
	ERROR        int = 1
	CLIENT_ERROR     = 2
	SERVER_ERROR     = 3
)

type MCError struct {
	msg  string
	code int
}

func NewMCError(code int, msg string) MCError {
	return MCError{
		msg:  msg,
		code: code,
	}
}

func (e MCError) Error() string {
	var err string
	switch e.code {
	case ERROR:
		err = "ERROR\r\n"
		break
	case CLIENT_ERROR:
		err = "CLIENT_ERROR " + e.msg + " \r\n"
		break
	case SERVER_ERROR:
		err = "SERVER_ERROR " + e.msg + " \r\n"
		break
	}
	return err
}
