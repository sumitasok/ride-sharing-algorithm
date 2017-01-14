package ride

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVehicleOccupancy(t *testing.T) {
	assert := assert.New(t)

	v := NewVehicle(4, location{Lat: 12.223234, Long: 77.23123})
	v.addRequestor(*NewRequestor("rider-1", 1, time.Now(), location{Lat: 12.223234, Long: 77.23123}, location{Lat: 12.223234, Long: 77.23123}))

	assert.Equal(int64(1), v.occupancyStatus())

	assert.True(true)
}

func TestVehicleRemainingOccupancy(t *testing.T) {
	assert := assert.New(t)

	v := NewVehicle(4, location{Lat: 12.223234, Long: 77.23123})
	v.addRequestor(*NewRequestor("rider-1", 1, time.Now(), location{Lat: 12.223234, Long: 77.23123}, location{Lat: 12.223234, Long: 77.23123}))

	assert.Equal(int64(3), v.remainingOccupancy())

	assert.True(true)
}