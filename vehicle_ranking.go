package ride

import (
	"time"
	"sort"
	// "github.com/kr/pretty"
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


func GetVehiclesRanking(vs []vehicle, reqPickUpPin, reqDropPin pin) Ranking {
	out := make(chan DeviationResult,len(vs))

	now := time.Now()

	for _, v := range vs {
		go calculateDeviation(v,reqPickUpPin,reqDropPin,now, out)
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