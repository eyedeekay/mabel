package i2pini

import (
	"fmt"

	"github.com/eyedeekay/sam-forwarder/interface"
	"github.com/go-ini/ini"
)

import (
	"github.com/eyedeekay/eephttpd"
	"github.com/eyedeekay/httptunnel"
	"github.com/eyedeekay/httptunnel/multiproxy"
	"github.com/eyedeekay/reposam"
	"github.com/eyedeekay/sam-forwarder/tcp"
	"github.com/eyedeekay/sam-forwarder/udp"
	"github.com/eyedeekay/samtracker"
)

func SAMTunnel(sct *ini.Section) (samtunnel.SAMTunnel, error) {
	k, err := sct.GetKey("type")
	if err != nil {
		return nil, err
	}
	var rt samtunnel.SAMTunnel
	switch k.Value() {
	case "server":
		rt = new(samforwarder.SAMForwarder)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "http":
		rt = new(samforwarder.SAMForwarder)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "client":
		rt = new(samforwarder.SAMClientForwarder)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "httpclient":
		rt = new(i2phttpproxy.SAMHTTPProxy)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "browserclient":
		rt = new(i2pbrowserproxy.SAMMultiProxy)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "udpserver":
		rt = new(samforwarderudp.SAMSSUForwarder)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "udpclient":
		rt = new(samforwarderudp.SAMSSUClientForwarder)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "eephttpd":
		rt = new(eephttpd.EepHttpd)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "debrepo":
		rt = new(reposam.RepoSam)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	case "tracker":
		rt = new(samtracker.SamTracker)
		if err := sct.MapTo(rt); err != nil {
			return nil, fmt.Errorf("Failed to load tunnel, %s %s", sct.Name(), err.Error())
		}
	}
	return rt, nil
}

func SAMTunnelSlice(inifile string) ([]samtunnel.SAMTunnel, error) {
	var sts []samtunnel.SAMTunnel
	cfg, err := ini.Load(inifile)
	if err != nil {
		return nil, err
	}
	for _, sct := range cfg.Sections() {
		if sct.HasKey("type") {
			st, err := SAMTunnel(sct)
			if err != nil {
				sts = append(sts, st)
			}
		} else {
			return nil, fmt.Errorf("error processing section %s:\n  type field is required")
		}
	}
	return sts, nil
}
