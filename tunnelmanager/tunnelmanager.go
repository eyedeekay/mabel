package tm

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"

	"github.com/eyedeekay/mabel/tunnelcontroller"
	"github.com/eyedeekay/sam-forwarder/interface"
)

type TunnelManager struct {
	net.Listener
	*rpc.Server

	Name     string
	Tunnels  map[string]*tc.TunnelController
	Managers map[string]*TunnelManager
}

//List all tunnels in the manager but not any sub-managers
func (tm *TunnelManager) List() []string {
	var r []string
	for _, v := range tm.Tunnels {
		r = append(r, v.SAMTunnel.ID())
	}
	return r
}

func (tm *TunnelManager) ListAllBelow() []string {
	allgroups := tm.AllGroups()
	var r []string
	for _, v := range allgroups {
		r = append(r, v.List()...)
	}
	return r
}

//Groups lists all sub-management groups but not any of their subgroups
func (tm *TunnelManager) Groups() []*TunnelManager {
	var r []*TunnelManager
	for k, v := range tm.Managers {
		v.Name = k
		r = append(r, v)
	}
	return r
}

func Groups(tm []*TunnelManager) []*TunnelManager {
	var r []*TunnelManager
	for _, m := range tm {
		for k, v := range m.Managers {
			v.Name = k
			r = append(r, v)
		}
	}
	return r
}

//List all tunnels in a chosen subgroup
func (tm *TunnelManager) ListGroup(group string) []string {
	var r []string
	for _, v := range tm.Tunnels {
		r = append(r, v.SAMTunnel.ID())
	}
	return r
}

//GroupGroups lists all subgroups in a chosen subgroup
func (tm *TunnelManager) GroupGroups(group string) []*TunnelManager {
	var r []*TunnelManager
	if tm.Managers[group] == nil {
		return nil
	}
	for k, v := range tm.Managers[group].Managers {
		v.Name = k
		r = append(r, v)
	}
	return r
}

func (tm *TunnelManager) AllGroups() []*TunnelManager {
	var list []*TunnelManager
	list = append(list, tm.Groups()...)
	lastlist := -1
	for {
		list = append(list, Groups(list)...)
		if len(list) == lastlist {
			break
		}
		lastlist = len(list)
	}
	return list
}

//Find a tunnel by ID
func (tm *TunnelManager) Find(id string) (*tc.TunnelController, string, error) {
	for _, v := range tm.Tunnels {
		if v.ID() == id {
			return v, "", nil
		}
	}
	for group, v := range tm.Managers {
		for _, w := range v.Tunnels {
			if w.ID() == id {
				return w, group, nil
			}
		}
	}
	return nil, "", fmt.Errorf("Tunnel not found, did you\n\tenter the correct ID or\n\tmean to search instead?")
}

//FindAType of Tunnels, all tunnels of a type
func (tm *TunnelManager) FindAType(kind string) ([]*tc.TunnelController, error) {
	var ts []*tc.TunnelController
	for _, v := range tm.Tunnels {
		if v.GetType() == kind {
			ts = append(ts, v)
		}
	}
	for group, v := range tm.Managers {
		for _, w := range v.Tunnels {
			if w.GetType() == kind {
				ts = append(ts, w)
				ts[len(ts)-1].SetGroup(group)
			}
		}
	}
	if len(ts) > 0 {
		return ts, nil
	}
	return nil, fmt.Errorf("Tunnel not found, did you\n\tenter the correct ID or\n\tmean to search instead?")
}

//Move a tunnel from one organizational group to another
func (tm *TunnelManager) Move(id, group string) error {
	from, fromgroup, err := tm.Find(id)
	if err != nil {
		return err
	}
	if tm.Managers[group] == nil {
		tm.Managers[group] = &TunnelManager{
			Listener: tm.Listener,
			Server:   rpc.NewServer(),
			Tunnels:  make(map[string]*tc.TunnelController),
		}
	}
	tm.Managers[group].Tunnels[from.ID()] = from
	if fromgroup != "" {
		tm.Managers[fromgroup] = nil
	}
	return nil
}

//InitializeTunnelManager, makes sure listeners and memory and stuff are ready.
func InitializeTunnelManager(host string, port int) (*TunnelManager, error) {
	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return &TunnelManager{
		Listener: listener,
		Server:   rpc.NewServer(),
		Tunnels:  make(map[string]*tc.TunnelController),
	}, nil
}

//NewTunnelManagerFromMap sets up a tunnel manager from a map of SAMTunnels.
//Most will use NewTunnelManager anyway, slices will be easier for most people.
func NewTunnelManagerFromMap(host string, port int, tunnels map[string]samtunnel.SAMTunnel) (*TunnelManager, error) {
	tm, err := InitializeTunnelManager(host, port)
	if err != nil {
		return nil, err
	}
	for key, value := range tunnels {
		tc, err := tc.NewTunnelController(value)
		if err != nil {
			return nil, err
		}
		tm.Tunnels[key] = tc
		port += 1
	}
	tm.Register(tm)
	tm.HandleHTTP("/rpc/", "/debug/rpc/")
	return tm, nil
}

//NewTunnelManager sets up a tunnel manager from a slice of SAMTunnels
func NewTunnelManager(host string, port int, tunnels []samtunnel.SAMTunnel) (*TunnelManager, error) {
	tm, err := InitializeTunnelManager(host, port)
	if err != nil {
		return nil, err
	}
	for _, value := range tunnels {
		tc, err := tc.NewTunnelController(value)
		if err != nil {
			return nil, err
		}
		tm.Tunnels[value.ID()] = tc
		port += 1
	}
	tm.Register(tm)
	tm.HandleHTTP("/rpc/", "/debug/rpc/")
	return tm, nil
}

//Serve up the whole mammajamma
func (tm *TunnelManager) Serve() error {
	for {
		if conn, err := tm.Listener.Accept(); err != nil {
			return err
		} else {
			log.Println("new connection established")
			go tm.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}
