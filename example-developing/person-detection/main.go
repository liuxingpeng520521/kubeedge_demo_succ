package main

import (
	"context"
	"encoding/json"
	"github/babydeng/kubeedge-example/person-detection/utils"
	"log"
	"strconv"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	devices "github.com/kubeedge/kubeedge/pkg/apis/devices/v1alpha2"
	"k8s.io/client-go/rest"
)

// The default namespace in which the device instance resides
var namespace = "default"

// The default status
var originCmd = "OFF"

// The CRD client used to patch the device instance.
var crdClient *rest.RESTClient

// DeviceStatus is used to patch device power-status
type DeviceStatus struct {
	Status devices.DeviceStatus `json:"status"`
}

func connectToMqtt() MQTT.Client {
	mqttUrl := "tcp://10.177.29.77:11883"
	log.Println("MQTT URL: ", mqttUrl)
	opts := MQTT.NewClientOptions()
	opts.AddBroker(mqttUrl)

	cli := MQTT.NewClient(opts)

	token := cli.Connect()
	if token.Wait() && token.Error() != nil {
		glog.Info(token.Error())
	} else {
		log.Println("Connected to MQTT broker")
	}

	return cli
}

// OnSubMessageReceived callback function which is called when message is received
func OnSubMessageReceived(client MQTT.Client, message MQTT.Message) {
	msg := string(message.Payload())
	log.Println("get data: ", msg)
	if msg == "IN" {
		UpdateDeviceTwinWithDesiredTrack("ON", "led-light-instance-01")
		UpdateDeviceTwinWithDesiredTrack("ON", "buzzer-instance-01")
	} else if msg == "OUT" {
		UpdateDeviceTwinWithDesiredTrack("OFF", "led-light-instance-01")
		UpdateDeviceTwinWithDesiredTrack("OFF", "buzzer-instance-01")
	}
}

func sub(client MQTT.Client) {
	topic := "/device/camera"
	Token_client := client.Subscribe(topic, 0, OnSubMessageReceived)
	if Token_client.Wait() && Token_client.Error() != nil {
		log.Println("subscribe() Error", Token_client.Error())
	}
}
func init() {

	// connect to mqtt broker
	cli := connectToMqtt()
	sub(cli)
	// Create a client to talk to the K8S API server to patch the device CRDs
	kubeConfig, err := utils.KubeConfig()
	if err != nil {
		log.Fatalf("Failed to create KubeConfig, error : %v", err)
	}
	log.Println("Get kubeConfig successfully")

	crdClient, err = utils.NewCRDClient(kubeConfig)
	if err != nil {
		log.Fatalf("Failed to create device crd client , error : %v", err)
	}
	log.Println("Get crdClient successfully")
}

func buildStatusWithDesiredTrack(cmd string) devices.DeviceStatus {
	metadata := map[string]string{
		"timestamp": strconv.FormatInt(time.Now().Unix()/1e6, 10),
		"type":      "string",
	}
	twins := []devices.Twin{{PropertyName: "power-status", Desired: devices.TwinProperty{Value: cmd, Metadata: metadata}, Reported: devices.TwinProperty{Value: cmd, Metadata: metadata}}}
	devicestatus := devices.DeviceStatus{Twins: twins}
	return devicestatus
}

// UpdateDeviceTwinWithDesiredTrack patches the desired state of
// the device twin with the command.
func UpdateDeviceTwinWithDesiredTrack(cmd string, deviceID string) bool {
	if cmd == originCmd {
		return true
	}

	status := buildStatusWithDesiredTrack(cmd)
	deviceStatus := &DeviceStatus{Status: status}
	body, err := json.Marshal(deviceStatus)
	if err != nil {
		log.Printf("Failed to marshal device status %v", deviceStatus)
		return false
	}
	result := crdClient.Patch(utils.MergePatchType).Namespace(namespace).Resource(utils.ResourceTypeDevices).Name(deviceID).Body(body).Do(context.TODO())
	if result.Error() != nil {
		log.Printf("Failed to patch device status %v of device %v in namespace %v \n error:%+v", deviceStatus, deviceID, namespace, result.Error())
		return false
	} else {
		log.Printf("Turn %s %s", cmd, deviceID)
	}
	originCmd = cmd

	return true
}

func main() {
	select {}
}
