package gps

import (
	"math"
)

const (
	XPI = 3.14159265358979324 * 3000.0 / 180.0
	PI  = 3.1415926535897932384626 // π
	A   = 6378245.0                // 长半轴
	EE  = 0.00669342162296594323   // 扁率
)

//判断是否是中国境内
func InChina(lng, lat float64) bool {
	return lng > 73.66 && lng < 135.05 && lat > 3.86 && lat < 53.55
}

//WGS84转GCJ02(火星坐标系)；
//一般情况下，google地图、高德地图、腾讯地图，使用这个坐标；
func Wgs84ToGcj02(lng, lat float64) (float64, float64) {
	if !InChina(lng, lat) {
		return lng, lat
	}
	dLat := transformLat(lng-105.0, lat-35.0)
	dLng := transformLng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * PI
	magic := math.Sin(radLat)
	magic = 1 - EE*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * PI)
	dLng = (dLng * 180.0) / (A / sqrtMagic * math.Cos(radLat) * PI)
	return lng + dLng, lat + dLat
}

//GCJ02(火星坐标系)转WPS84
func Gcj02ToWps84(lng, lat float64) (float64, float64) {
	if !InChina(lng, lat) {
		return lng, lat
	}
	dLat := transformLat(lng-105.0, lat-35.0)
	dLng := transformLng(lng-105.0, lat-35.0)
	radLat := lat / 180.0 * EE
	magic := math.Sin(radLat)
	magic = 1 - EE*magic*magic
	sqrtMagic := math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((A * (1 - EE)) / (magic * sqrtMagic) * PI)
	dLng = (dLng * 180.0) / (A / sqrtMagic * math.Cos(radLat) * PI)
	return lng*2 - (lng + dLng), lat*2 - (lat + dLat)
}

//火星坐标转百度坐标
func Gcj02ToBd09(lng, lat float64) (float64, float64) {
	z := math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*XPI)
	theta := math.Atan2(lat, lng) + 0.000003*math.Cos(lng*XPI)
	return z*math.Cos(theta) + 0.0065, z*math.Sin(theta) + 0.006
}

//百度坐标转火星
func Bd09ToGcj02(lng, lat float64) (float64, float64) {
	x := lng - 0.0065
	y := lat - 0.006
	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPI)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*XPI)
	return z * math.Cos(theta), z * math.Sin(theta)
}

//wgs84转百度
func Wgs84ToBd09(lng, lat float64) (float64, float64) {
	return Gcj02ToBd09(Wgs84ToGcj02(lng, lat))
}

//百度转wgs84
func Bd09ToWgs84(lng, lat float64) (float64, float64) {
	return Bd09ToGcj02(Gcj02ToWps84(lng, lat))
}

func transformLat(lng, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*PI) + 40.0*math.Sin(lat/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*PI) + 320*math.Sin(lat*PI/30.0)) * 2.0 / 3.0
	return ret
}

func transformLng(lng, lat float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*PI) + 20.0*math.Sin(2.0*lng*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*PI) + 40.0*math.Sin(lng/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*PI) + 300.0*math.Sin(lng/30.0*PI)) * 2.0 / 3.0
	return ret
}
