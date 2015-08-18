package handler

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/jeffbmartinez/loadbalancer/host"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Balancer struct {
	Hosts []host.Host

	// bagOfHostIndexes is collection of indexes into Hosts that allows us
	// to select a random host, taking their weights into account.
	bagOfHostIndexes []int
}

func NewBalancer(hosts []host.Host) Balancer {
	balancer := Balancer{
		Hosts:            hosts,
		bagOfHostIndexes: makeBagOfHostIndexes(hosts),
	}

	return balancer
}

func (b Balancer) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	b.chooseRandomHost().ServeHTTP(response, request)
}

func (b Balancer) chooseRandomHost() host.Host {
	bagIndex := rand.Int() % len(b.bagOfHostIndexes)
	hostIndex := b.bagOfHostIndexes[bagIndex]

	return b.Hosts[hostIndex]
}

func makeBagOfHostIndexes(hosts []host.Host) []int {
	hostIndexes := make([]int, 0)

	for hostIndex, host := range hosts {
		for i := 0; i < host.Weight; i++ {
			hostIndexes = append(hostIndexes, hostIndex)
		}
	}

	return hostIndexes
}
