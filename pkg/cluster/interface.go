package cloud

type CloudCluster interface {
	// CreateCluster creates a new cloud cluster
	CreateCluster() error
	// // DeleteCluster deletes a cloud cluster
	// DeleteCluster() error
	// // GetClusterStatus returns the status of a cloud cluster
	// GetClusterStatus() (string, error)
	// // GetClusterConfig returns the config of a cloud cluster
	// GetClusterConfig() (string, error)
	// // GetClusterKubeConfig returns the kubeconfig of a cloud cluster
	// GetClusterKubeConfig() (string, error)
}
