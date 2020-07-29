package service

import (
	...
)

type GeoLocation struct {
	Lat float64
	Long float64
	FriendlyName string
}

type GeoIPLocator struct {
	mongo                 *store.Mongo
	globalGeoIPCache      *store.MySuperCoolCacheImpl
	geoLocatorAPIHost     string
	geoLocatorAPIUsername string
	geoLocatorAPIPassword string
}

func NewGeoIPLocator(mongo *store.Mongo, globalCache *store.MySuperCoolCacheImpl, geoLocatorAPIHost, geoLocatorAPIUsername, geoLocatorAPIPassword string) *GeoIPLocator {
	return &GeoIPLocator{
		mongo:                 mongo
		globalGeoIPCache:      globalCache,
		geoLocatorAPIHost:     geoLocatorAPIHost,
		geoLocatorAPIUsername: geoLocatorAPIUsername,
		geoLocatorAPIPassword: geoLocatorAPIPassword,
	}
}

func (l *GeoIPLocator) GeoLocateIP(ip string) (GeoLocation, error) {

	location, found := l.globalGeoIPCache.Get(ip)
	if !found {
		return GeoLocation{}, errors.New("Could not geo locate IP address :(")
	}

	location, err := l.mongo.FindLocationForIP(ip)
	if err != nil {
		return GeoLocation{}, errors.New("Could not geo locate IP Address :(")
	}

	location, err = makeGetRequest(geoLocatorAPIHost + "/geolocate?ip="+ip+"&username="+l.geoLocatorAPIUsername)
	if err != nil {
		return GeoLocation{}, err
	}

	return location, nil
}
