package chdfs

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/glog"
	csicommon "github.com/kubernetes-csi/drivers/pkg/csi-common"
	"k8s.io/utils/mount"
)

const version = "v1.0.0"

type driver struct {
	csiDriver *csicommon.CSIDriver
	endpoint  string
}

// NewDriver creates a new CSI driver for CHDFS.
func NewDriver(endpoint, driverName, nodeID string) *driver {
	glog.Infof("Driver: %v version: %v", driverName, version)

	csiDriver := csicommon.NewCSIDriver(driverName, version, nodeID)
	csiDriver.AddVolumeCapabilityAccessModes([]csi.VolumeCapability_AccessMode_Mode{
		csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
	})

	return &driver{
		csiDriver: csiDriver,
		endpoint:  endpoint,
	}
}

func NewNodeServer(driver *csicommon.CSIDriver) csi.NodeServer {
	return &nodeServer{
		DefaultNodeServer: csicommon.NewDefaultNodeServer(driver),
		mounter:           mount.New(""),
	}
}

func (d *driver) Start() {
	server := csicommon.NewNonBlockingGRPCServer()
	server.Start(d.endpoint, csicommon.NewDefaultIdentityServer(d.csiDriver), nil,
		NewNodeServer(d.csiDriver))
	server.Wait()
}
