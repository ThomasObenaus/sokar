# CapacityPlanner (CAP)

The CapacityPlanner is responsible to decide how many instances/ allocations of the scaling-object shall be deployed at a certain point in time. When the CAP is asked to determine the "planned scale" he regards all information available based on the selected CAP mode. Usually the CAP is asked to calculate this new count of the scaling-object in case sokar decided to scale up or down.

## Features

- [Scheduled Scaling](ScheduledScaling.md)
- Configurable [Planning Modes](PlanningModes.md)
