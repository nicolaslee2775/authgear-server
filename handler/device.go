package handler

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/oursky/ourd/oddb"
	"github.com/oursky/ourd/oderr"
	"github.com/oursky/ourd/router"
	"github.com/oursky/ourd/uuid"
)

type deviceRegisterPayload struct {
	ID          string
	Type        string
	DeviceToken string `mapstructure:"device_token"`
}

func (p *deviceRegisterPayload) Validate() error {
	if p.Type == "" {
		return errors.New("empty device type")
	} else if p.Type != "ios" && p.Type != "android" {
		return fmt.Errorf("unknown device type = %v", p.Type)
	}
	if p.DeviceToken == "" {
		return errors.New("empty device token")
	}

	return nil
}

// DeviceReigsterResult is the result put onto response.Result on
// successful call of DeviceRegisterHandler
type DeviceReigsterResult struct {
	ID string `json:"id"`
}

// DeviceRegisterHandler creates or updates a device and associates it to a user
//
// Example to create a new device:
//
//	curl -X POST -H "Content-Type: application/json" \
//	  -d @- http://localhost:3000/ <<EOF
//	{
//		"action": "device:register",
//		"access_token": "some-access-token",
//		"type": "ios",
//		"device_token": "some-device-token"
//	}
//	EOF
//
// Example to update an existing device:
//
//	curl -X POST -H "Content-Type: application/json" \
//	  -d @- http://localhost:3000/ <<EOF
//	{
//		"action": "device:register",
//		"access_token": "some-access-token",
//		"id": "existing-device-id",
//		"type": "ios",
//		"device_token": "new-device-token"
//	}
//	EOF
//
func DeviceRegisterHandler(rpayload *router.Payload, response *router.Response) {
	payload := deviceRegisterPayload{}
	if err := mapstructure.Decode(rpayload.Data, &payload); err != nil {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: "+err.Error())
		return
	}
	if err := payload.Validate(); err != nil {
		response.Err = oderr.New(oderr.RequestInvalidErr, "invalid request: "+err.Error())
		return
	}

	conn := rpayload.DBConn

	device := oddb.Device{}
	deviceID := payload.ID
	if deviceID == "" { // new device
		device.ID = uuid.New()
	} else { // update device
		if err := conn.GetDevice(deviceID, &device); err != nil {
			var errToReturn oderr.Error
			if err == oddb.ErrDeviceNotFound {
				errToReturn = oderr.New(oderr.RequestInvalidErr, "invalid request: updating device does not exist")
			} else {
				errToReturn = oderr.NewFmt(oderr.PersistentStorageErr, "persistent error: failed to fetch device: %v", deviceID)
			}
			response.Err = errToReturn
			return
		}
	}

	userinfoID := rpayload.UserInfo.ID

	device.Type = payload.Type
	device.Token = payload.DeviceToken
	device.UserInfoID = userinfoID

	if err := conn.SaveDevice(&device); err != nil {
		response.Err = oderr.New(oderr.PersistentStorageErr, "persistent error: failed to save device")
	} else {
		response.Result = DeviceReigsterResult{device.ID}
	}
}