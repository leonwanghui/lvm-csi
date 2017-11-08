package lvm

import (
	"errors"
	"log"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
//                            Controller Service                              //
////////////////////////////////////////////////////////////////////////////////

// CreateVolume implementation
func (p *Plugin) CreateVolume(
	ctx context.Context,
	req *csi.CreateVolumeRequest) (
	*csi.CreateVolumeResponse, error) {

	log.Println("start to CreateVolume, req is:", req)
	defer log.Println("end to CreateVolume")

	if req.CapacityRange == nil {
		log.Println("Capacity required!")
		return nil, errors.New("Capacity required!")
	}

	v, err := lvmDriver.CreateVolume(req.Name, req.CapacityRange.RequiredBytes)
	if err != nil {
		return nil, err
	}

	return &csi.CreateVolumeResponse{
		Reply: &csi.CreateVolumeResponse_Result_{
			Result: &csi.CreateVolumeResponse_Result{
				VolumeInfo: v,
			},
		},
	}, nil
}

// DeleteVolume implementation
func (p *Plugin) DeleteVolume(
	ctx context.Context,
	req *csi.DeleteVolumeRequest) (
	*csi.DeleteVolumeResponse, error) {

	log.Println("start to DeleteVolume")
	defer log.Println("end to DeleteVolume")

	lvPath, ok := req.UserCredentials.Data["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume metadata!")
		log.Println(err)
		return nil, err
	}

	if err := lvmDriver.DeleteVolume(lvPath); err != nil {
		log.Println("Failed to delete volume in lvm driver!")
		return nil, err
	}

	return &csi.DeleteVolumeResponse{
		Reply: &csi.DeleteVolumeResponse_Result_{
			Result: &csi.DeleteVolumeResponse_Result{},
		},
	}, nil
}

// ControllerPublishVolume implementation
func (p *Plugin) ControllerPublishVolume(
	ctx context.Context,
	req *csi.ControllerPublishVolumeRequest) (
	*csi.ControllerPublishVolumeResponse, error) {

	log.Println("start to ControllerPublishVolume")
	defer log.Println("end to ControllerPublishVolume")

	lvPath, ok := req.VolumeAttributes["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume metadata!")
		log.Println(err)
		return nil, err
	}
	initiator, ok := req.VolumeAttributes["initiator"]
	if !ok {
		err := errors.New("Failed to find initiator field in volume metadata!")
		log.Println(err)
		return nil, err
	}

	connInfo, err := lvmDriver.InitializeConnection(lvPath, initiator)
	if err != nil {
		log.Println("Failed to initialize volume connection in lvm driver!")
		return nil, err
	}

	return &csi.ControllerPublishVolumeResponse{
		Reply: &csi.ControllerPublishVolumeResponse_Result_{
			Result: &csi.ControllerPublishVolumeResponse_Result{
				PublishVolumeInfo: connInfo,
			},
		},
	}, nil
}

// ControllerUnpublishVolume implementation
func (p *Plugin) ControllerUnpublishVolume(
	ctx context.Context,
	req *csi.ControllerUnpublishVolumeRequest) (
	*csi.ControllerUnpublishVolumeResponse, error) {

	log.Println("start to ControllerUnpublishVolume")
	defer log.Println("end to ControllerUnpublishVolume")

	lvPath, ok := req.UserCredentials.Data["lvPath"]
	if !ok {
		err := errors.New("Failed to find logic volume path in volume metadata!")
		log.Println(err)
		return nil, err
	}

	if err := lvmDriver.TerminateConnection(lvPath, ""); err != nil {
		log.Println("Failed to terminate volume connection in lvm driver!")
		return nil, err
	}

	return &csi.ControllerUnpublishVolumeResponse{
		Reply: &csi.ControllerUnpublishVolumeResponse_Result_{
			Result: &csi.ControllerUnpublishVolumeResponse_Result{},
		},
	}, nil
}

// ValidateVolumeCapabilities implementation
func (p *Plugin) ValidateVolumeCapabilities(
	ctx context.Context,
	req *csi.ValidateVolumeCapabilitiesRequest) (
	*csi.ValidateVolumeCapabilitiesResponse, error) {

	log.Println("start to ValidateVolumeCapabilities")
	defer log.Println("end to ValidateVolumeCapabilities")

	if strings.TrimSpace(req.VolumeId) == "" {
		return &csi.ValidateVolumeCapabilitiesResponse{
			Reply: &csi.ValidateVolumeCapabilitiesResponse_Error{
				Error: &csi.Error{
					Value: &csi.Error_ValidateVolumeCapabilitiesError_{
						ValidateVolumeCapabilitiesError: &csi.Error_ValidateVolumeCapabilitiesError{
							ErrorCode:        csi.Error_ValidateVolumeCapabilitiesError_INVALID_VOLUME_INFO,
							ErrorDescription: "invalid volume id",
						},
					},
				},
			},
		}, nil
	}

	for _, capabilities := range req.VolumeCapabilities {
		if capabilities.GetMount() != nil {
			return &csi.ValidateVolumeCapabilitiesResponse{
				Reply: &csi.ValidateVolumeCapabilitiesResponse_Result_{
					Result: &csi.ValidateVolumeCapabilitiesResponse_Result{
						Supported: false,
						Message:   "opensds does not support mounted volume",
					},
				},
			}, nil
		}
	}

	return &csi.ValidateVolumeCapabilitiesResponse{
		Reply: &csi.ValidateVolumeCapabilitiesResponse_Result_{
			Result: &csi.ValidateVolumeCapabilitiesResponse_Result{
				Supported: true,
				Message:   "supported",
			},
		},
	}, nil
}

// ListVolumes implementation
func (p *Plugin) ListVolumes(
	ctx context.Context,
	req *csi.ListVolumesRequest) (
	*csi.ListVolumesResponse, error) {

	log.Println("start to ListVolumes")
	defer log.Println("end to ListVolumes")

	return nil, errors.New("Not implemented!")
}

// GetCapacity implementation
func (p *Plugin) GetCapacity(
	ctx context.Context,
	req *csi.GetCapacityRequest) (
	*csi.GetCapacityResponse, error) {

	log.Println("start to GetCapacity")
	defer log.Println("end to GetCapacity")

	return nil, errors.New("Not implemented!")
}

// ControllerProbe implementation
func (p *Plugin) ControllerProbe(
	ctx context.Context,
	req *csi.ControllerProbeRequest) (
	*csi.ControllerProbeResponse, error) {

	log.Println("start to ControllerProbe")
	defer log.Println("end to ControllerProbe")

	return nil, errors.New("Not implemented!")
}

// ControllerGetCapabilities implementation
func (p *Plugin) ControllerGetCapabilities(
	ctx context.Context,
	req *csi.ControllerGetCapabilitiesRequest) (
	*csi.ControllerGetCapabilitiesResponse, error) {

	log.Println("start to ControllerGetCapabilities")
	defer log.Println("end to ControllerGetCapabilities")

	return &csi.ControllerGetCapabilitiesResponse{
		Reply: &csi.ControllerGetCapabilitiesResponse_Result_{
			Result: &csi.ControllerGetCapabilitiesResponse_Result{
				Capabilities: []*csi.ControllerServiceCapability{
					&csi.ControllerServiceCapability{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
							},
						},
					},
					&csi.ControllerServiceCapability{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
							},
						},
					},
					&csi.ControllerServiceCapability{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_LIST_VOLUMES,
							},
						},
					},
					&csi.ControllerServiceCapability{
						Type: &csi.ControllerServiceCapability_Rpc{
							Rpc: &csi.ControllerServiceCapability_RPC{
								Type: csi.ControllerServiceCapability_RPC_GET_CAPACITY,
							},
						},
					},
				},
			},
		},
	}, nil
}
