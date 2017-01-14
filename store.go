package ride

import (
	"fmt"
	"gopkg.in/redis.v5"
	"encoding/gob"
	"bytes"
	"errors"
)

var(
	redisST = NewRedisStore("localhost:6379", "")
)

func Store() {

}

type distanceUnit string

const (
	KM distanceUnit = "km"
	M  distanceUnit = "m"
)

type store interface {
	FetchAllVehicles() []vehicle
	FetchAllByRadius(string, float32, distanceUnit) []vehicle
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string) *redisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>
	return &redisStore{client}
}

func (r redisStore) AddVehicle(key, name string, long, lat float64) (int64, error) {
	geoLocation := &redis.GeoLocation{
		Name:      name,
		Longitude: long,
		Latitude:  lat,
	}

	intCmd := r.client.GeoAdd(key, geoLocation)
	return intCmd.Result()
}

func (r redisStore) FetchAllByRadius(key string, long, lat, radius float64, unit distanceUnit) ([]redis.GeoLocation, error) {
	radiusQuery := &redis.GeoRadiusQuery{
		Radius:   radius,
		Unit:     string(unit),
		WithDist: true,
		WithCoord:true,
		Sort: "ASC",
	}
	geoLocations := r.client.GeoRadius(key, long, lat, radiusQuery)
	return geoLocations.Result()
}

func (r redisStore) FetchVehicleDetail(keys ...string) ([]vehicle, error){
	results, err := r.client.MGet(keys...).Result()
	if err != nil {
		return nil,err
	}
	vs := []vehicle{}
	for _, result := range results {
		v  := vehicle{}
		vBuff := bytes.NewBufferString(result.(string))
		dec := gob.NewDecoder(vBuff)
		err = dec.Decode(&v)
		if err != nil {
			return nil,err
		}
		vs = append(vs, v)
	}
	return vs,nil
}

func (r redisStore) InsertVehicles(vs... vehicle) (string,error)  {

	var vStr []interface{}
	for _, v := range vs {
		var vehicleBuff bytes.Buffer
		enc := gob.NewEncoder(&vehicleBuff)

		err := enc.Encode(v)
		if err != nil {
			fmt.Println("Err:::", err)
			return "", errors.New("Can't encode to gob")
		}
		vStr = append(vStr,v.ID, vehicleBuff.String())
	}

	return r.client.MSet(vStr...).Result()
}

func (r redisStore) PickupRider(vehicle_id, rider_id string) error {
	vehicles, err := redisST.FetchVehicleDetail(vehicle_id)
	if err != nil {
		return err
	}

	if len(vehicles) == 0 {
		return errors.New("Vehicle not found")
	}

	v := vehicles[0]
	err = v.Pickup(rider_id)
	if err != nil {
		return err
	}

	_, err = redisST.InsertVehicles(v)
	if err != nil {
		return err
	}

	return nil
}

func (r redisStore) RemoveVehicle(key, name string) (int64, error) {
	intCmd := r.client.ZRem(key, name)
	return intCmd.Result()
}