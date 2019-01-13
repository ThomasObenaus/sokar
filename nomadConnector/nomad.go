package nomadConnector

type connectorImpl struct {
	jobName string
}

func (nc *connectorImpl) ScaleBy(amount int) error {
	return nil
}
