package common

import "fmt"

type MicronetReconnectTimeoutError struct {
	NetConf
}

func (e MicronetReconnectTimeoutError) Error() string {
	return fmt.Sprintf("connexion timeout to %s:%s", e.Ip, e.Port)
}
