package lvm

import (
	"log"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
//                            Node Service                                    //
////////////////////////////////////////////////////////////////////////////////

// NodePublishVolume implementation
func (p *Plugin) NodePublishVolume(
	ctx context.Context,
	req *csi.NodePublishVolumeRequest) (
	*csi.NodePublishVolumeResponse, error) {
	// TODO
	return nil, nil
}

// NodeUnpublishVolume implementation
func (p *Plugin) NodeUnpublishVolume(
	ctx context.Context,
	req *csi.NodeUnpublishVolumeRequest) (
	*csi.NodeUnpublishVolumeResponse, error) {
	// TODO
	return nil, nil
}

// GetNodeID implementation
func (p *Plugin) GetNodeID(
	ctx context.Context,
	req *csi.GetNodeIDRequest) (
	*csi.GetNodeIDResponse, error) {

	log.Println("start to GetNodeID")
	defer log.Println("end to GetNodeID")

	// TODO
	return nil, nil
}

// NodeProbe implementation
func (p *Plugin) NodeProbe(
	ctx context.Context,
	req *csi.NodeProbeRequest) (
	*csi.NodeProbeResponse, error) {

	log.Println("start to ProbeNode")
	defer log.Println("end to ProbeNode")

	// TODO
	return nil, nil
}

// NodeGetCapabilities implementation
func (p *Plugin) NodeGetCapabilities(
	ctx context.Context,
	req *csi.NodeGetCapabilitiesRequest) (
	*csi.NodeGetCapabilitiesResponse, error) {

	log.Println("start to NodeGetCapabilities")
	defer log.Println("end to NodeGetCapabilities")

	// TODO
	return nil, nil
}
