# CapacityPlanner (CAP)

The CapacityPlanner is responsible to decide how many instances/ allocations of the scaling-object shall be deployed at a certain point in time. When the CAP is asked to determine the "planned scale" he regards all information available based on the selected CAP mode. Usually the CAP is asked to calculate this new count of the scaling-object in case sokar decided to scale up or down.

## Modes

There are 3 modes available in order to calculate the new count of the scaling-object:

- Constant Mode **(default mode)**
- Linear Mode
- Stepwise Mode

### Constant Mode

In this mode the CAP:

- Adds a constant offset to the current count of the scaling-object if an upscale is needed.
- Subtracts a constant offset to the current count of the scaling-object if an upscale is needed.

#### Inputs

- `current_scale`: Current count of the scaling-object.
- `offset`: Config-Parameter.

#### Calculation

- `planned_scale = current_scale + offset`

#### Example

Start sokar with CAP in constant mode using an offset of 2 (scales n+2).
`./sokar --cap.constant-mode.enable --cap.constant-mode.offset=2`

### Linear Mode

In this mode the CAP scales linearly based on the current value of the `scaleFactor`. The value for adjusting the count is calculated by multiplying the current count of the scaling-object with the `scaleFactor`. This value is then added to the current count to calculate the "planned scale". This means if the `scaleFactor` is 1.0 the count of the scaling-object is increased by 100% (doubled).

Even if the current count is 0 the CAP will use an adjustment value of 0 in order to avoid starvation at a scale of 0.

To control a bit the impact of the `scaleFactor` on the planning the parameter `scaleFactorWeight` can be used. Internally the `scaleFactorWeight` and `scaleFactor` are multiplied.

#### Inputs

- `current_scale`: Current count of the scaling-object.
- `scaleFactor`: Current steepness of the change of the incoming alert weights.
- `scaleFactorWeight`: Config-Parameter.

#### Calculation

- `planned_scale = current_scale * (1 + scaleFactor * scaleFactorWeight)`

#### Example

Start sokar with CAP in linear mode using a `scaleFactorWeight` of 0.7.
`./sokar --cap.constant-mode.enable=false --cap.linear-mode.enable --cap.linear-mode.scaleFactorWeight=0.7`

### Stepwise Mode

- Not implemented yet.
