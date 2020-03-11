package tc

import (
	"github.com/eyedeekay/sam-forwarder/interface"
)

type TunnelController struct {
	samtunnel.SAMTunnel
	group string
}

func (tc *TunnelController) SetGroup(group string) {
	tc.group = group
}

func (tc *TunnelController) GetGroup() string {
	temp := tc.group
	tc.group = ""
	return temp
}

func NewTunnelController(tunnel samtunnel.SAMTunnel) (*TunnelController, error) {
	tc := TunnelController{
		SAMTunnel: tunnel,
	}
	return &tc, nil
}
