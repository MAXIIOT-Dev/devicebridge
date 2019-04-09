package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxiiot/vbaseBridge/storage"
)

var (
	VbID         = "b360662ed36b4633b3fcac95b4dd6477"
	CategoryCode = "residentialArea"
	CategoryName = "居民小区"
)

// VbaseDevice
type VbaseDevice struct {
	VbID         string `json:"vbid"`
	SrcID        string `json:"srcid"`
	Name         string `json:"name"`
	CategoryCode string `json:"categoryCode"`
	CategoryName string `json:"categoryName"`
	Addr         string `json:"addr"`
	Point        *Point `json:"point"`
}

// Point vbase point
type Point struct {
	ID       string   `json:"id"`
	Geometry Geometry `json:"geometry"`
}

// Geometry vbase geometry
type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// VbaseList vbase list
// @summary vbase地图设备列表
// @description vbase地图设备列表
// @tags vbaseapi
// @accept json
// @produce json
// @param status query string true "status online/offline"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /vbase/list [get]
func VbaseList(c *gin.Context) {
	status := c.DefaultQuery("status", "offline")
	if status != storage.VDSOnline && status != storage.VDSOffline {
		Response(c, http.StatusBadRequest, 0, 1, "status allow (online/offline)", nil)
		return
	}
	devs, err := storage.GetVbaseDevices(status)
	if err != nil {
		Response(c, http.StatusInternalServerError, 0, 1, err.Error(), nil)
		return
	}

	vdevs := make([]VbaseDevice, 0, len(devs))
	for _, dev := range devs {
		vdev := VbaseDevice{
			VbID:         VbID,
			SrcID:        dev.DeviceEUI.String(),
			Name:         dev.Name,
			CategoryCode: CategoryCode,
			CategoryName: CategoryName,
		}
		if dev.Location != nil {
			point := Point{
				Geometry: Geometry{
					Type:        "Point",
					Coordinates: []float64{dev.Location.Longitude, dev.Location.Latitude}, //[]float64{113.917974, 22.582150},
				},
			}
			vdev.Point = &point
		}
		vdevs = append(vdevs, vdev)
	}

	Response(c, http.StatusOK, 0, 0, "success", vdevs)
}

// VbaseCount returns device count
// @summary vbase地图设备数量
// @description vbase地图设备数量
// @tags vbaseapi
// @accept json
// @produce json
// @param status query string true "status online/offline"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /vbase/count [get]
func VbaseCount(c *gin.Context) {
	status := c.DefaultQuery("status", "offline")
	if status != storage.VDSOnline && status != storage.VDSOffline {
		Response(c, http.StatusBadRequest, 0, 1, "status allow (online/offline)", nil)
		return
	}

	count, err := storage.GetVbaseDevicesCount(status)
	if err != nil {
		Response(c, http.StatusInternalServerError, 0, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, 0, "success", count)
}

// VbaseTrack vbase地图设备跟踪
// @summary vbase地图设备跟踪
// @description vbase地图设备跟踪
// @tags vbaseapi
// @accept json
// @produce json
// @param id query string true "device id(eui)"
// @param st query int true "start timestamp"
// @param et query int true "end timestamp"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /vbase/track [get]
func VbaseTrack(c *gin.Context) {
	id := c.Query("id")
	st := c.Query("st")
	et := c.Query("et")

	if id == "" || st == "" || et == "" {
		Response(c, http.StatusBadRequest, 0, 1, "param erorr,should provide (id/st/et)", nil)
		return
	}

	ist, err := strconv.ParseInt(st, 10, 0)
	if err != nil || ist <= 0 {
		Response(c, http.StatusBadRequest, 0, 1, "st should integer", nil)
		return
	}
	iet, err := strconv.ParseInt(et, 10, 0)
	if err != nil || iet <= 0 {
		Response(c, http.StatusBadRequest, 0, 1, "et should interger", nil)
		return
	}

	// vbase st/et参数单位为毫秒
	start := time.Unix(ist/1000, 0)
	end := time.Unix(iet/1000, 0)

	tracks, err := storage.GetDeviceTrack(id, start, end)
	if err != nil {
		Response(c, http.StatusInternalServerError, 0, 1, err.Error(), nil)
		return
	}

	result := make([]string, 0, len(tracks))
	for _, track := range tracks {
		// vbase 轨迹时间单位为毫秒
		res := fmt.Sprintf("%d,%.6f,%.6f,%d", track.CreatedAt.Unix()*1000, track.Location.Latitude, track.Location.Longitude, track.Altitude)
		result = append(result, res)
	}
	results := make([][]string, 1)
	results[0] = result

	c.JSON(http.StatusOK, results)
}

// VbaseDetail vbase地图设备详情
// @summary vbase地图设备详情
// @description vbase地图设备详情
// @tags vbaseapi
// @accept json
// @produce json
// @param id query string true "device id(eui)"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @router /vbase/detail [get]
func VbaseDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		Response(c, http.StatusBadRequest, 0, 1, "param erorr,should provide (id/st/et)", nil)
		return
	}

	state, err := storage.GetDeviceState(id)
	if err != nil {
		if err == sql.ErrNoRows {
			Response(c, http.StatusOK, 0, 0, "", gin.H{
				"id":     id,
				"prop":   map[string]string{"设备ID": id},
				"sensor": map[string]string{"设备详情": "无"},
			})
			return
		}
		Response(c, http.StatusInternalServerError, 0, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, 0, "success", state.Detail)
}
