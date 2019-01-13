package nomadConnector

type NomadCfg struct {
}

// NomadConnector defines the interface of the component being able to communicate with nomad
type NomadConnector interface {
	ScaleBy(jobName string, amount int) error
}

// New creates a new nomad connector
func (cfg *NomadCfg) New() NomadConnector {
	nc := &nomadConnectorImpl{}

	return nc
}
