// data script somethings
package utils

import (
	luajson "github.com/layeh/gopher-json"
	"github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
)

func ScriptDeal(code string, msg []byte, topic string) (string, error) {
	/*
	   function encodeInp(msg, topic)
	       -- This function is an encoding function that converts the input message into a format recognizable by the platform or the device
	       -- Parameters:
	       -- msg: The input message (either a subscribed or reported message), type string
	       -- topic: The topic of the message, type string
	       -- Return value: Returns a string which is the encoded message

	       -- Please implement the encoding logic based on actual requirements

	       -- To convert between string and jsonObj, the json library needs to be imported:
	       -- local json = require("json")

	       -- Example: Converting a string to jsonObj:
	       -- local jsonTable = json.decode(msgString)

	       -- Example: Converting a jsonObj back to string:
	       -- local json_str = json.encode(jsonTable)

	       -- Below is the sample code:

	       -- Only supports the following json package import

	       -- After processing, convert the object back to string format
	       local json = require("json")
	       local jsonTable = json.decode(jsonString)

	       -- If the service_id in the first service is "CO2", update the property "current"
	       if jsonTable.services[1].service_id == "CO2" then
	           jsonTable.services[1].properties.current = 200
	       end

	       -- Convert the updated jsonTable back to a JSON string
	       local newJsonString = json.encode(jsonTable)

	       -- Return the newly encoded JSON string
	       return newJsonString
	   end
	*/

	L := lua.NewState()
	defer L.Close()

	L.PreloadModule("json", luajson.Loader)

	err := L.DoString(code)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	encodeInp := L.GetGlobal("encodeInp")
	err = L.CallByParam(lua.P{
		Fn:      encodeInp,
		NRet:    1,
		Protect: true,
	}, lua.LString(msg), lua.LString(topic))

	if err != nil {
		logrus.Error("Error executing Lua script:", err)
		return "", err
	}

	result := L.Get(-1)
	if result.Type() != lua.LTString {
		logrus.Error("Lua script must return a string")
		return "", err
	}
	return result.String(), nil
}
