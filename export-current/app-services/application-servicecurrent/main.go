//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"fmt"
//        "log"
	"os"
//        "strconv"
	"encoding/json"
        "strings"
	"github.com/edgexfoundry/app-functions-sdk-go/appcontext"
	"github.com/edgexfoundry/app-functions-sdk-go/appsdk"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/transforms"
	"github.com/edgexfoundry/app-functions-sdk-go/pkg/util"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

const (
	serviceKey = "samplePrintEdgeXFoundryDataToConsole"
)
type Sensor_Data struct {
	Data_Id string `json:"Id"`
	Data_Layers int64 `json:"Total Layers"`
	Data_Name string `json:"Job Name"`
	Data_Times_Begin int64 `json:"Times Begin"`
	Data_Times_Remaining int64 `json:"Times Remaining"`
	Data_LayerIndex int64 `json:"Current Layer"`
}

type JSONString struct {
	ID string `json:"id"`
	Device string `json:"device"`
	Origin int64 `json:"origin"`
	Readings []Readings `json:"readings"`
}

type Readings struct {
        ID string `json:"id"`
	Created int64 `json:"created"`
	Origin string `json:"origin"`
	Device string `json:"device"`
	Name string `json:"name"`
	Value string `json:"value"`
        ValueType string `json:"valueType"`
}
type CurrentValue struct {
        Data Data `json:"data"`
	Id string `json:"id"`
	LayerIndex int64 `json:"LayerIndex"`
	Name	string `json:"Name"`
	TimeCostInfo TimeCostInfo `json:"TimeCostInfo"`
}

type Data struct {
        Files []Files `json:"files"`
        Id string `json:"id"`
	Layers int64 `json:"layers"`
	Name string `json:"name"`
	Operator string `json:"operator"`
	Size []int64 `json:"size"`
	Tickness float64 `json:"tickness"`
	Times Times `json:"times"`

}

type TimeCostInfo struct {
	Other int64 `json:"Other"`
	Pause int64 `json:"Pause"`
	PrepareData int64 `json:"PrepareData"`
	Scan int64 `json:"Scan"`
	SpreadPower int64 `json:"SpreadPower"`
}

type Files struct {
	File string `json:"file"`
	Matrix []int64 `json:"matrix"`
}

type Times struct {
	Begin int64  `json:"begin"`
        Remaining int64 `json:"remaining"`
}

var counter int 

func main() {

	// 1) First thing to do is to create an instance of the EdgeX SDK and initialize it.
	edgexSdk := &appsdk.AppFunctionsSDK{ServiceKey: serviceKey}
	if err := edgexSdk.Initialize(); err != nil {
		message := fmt.Sprintf("SDK initialization failed: %v\n", err)
		if edgexSdk.LoggingClient != nil {
			edgexSdk.LoggingClient.Error(message)
		} else {
			fmt.Println(message)
		}
		os.Exit(-1)
	}

	// 2) shows how to access the application's specific configuration settings.
	deviceNames, err := edgexSdk.GetAppSettingStrings("DeviceNames")
	if err != nil {
		edgexSdk.LoggingClient.Error(err.Error())
		os.Exit(-1)
	}
	edgexSdk.LoggingClient.Info(fmt.Sprintf("Filtering for devices %v", deviceNames))


	// Since we are using MQTT, we'll also need to set up the addressable model to
	// configure it to send to our broker. If you don't have a broker setup you can pull one from docker i.e:
	// docker run -it -p 1883:1883 -p 9001:9001  eclipse-mosquitto
	addressable := models.Addressable{
		Address:   "172.17.48.170",
//                Address:   "192.168.1.249h",
		Port:      1883,
		Protocol:  "tcp",
		Publisher: "MyAppCurrent",
		User:      "",
		Password:  "",
//		Topic:     "test/topic",
                Topic:     "topic/current",
	}

	// Using default settings, so not changing any fields in MqttConfig
	mqttConfig := transforms.MqttConfig{}

	// Make sure you change KeyFile and CertFile here to point to actual key/cert files
	// or an error will be logged for failing to load key/cert files
	// If you don't use key/cert for MQTT authentication, just pass nil to NewMQTTSender() as following:
	//	mqttSender := transforms.NewMQTTSender(edgexSdk.LoggingClient, addressable, nil, mqttConfig)
	pair := transforms.KeyCertPair{
		KeyFile:  "PATH_TO_YOUR_KEY_FILE",
		CertFile: "PATH_TO_YOUR_CERT_FILE",
	}

	mqttSender := transforms.NewMQTTSender(edgexSdk.LoggingClient, addressable, &pair, mqttConfig, false)

	
	
	// 3) This is our pipeline configuration, the collection of functions to
	// execute every time an event is triggered.
	fmt.Println("Setting the pipeline functions\n\n")
	fmt.Println("Preparing to print EdgeX Foundry Reading Data to console\n\n\n")
	edgexSdk.SetFunctionsPipeline(
		transforms.NewFilter(deviceNames).FilterByDeviceName,
//		SetOutputData,
//                CustomJson,
 		transforms.NewConversion().TransformToJSON,
//		SetOutputData,
                funzioneditrasformazione,
//		printEdgeXReadingDataToConsole,
//		SetOutputData,
		mqttSender.MQTTSend,
	)

	// 4) Lastly, we'll go ahead and tell the SDK to "start" and begin listening for events
	// to trigger the pipeline.
	err = edgexSdk.MakeItRun()
	if err != nil {
		edgexSdk.LoggingClient.Error("MakeItRun returned error: ", err.Error())
		os.Exit(-1)
	}

	// Do any required cleanup here

	os.Exit(0)
}

      
func funzioneditrasformazione(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	// Save the event reading from EdgeX Foundry to
	// the variable myJsonString
	myJsonString := params[0].(string)
//	fmt.Println("\n\n", "contenuto topic event", myJsonString)

	// Convert myJsonString to bytes
	jsonData := []byte(myJsonString)

        	
	// Define a helping variable of
	// type JSONString - defined structure
	var JsonData JSONString
        	
	// Unmarshal jsonData
	err := json.Unmarshal(jsonData, &JsonData)
     
	if err != nil {
	fmt.Println(err)
	}
        if JsonData.Device == "DISPOSITIVO_CURRENT_PRINTER"{
	        // Value e valueacc sono stringhe
	        fmt.Println("\n\n", "STRINGA JSONacc ISOLATO", JsonData.Readings[0].Value)
//       	 fmt.Println("\n\n", JsonData.Device)

        	valuecurrent1 := strings.NewReplacer("'", "\"")
        	valuecurrent := valuecurrent1.Replace(JsonData.Readings[0].Value)
//        	fmt.Println("\n\n", "STRINGA JSONPIR ISOLATO CON DOPPI APICI", valueacc)
        	// save the parsed data in (valueacc) in the helping variable
        	// ValueAcc (and convert it to byte type) poichè serve per l'unmarshal 

        
        	ValueCurrent := []byte(string(valuecurrent))

        	fmt.Println("\n\n", "prova", ValueCurrent)
      
  
		// Define a helping variable ValueJsonData

		var ValueJsonData CurrentValue

		// Unmarshal the valueacc string in ValueJsonData
		Err := json.Unmarshal(ValueCurrent, &ValueJsonData)

		if Err != nil {
			fmt.Println(Err)
		}

       
		fmt.Println("\n\n", "unmarshal della struct JSONCURRENT", ValueJsonData)
//        
        	

        	var data Sensor_Data



        	data.Data_Id = ValueJsonData.Data.Id
		data.Data_Layers = ValueJsonData.Data.Layers
        	data.Data_Name = ValueJsonData.Data.Name
		data.Data_Times_Begin = ValueJsonData.Data.Times.Begin
		data.Data_Times_Remaining = ValueJsonData.Data.Times.Remaining
		data.Data_LayerIndex = ValueJsonData.LayerIndex

        	fmt.Println("\n\n", "stampa struct data di export", data)

		SendTheFollowingDataToHighLeveApplications, err := json.Marshal(data)
        	if err != nil {
        		fmt.Println(err)
        	}
        	NewJson := string( SendTheFollowingDataToHighLeveApplications)
        	fmt.Println("\n\n", "stampa stringa di export", NewJson)
       
        

        	return true, NewJson
//          return JsonData.Device
}

 
        return false, params[0] 
//	  return true, nil      
}

// SetOutputData sets the output data to that passed in from the previous function.
// It will return an error and stop the pipeline if the input data is not of type []byte, string or json.Mashaler
func SetOutputData(edgexcontext *appcontext.Context, params ...interface{}) (bool, interface{}) {
	fmt.Println("\nSetting output data :)")
	edgexcontext.LoggingClient.Debug("Setting output data")

	if len(params) < 1 {
		// We didn't receive a result
		return false, nil
	}

	data, err := util.CoerceType(params[0])
	if err != nil {
		return false, err
	}
	// By setting this the data will be posted back to to configured trigger response, i.e. message bus
	edgexcontext.OutputData = data
	fmt.Println("\n",string(data), "\n")

	return false, params[0]
}



