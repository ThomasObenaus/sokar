package nomadConnector

func (nc *connectorImpl) ScaleBy(amount int) error {
	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...")

	nc.log.Info().Str("job", nc.jobName).Int("amount", amount).Msg("Scaling ...done")
	return nil
}
