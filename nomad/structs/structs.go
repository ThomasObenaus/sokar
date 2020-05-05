package structs

const (
	// DeploymentStatuses are the various states a deployment can be be in
	DeploymentStatusRunning    = "running"
	DeploymentStatusPaused     = "paused"
	DeploymentStatusFailed     = "failed"
	DeploymentStatusSuccessful = "successful"
	DeploymentStatusCancelled  = "cancelled"

	// DeploymentStatusDescriptions are the various descriptions of the states a
	// deployment can be in.
	DeploymentStatusDescriptionRunning               = "Deployment is running"
	DeploymentStatusDescriptionRunningNeedsPromotion = "Deployment is running but requires promotion"
	DeploymentStatusDescriptionPaused                = "Deployment is paused"
	DeploymentStatusDescriptionSuccessful            = "Deployment completed successfully"
	DeploymentStatusDescriptionStoppedJob            = "Cancelled because job is stopped"
	DeploymentStatusDescriptionNewerJob              = "Cancelled due to newer version of job"
	DeploymentStatusDescriptionFailedAllocations     = "Failed due to unhealthy allocations"
	DeploymentStatusDescriptionProgressDeadline      = "Failed due to progress deadline"
	DeploymentStatusDescriptionFailedByUser          = "Deployment marked as failed"
)

const (
	JobStatusPending = "pending" // Pending means the job is waiting on scheduling
	JobStatusRunning = "running" // Running means the job has non-terminal allocations
	JobStatusDead    = "dead"    // Dead means all evaluation's and allocations are terminal
)
