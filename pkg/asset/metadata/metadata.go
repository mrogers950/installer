package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

const (
	// MetadataFilename is name of the file where clustermetadata is stored.
	MetadataFilename  = "metadata.json"
	metadataAssetName = "Cluster Metadata"
)

// Metadata depends on cluster and installconfig,
type Metadata struct {
	metadata *types.ClusterMetadata
	file     *asset.File
}

var _ asset.WritableAsset = (*Metadata)(nil)

// Name returns the human-friendly name of the asset.
func (m *Metadata) Name() string {
	return metadataAssetName
}

// Dependencies returns the dependency of the MetaData.
func (m *Metadata) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		&cluster.Cluster{},
	}
}

// Generate generates the metadata.yaml file.
func (m *Metadata) Generate(parents asset.Parents) error {
	installConfig := &installconfig.InstallConfig{}
	cluster := &cluster.Cluster{}
	parents.Get(installConfig, cluster)

	m.metadata = &types.ClusterMetadata{
		ClusterName: installConfig.Config.ObjectMeta.Name,
	}
	switch {
	case installConfig.Config.Platform.AWS != nil:
		m.metadata.ClusterPlatformMetadata.AWS = &types.ClusterAWSPlatformMetadata{
			Region: installConfig.Config.Platform.AWS.Region,
			Identifier: map[string]string{
				"tectonicClusterID": installConfig.Config.ClusterID,
			},
		}
	case installConfig.Config.Platform.OpenStack != nil:
		m.metadata.ClusterPlatformMetadata.OpenStack = &types.ClusterOpenStackPlatformMetadata{
			Region: installConfig.Config.Platform.OpenStack.Region,
			Identifier: map[string]string{
				"tectonicClusterID": installConfig.Config.ClusterID,
			},
		}
	case installConfig.Config.Platform.Libvirt != nil:
		m.metadata.ClusterPlatformMetadata.Libvirt = &types.ClusterLibvirtPlatformMetadata{
			URI: installConfig.Config.Platform.Libvirt.URI,
		}
	default:
		return fmt.Errorf("no known platform")
	}

	data, err := json.Marshal(m.metadata)
	if err != nil {
		return errors.Wrap(err, "failed to Marshal ClusterMetadata")
	}
	m.file = &asset.File{
		Filename: MetadataFilename,
		Data:     data,
	}

	return nil
}

// Files returns the files generated by the asset.
func (m *Metadata) Files() []*asset.File {
	return []*asset.File{m.file}
}
