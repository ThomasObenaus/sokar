package nomadWorker

import "fmt"

func (c *Connector) downscale(datacenter string, desiredCount uint) error {

	// 1. Select a candidate for downscaling -> returns [needs node id]
	// 2. Make the node ineligible [needs node id]
	// 3. Drain the node [needs node id]
	// 4. Terminate the node using the AWS ASG [needs instance id]

	return fmt.Errorf("Downscaling is not yet implemented")
}

func selectCandidateForDownscaling(datacenter string) (nodeID string, err error) {
	return "", nil
}
