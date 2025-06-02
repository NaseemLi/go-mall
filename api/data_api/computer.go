package dataapi

import (
	"fast_gin/utils/computer"
	"fast_gin/utils/res"

	"github.com/gin-gonic/gin"
)

type ComputerResponse struct {
	CpuPercent  float64 `json:"cpuPercent"`
	MemPercent  float64 `json:"memPercent"`
	DiskPercent float64 `json:"diskPercent"`
}

func (DataApi) ComputerView(c *gin.Context) {
	data := ComputerResponse{
		CpuPercent:  computer.GetCpuPercent(),
		MemPercent:  computer.GetMemPercent(),
		DiskPercent: computer.GetDiskPercent(),
	}

	res.OkWithData(data, c)
}
