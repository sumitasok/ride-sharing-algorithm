package ride

import (
	"time"
	"sort"
	"errors"
)

const (
	RADIUS = 30
)


type ResponseVehicle struct {
	Accept bool
	Rank int
}
type Ranking []DeviationResult


func (r Ranking) Len() int {
	return len(r)
}

func (r Ranking) Less(i, j int) bool {
	return r[i].VehicleDeviation < r[j].VehicleDeviation;
}

func (r Ranking) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}


func GetVehiclesRanking(vs []vehicle, reqID string, reqPickUpPin, reqDropPin pin) Ranking {
	out := make(chan DeviationResult,len(vs))

	now := time.Now()

	for _, v := range vs {
		go calculateDeviation(v,reqID,reqPickUpPin,reqDropPin,now, out)
	}

	ranking := Ranking{}
	for i := 0; i < len(vs);i++ {
		ranking = append(ranking,<- out)
	}
	close(out)

	sort.Sort(ranking)
	// pretty.Println("RANKING:::", ranking)

	return ranking
}

func AssignVehicles(req requestor, vs []vehicle) (DeviationResult,error) {
	reqPickUpPin :=  *NewPinFromRequestor(req, pickup) 	// New pin for upcoming rider's pickup
	reqDropPin :=  *NewPinFromRequestor(req, drop)		// New pin for upcoming rider's drop


	ranks := GetVehiclesRanking(vs, req.Identifier, reqPickUpPin, reqDropPin)

	nbrOfReqVeh := len(ranks)

	respChan := make(chan ResponseVehicle, nbrOfReqVeh)
	if len(ranks) < 0 {
		return DeviationResult{}, errors.New("No Vehicle found")
	}
	nbrOfReqVeh = 0
	for i := range ranks {
		go sendRequestToVehicle(ranks[i].V, i,respChan)
		nbrOfReqVeh++
	}

	var resp ResponseVehicle

	for i := 0; i < nbrOfReqVeh; i++ {
		resp = <- respChan
		if resp.Accept {
			break
		}
	}

	if !resp.Accept {
		return DeviationResult{}, errors.New("No Vehicle found")
	}

	selRank := ranks[resp.Rank]
	//adding incoming requestor to vehicles requestor list
	if selRank.V.Riders == nil {
		selRank.V.Riders = make(map[string]*requestor)
	}
	req.State = rideRequested
	req.DirectDropTime = selRank.DirectDropTime
	req.ExpectedPickUpTime = selRank.PickUpTime
	req.ExpectedPickUpTime = selRank.DropTime
	selRank.V.Riders[req.Identifier] = &req
	selRank.V.CurrentRoute = routes{Pins:selRank.Route}
	return selRank, nil
}

func sendRequestToVehicle(v vehicle, rank int,out chan ResponseVehicle) {
	time.Sleep(time.Second*20*time.Duration(rank))
	//TBD Calling vehicle and getting vehicle response
	Accept := true
	out <-  ResponseVehicle{
		Accept,
		rank,
	}
}