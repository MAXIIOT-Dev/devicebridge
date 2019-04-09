package lorahandler

import "math"

const (
	PI = 3.14159265358979324
	a  = 6378245.0              //  a: 卫星椭球坐标投影到平面地图坐标系的投影因子。
	ee = 0.00669342162296594323 //  ee: 椭球的偏心率。
)

// gps84ToGcj02 GPS84坐标系转GCJ02
// 参考资料 https://www.cnblogs.com/94cool/p/4266907.html
func gps84ToGcj02(wgsLat, wgsLon float64) (float64, float64) {
	if outOfChina(wgsLat, wgsLon) {
		return wgsLat, wgsLon
	}

	var dLat = transformLat(wgsLon-105.0, wgsLat-35.0)
	var dLon = transformLon(wgsLon-105.0, wgsLat-35.0)
	var radLat = wgsLat / 180.0 * PI
	var magic = math.Sin(radLat)
	magic = 1 - ee*magic*magic
	var sqrtMagic = math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * PI)
	dLon = (dLon * 180.0) / (a / sqrtMagic * math.Cos(radLat) * PI)
	return wgsLat + dLat, wgsLon + dLon
}

func outOfChina(lat float64, lon float64) bool {
	if lon < 72.004 || lon > 137.8347 {
		return true
	}
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	return false
}

func transformLat(x, y float64) float64 {
	var ret = -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*PI) + 20.0*math.Sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(y*PI) + 40.0*math.Sin(y/3.0*PI)) * 2.0 / 3.0
	ret += (160.0*math.Sin(y/12.0*PI) + 320*math.Sin(y*PI/30.0)) * 2.0 / 3.0
	return ret
}

func transformLon(x, y float64) float64 {
	var ret = 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20.0*math.Sin(6.0*x*PI) + 20.0*math.Sin(2.0*x*PI)) * 2.0 / 3.0
	ret += (20.0*math.Sin(x*PI) + 40.0*math.Sin(x/3.0*PI)) * 2.0 / 3.0
	ret += (150.0*math.Sin(x/12.0*PI) + 300.0*math.Sin(x/30.0*PI)) * 2.0 / 3.0
	return ret
}
