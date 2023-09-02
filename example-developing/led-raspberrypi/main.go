package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github/babydeng/kubeedge-example/led-raspberrypi/configuration"
	lightdriver "github/babydeng/kubeedge-example/led-raspberrypi/light_driver"
	"os"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
)

const (
	modelName                 = "LED-LIGHT"
	powerStatus               = "power-status"
	pinNumberConfig           = "GPIO-PIN-NUMBER"
	DeviceETPrefix            = "$hw/events/device/"
	DeviceETStateUpdateSuffix = "/state/update"
	TwinETUpdateSuffix        = "/twin/update"
	TwinETCloudSyncSuffix     = "/twin/cloud_updated"
	TwinETGetResultSuffix     = "/twin/get/result"
	TwinETGetSuffix           = "/twin/get"
	MqttURL                   = "tcp://10.177.21.124:1883"
	DeviceName                = "led-light-instance-01"
)

var cli MQTT.Client
var Token_client Token
var wg sync.WaitGroup
var deviceTwinResult DeviceTwinUpdate
var deviceID string
var pinNumber float64

// Token interface to validate the MQTT connection.
type Token interface {
	Wait() bool
	WaitTimeout(time.Duration) bool
	Error() error
}

// DeviceStateUpdate is the structure used in updating the device state
type DeviceStateUpdate struct {
	State string `json:"state,omitempty"`
}

// BaseMessage the base struct of event message
type BaseMessage struct {
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

// TwinValue the struct of twin value
type TwinValue struct {
	Value    *string        `json:"value, omitempty"`
	Metadata *ValueMetadata `json:"metadata,omitempty"`
}

// ValueMetadata the meta of value
type ValueMetadata struct {
	Timestamp int64 `json:"timestamp, omitempty"`
}

// TypeMetadata the meta of value type
type TypeMetadata struct {
	Type string `json:"type,omitempty"`
}

// TwinVersion twin version
type TwinVersion struct {
	CloudVersion int64 `json:"cloud"`
	EdgeVersion  int64 `json:"edge"`
}

// MsgTwin the struct of device twin
type MsgTwin struct {
	Expected        *TwinValue    `json:"expected,omitempty"`
	Actual          *TwinValue    `json:"actual,omitempty"`
	Optional        *bool         `json:"optional,omitempty"`
	Metadata        *TypeMetadata `json:"metadata,omitempty"`
	ExpectedVersion *TwinVersion  `json:"expected_version,omitempty"`
	ActualVersion   *TwinVersion  `json:"actual_version,omitempty"`
}

// DeviceTwinUpdate the struct of device twin update
type DeviceTwinUpdate struct {
	BaseMessage
	Twin map[string]*MsgTwin `json:"twin"`
}

// usage is responsible for setting up the default settings of all defined command-line flags for glog.
func usage() {
	flag.PrintDefaults()
	os.Exit(2)
}

// init for getting command line arguments for glog and initiating the MQTT connection
func init() {
	flag.Usage = usage
	// NOTE: This next line is key you have to call flag.Parse() for the command line
	// options or "flags" that are defined in the glog module to be picked up.
	flag.Parse()

	glog.Info("MQTT URL: ", MqttURL)
	// connect to mqtt broker
	connectToMqtt(MqttURL)

	err := LoadConfigMap()
	if err != nil {
		glog.Error(errors.New("Error while reading from config map " + err.Error()))
		os.Exit(1)
	}
}

func LoadConfigMap() error {
	var ok bool
	readConfigMap := configuration.DeviceProfile{}
	err := readConfigMap.ReadFromConfigMap()
	if err != nil {
		return errors.New("Error while reading from config map " + err.Error())
	}
	for _, device := range readConfigMap.DeviceInstances {
		if strings.ToUpper(device.Model) == modelName && strings.ToUpper(device.Name) == strings.ToUpper(DeviceName) {
			deviceID = device.ID
		}
	}
	for _, deviceModel := range readConfigMap.DeviceModels {
		if strings.ToUpper(deviceModel.Name) == modelName {
			for _, property := range deviceModel.Properties {
				if strings.ToUpper(property.Name) == pinNumberConfig {
					if pinNumber, ok = property.DefaultValue.(float64); ok == false {
						return errors.New(" Error in reading pin number from config map")
					}
					glog.Info("The pin number is ", pinNumber)
				}
			}
		}
	}
	return nil
}

// changeDeviceState function is used to change the state of the device
func changeDeviceState(state string) {
	glog.Info("Changing the state of the device to online")
	var deviceStateUpdateMessage DeviceStateUpdate
	deviceStateUpdateMessage.State = state
	stateUpdateBody, err := json.Marshal(deviceStateUpdateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	deviceStatusUpdate := DeviceETPrefix + deviceID + DeviceETStateUpdateSuffix
	Token_client = cli.Publish(deviceStatusUpdate, 0, false, stateUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device state update  is ", Token_client.Error())
	}
}

// changeTwinValue sends the updated twin value to the edge through the MQTT broker
func changeTwinValue(updateMessage DeviceTwinUpdate) {
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	deviceTwinUpdate := DeviceETPrefix + deviceID + TwinETUpdateSuffix
	Token_client = cli.Publish(deviceTwinUpdate, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device twin update is ", Token_client.Error())
	}
}

// OnSubMessageReceived callback function which is called when message is received
func OnSubMessageReceived(client MQTT.Client, message MQTT.Message) {
	err := json.Unmarshal(message.Payload(), &deviceTwinResult)
	// glog.Info("get data: ", message.Payload())
	if err != nil {
		glog.Error("Error in unmarshalling:  ", err)
	}
}

// createActualUpdateMessage function is used to create the device twin update message
func createActualUpdateMessage(actualValue string) DeviceTwinUpdate {
	var deviceTwinUpdateMessage DeviceTwinUpdate
	actualMap := map[string]*MsgTwin{powerStatus: {Actual: &TwinValue{Value: &actualValue}, Metadata: &TypeMetadata{Type: "Updated"}}}
	deviceTwinUpdateMessage.Twin = actualMap
	return deviceTwinUpdateMessage
}

// getTwin function is used to get the device twin details from the edge
func getTwin(updateMessage DeviceTwinUpdate) {
	getTwin := DeviceETPrefix + deviceID + TwinETGetSuffix
	twinUpdateBody, err := json.Marshal(updateMessage)
	if err != nil {
		glog.Error("Error:   ", err)
	}
	Token_client = cli.Publish(getTwin, 0, false, twinUpdateBody)
	if Token_client.Wait() && Token_client.Error() != nil {
		glog.Error("client.publish() Error in device twin get  is ", Token_client.Error())
	}
}

// subscribe function subscribes  the device twin information through the MQTT broker
func subscribe() {
	for {
		getTwinResult := DeviceETPrefix + deviceID + TwinETGetResultSuffix
		Token_client = cli.Subscribe(getTwinResult, 0, OnSubMessageReceived)
		if Token_client.Wait() && Token_client.Error() != nil {
			glog.Error("subscribe() Error in device twin result get  is ", Token_client.Error())
		}
		time.Sleep(1 * time.Second)
		if deviceTwinResult.Twin != nil {
			wg.Done()
			break
		}
	}
}

// equateTwinValue is responsible for equating the actual state of the device to the expected state that has been set
func equateTwinValue(updateMessage DeviceTwinUpdate) {
	glog.Info("Watching on the device twin values for device: ", DeviceName)
	wg.Add(1)
	go subscribe()
	getTwin(updateMessage)
	wg.Wait()
	if deviceTwinResult.Twin[powerStatus].Expected != nil && ((deviceTwinResult.Twin[powerStatus].Actual == nil) && deviceTwinResult.Twin[powerStatus].Expected != nil || (*deviceTwinResult.Twin[powerStatus].Expected.Value != *deviceTwinResult.Twin[powerStatus].Actual.Value)) {
		glog.Info("Expected Value : ", *deviceTwinResult.Twin[powerStatus].Expected.Value)
		if deviceTwinResult.Twin[powerStatus].Actual == nil {
			glog.Info("Actual Value: ", deviceTwinResult.Twin[powerStatus].Actual)
		} else {
			glog.Info("Actual Value: ", *deviceTwinResult.Twin[powerStatus].Actual.Value)
		}
		glog.Info("Equating the actual  value to expected value")
		switch strings.ToUpper(*deviceTwinResult.Twin[powerStatus].Expected.Value) {
		case "ON":
			glog.Info("Turning ON the light")
			//Turn On the light by supplying power on the pin specified
			lightdriver.TurnON(int64(pinNumber))

		case "OFF":
			glog.Info("Turning OFF the light")
			//Turn Off the light by cutting off power on the pin specified
			lightdriver.TurnOff(int64(pinNumber))

		default:
			panic("OOPS!!!!! Attempt to perform invalid operation " + *deviceTwinResult.Twin[powerStatus].Expected.Value + " on LED light")
		}
		updateMessage = createActualUpdateMessage(*deviceTwinResult.Twin[powerStatus].Expected.Value)
		changeTwinValue(updateMessage)
	} else {
		glog.Info("Actual values are in sync with Expected value")
	}
}

func connectToMqtt(mqttUrl string) MQTT.Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttUrl)

	cli = MQTT.NewClient(opts)

	token := cli.Connect()
	if token.Wait() && token.Error() != nil {
		glog.Info(token.Error())
	} else {
		glog.Info("Connected to MQTT broker")
	}

	return cli
}

func main() {
	changeDeviceState("online")
	updateMessage := createActualUpdateMessage("unknown")
	for {
		equateTwinValue(updateMessage)
	}
}
