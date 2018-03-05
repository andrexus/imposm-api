package service

type NearbyPoint struct {
	OSMId     int     `json:"osmId"`
	Name      string  `json:"name"`
	PointType string  `json:"pointType,omitempty"`
	Distance  float32 `json:"distance"`
}

type Location struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
}

type NearbyPointsSearchRequest struct {
	Location       Location `json:"location"`
	Radius         int      `json:"radius"`
	TransportTypes []string `json:"transportTypes"`
	Distinct       bool     `json:"distinct"`
}

type ImposmService interface {
	FindNearbyTransportPoints(r *NearbyPointsSearchRequest) ([]NearbyPoint, error)
	FindNearbyStreets(r *NearbyPointsSearchRequest) ([]NearbyPoint, error)
}

type ImposmServiceImpl struct {
	Repository ImposmRepository
}

func NewImposmService(repository ImposmRepository) *ImposmServiceImpl {
	return &ImposmServiceImpl{
		Repository: repository,
	}
}

func (c *ImposmServiceImpl) FindNearbyTransportPoints(r *NearbyPointsSearchRequest) ([]NearbyPoint, error) {
	return c.Repository.FindNearbyTransportPoints(r.Location, uint(r.Radius), r.Distinct, r.TransportTypes)
}

func (c *ImposmServiceImpl) FindNearbyStreets(r *NearbyPointsSearchRequest) ([]NearbyPoint, error) {
	return c.Repository.FindNearbyStreets(r.Location, uint(r.Radius), r.Distinct)
}
