package nomadworker

// SetJobCount will scale the nomad workers to the desired count (amount of instances)
func (c *Connector) SetJobCount(datacenter string, count uint) error {
	c.log.Warn().Msgf("nomadworker.Connector.SetJobCount(%s, %d) not implemented yet.", datacenter, count)
	return nil
}

// GetJobCount will return the count of the nomad workers
func (c *Connector) GetJobCount(datacenter string) (uint, error) {
	c.log.Warn().Msgf("nomadworker.Connector.GetJobCount(%s) not implemented yet. Will return 0.", datacenter)
	return 0, nil
}

// IsJobDead will return if the nomad workers of the actual data-center are still available.
func (c *Connector) IsJobDead(datacenter string) (bool, error) {
	c.log.Warn().Msgf("nomadworker.Connector.IsJobDead(%s) not implemented yet. Will return false.", datacenter)
	return false, nil
}
