package service

import (
	"database/sql"
	"github.com/lib/pq"
)

var defaultTransportTypes = []string{"bus_stop", "halt", "station", "tram_stop"}

//noinspection SqlResolve
const queryTransportPointsNearby = `
SELECT
	  osm_id,
	  name,
	  type,
	  min(ST_Distance(geometry, ref_point)) distance
	FROM (
		   SELECT
			 osm_id,
			 name,
			 type,
			 geometry,
			 ST_Transform(ST_SetSRID(ST_MakePoint($1, $2), 4326), 3857) ref_point
		   FROM osm_transport_points
		 ) AS t
	WHERE
	  type = ANY($4)
	  AND ST_DWithin(geometry, ref_point, $3)
	GROUP BY osm_id, name, type
	ORDER BY
	  distance;
`

//noinspection SqlResolve
const queryStreetsNearby = `
	SELECT
	  osm_id,
	  name,
	  type,
	  min(ST_Distance(geometry, ref_point)) distance
	FROM (
		   SELECT
			 osm_id,
			 name,
			 geometry,
			 type,
			 ST_Transform(ST_SetSRID(ST_MakePoint($1, $2), 4326), 3857) ref_point
		   FROM osm_roads
		 ) AS t
	WHERE
	  name != ''
	  AND ST_DWithin(geometry, ref_point, $3)
	GROUP BY osm_id, name, type
	ORDER BY
	  distance;
`

type ImposmRepository interface {
	FindNearbyTransportPoints(location Location, radius uint, transportTypes []string) ([]NearbyPoint, error)
	FindNearbyStreets(location Location, radius uint) ([]NearbyPoint, error)
}

type PostGISImposmRepository struct {
	db *sql.DB
}

func NewImposmRepository(db *sql.DB) *PostGISImposmRepository {
	return &PostGISImposmRepository{db}
}

func (r *PostGISImposmRepository) FindNearbyTransportPoints(location Location, radius uint, transportTypes []string) ([]NearbyPoint, error) {
	//noinspection GoPreferNilSlice
	items := []NearbyPoint{}

	transportPointTypes := transportTypes
	if transportPointTypes == nil {
		transportPointTypes = defaultTransportTypes
	}
	rows, err := r.db.Query(queryTransportPointsNearby, location.Lon, location.Lat, radius, pq.Array(transportPointTypes))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var osmId int
		var name string
		var pType string
		var distance float32
		err = rows.Scan(&osmId, &name, &pType, &distance)
		item := &NearbyPoint{
			OSMId:     osmId,
			Name:      name,
			PointType: pType,
			Distance:  distance,
		}
		items = append(items, *item)
	}

	return items, nil
}

func (r *PostGISImposmRepository) FindNearbyStreets(location Location, radius uint) ([]NearbyPoint, error) {
	//noinspection GoPreferNilSlice
	items := []NearbyPoint{}

	rows, err := r.db.Query(queryStreetsNearby, location.Lon, location.Lat, radius)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var osmId int
		var name string
		var rType string
		var distance float32
		err = rows.Scan(&osmId, &name, &rType, &distance)
		item := &NearbyPoint{
			OSMId:    osmId,
			Name:     name,
			Distance: distance,
		}
		items = append(items, *item)
	}

	return items, nil
}
