package common

const (
	PING string = "PING"
	PONG string = "PONG"
)

type NetConf struct {
	Name     string
	Ip       string
	Port     string
	Protocol string
}

type Ping struct {
	Data string
}

type Pong struct {
	Data string
}

type PingHandler struct{}

func (p *PingHandler) Ping(req *Ping, res *Pong) error {
	if req.Data == PING {
		res.Data = PONG
	}

	return nil
}

type SubscribeRequest struct {
	Subscriber NetConf
	Publisher  NetConf
}

type SubscribeResponse struct {
	Ok bool
}
