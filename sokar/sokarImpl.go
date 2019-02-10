package sokar

func (sk *Sokar) Run() {

	// main loop
	go func() {
		sk.logger.Info().Msg("Sokar main loop started")

	loop:
		for {
			select {
			case <-sk.stopChan:
				// send the stop message a second time to complete waiting join calls
				sk.stopChan <- struct{}{}
				break loop
			}
		}
		sk.logger.Info().Msg("Sokar main loop left")
	}()

}

func (sk *Sokar) Stop() {
	sk.logger.Info().Msg("Teardown requested")
	// send the stop message
	sk.stopChan <- struct{}{}
}

func (sk *Sokar) Join() {
	<-sk.stopChan
}
