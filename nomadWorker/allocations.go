package nomadWorker

import (
	"fmt"

	nomadApi "github.com/hashicorp/nomad/api"
)

type allocInfo struct {
	numAllocations uint
	resources
}

type resources struct {
	cpu      int
	diskMB   int
	memoryMB int
}

func getNumAllocationsInStatus(nodesIF Nodes, nodeID string, status string) (*allocInfo, error) {
	allocs, _, err := nodesIF.Allocations(nodeID, nil)
	if err != nil {
		return nil, err
	}

	allocInformation := &allocInfo{}
	for _, alloc := range allocs {
		if alloc == nil || alloc.Job == nil || alloc.Job.Status == nil {
			continue
		}
		if alloc.Resources == nil || alloc.Resources.CPU == nil || alloc.Resources.MemoryMB == nil || alloc.Resources.DiskMB == nil {
			continue
		}

		if *alloc.Job.Status != nomadApi.AllocClientStatusRunning {
			continue
		}

		allocInformation.numAllocations++
		allocInformation.cpu += *alloc.Resources.CPU
		allocInformation.diskMB += *alloc.Resources.DiskMB
		allocInformation.memoryMB += *alloc.Resources.MemoryMB
	}
	return allocInformation, nil
}

func (r *resources) String() string {
	return fmt.Sprintf("cpu=%d,disk=%d MB,memory=%d MB", r.cpu, r.diskMB, r.memoryMB)
}
