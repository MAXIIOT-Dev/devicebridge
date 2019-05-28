package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/maxiiot/devicebridge/backend/server"
	"github.com/maxiiot/devicebridge/storage"

	"github.com/gin-gonic/gin"
)

// Device for request device.
type Device struct {
	DeviceEUI    string `json:"device_eui" binding:"required"`
	ProtocolType string `json:"protocol_type" example:"optional(humiture/smoke/digital)"`
}

func (dev *Device) validate() error {
	dev.ProtocolType = strings.ToLower(dev.ProtocolType)
	if dev.ProtocolType != storage.ProtocolHumiture && dev.ProtocolType != storage.ProtocolSmoke &&
		dev.ProtocolType != storage.ProtocolDefault {
		return errors.New("Unsupported device protocol.optional(humiture/smoke/digital)")
	}
	return nil
}

func (dev Device) toStorageDevice() (storage.Device, error) {
	var sDev storage.Device
	var devEUI storage.EUI64
	err := devEUI.UnmarshalText([]byte(dev.DeviceEUI))
	if err != nil {
		return sDev, err
	}
	sDev.DeviceEUI = devEUI
	sDev.ProtocolType = dev.ProtocolType
	return sDev, nil
}

// @summary 新增设备
// @description 新增设备
// @tags device
// @accept json
// @produce json
// @param device body controllers.Device true "create device info"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /device [post]
func CreateDevice(c *gin.Context) {
	var dev Device
	err := c.ShouldBind(&dev)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	if err = dev.validate(); err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	sDev, err := dev.toStorageDevice()
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	err = storage.CreateDevice(sDev)
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	server.Serv.OnDeviceChange(map[string]bool{dev.DeviceEUI: true})

	Response(c, http.StatusOK, 0, "success", nil)
}

// @summary 设备列表
// @description 设备列表
// @tags device
// @accept json
// @produce json
// @param page query int true "page"
// @param perpage query int true "perpage"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /device [get]
func ListDevice(c *gin.Context) {
	_page := c.DefaultQuery("page", "1")
	_perpage := c.DefaultQuery("perpage", "10")

	page, err := strconv.Atoi(_page)
	if err != nil || page < 1 {
		Response(c, http.StatusBadRequest, 1, "page must be >=1", nil)
		return
	}

	perpage, err := strconv.Atoi(_perpage)
	if err != nil || perpage <= 0 {
		Response(c, http.StatusBadRequest, 1, "perpage must be >0", nil)
		return
	}

	limit := perpage
	offset := perpage * (page - 1)

	devs, err := storage.GetDevices(limit, offset)
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	count, err := storage.GetDevicesCount()
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, "success", gin.H{
		"total":   count,
		"devices": devs,
	})
}

// @summary 设备明细
// @description  设备明细
// @tags device
// @accept json
// @produce json
// @param dev_eui path string true "device eui"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /device/{dev_eui} [get]
func GetDevice(c *gin.Context) {
	eui := c.Param("dev_eui")
	if eui == "" {
		Response(c, http.StatusBadRequest, 1, "please provide dev_eu param", nil)
		return
	}

	dev, err := storage.GetDeviceByEUI(eui)
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, "success", dev)
}

// @summary 修改设备
// @description 修改设备
// @tags device
// @accept json
// @produce json
// @param device body controllers.Device true "update device info"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /device [put]
func UpdateDevice(c *gin.Context) {
	var dev Device
	err := c.ShouldBind(&dev)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	if err := dev.validate(); err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	sDev, err := dev.toStorageDevice()
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	err = storage.UpdateDevice(sDev)
	if err != nil {
		Response(c, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}

	Response(c, http.StatusOK, 0, "success", nil)
}

// @summary 删除设备
// @description 删除设备
// @tags device
// @accept json
// @produce json
// @param dev_eui path string true "device eui"
// @success 200 {object} controllers.ResponseData
// @failure 500 {object} controllers.ResponseData
// @security ApiKeyAuth
// @router /device/{dev_eui} [delete]
func DeleteDevice(c *gin.Context) {
	eui := c.Param("dev_eui")
	if eui == "" {
		Response(c, http.StatusBadRequest, 1, "please provide dev_eu param", nil)
		return
	}

	err := storage.DeleteDevice(eui)
	if err != nil {
		Response(c, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	server.Serv.OnDeviceChange(map[string]bool{eui: false})

	Response(c, http.StatusOK, 0, "success", nil)
}
