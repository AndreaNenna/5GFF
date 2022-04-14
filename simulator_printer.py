import paho.mqtt.client as mqtt
import json
import time

mqtt_host = "172.17.48.170"
mqtt_port = 1883
mqtt_topic_device = "{ver}/Printers/{Id}/Device"
mqtt_topic_status = "{ver}/Printers/{Id}/Status"
mqtt_topic_current = "{ver}/Printers/{Id}/Task/Current"
mqtt_keepalive_interval = 60
mqtt_topic_current_finished = "{ver}/Printers/{Id}/Task/CurrentLayerFinished"


mqtt_msg_device = json.dumps({"Sn": "SN0000", "Model": "EP-M250A", "Type": "SLM", "Material": "stainless steel", "Size": [1,2,2], "Version": {"Driver": "1.0", "GUI": "1.0", "Remote": "1.0"}, "Shipping": "2018/3/1", "LicenseRemaining": "2019/8/2"})
mqtt_msg_status = json.dumps({"Status": 1})
mqtt_msg_current = json.dumps({"Data":{"Files":[{"File":"test","Matrix":[1,1,1,1,1]}], "Id": "{123}","Layers": 123,"Name": "{abc}","Operator": "USER","Size": [1,2,3],"Tickness": 1.234,"Times": {"Begin": 123456, "Remaining": 123}}})
mqtt_msg_current_finished = json.dumps({"Id": "1234", "LayerIndex": 7, "Name": "abc", "TimeCostInfo": {"Other": 100, "Pause": 100, "PrepareData": 100, "Scan": 100, "SpreadPower": 100}})
def on_publish(client, userdata, mid):
	print("Published")

def on_connect(client, userdata, flags, rc):
	client.subscribe(mqtt_topic_device, 0)
	client.publish(mqtt_topic_device, mqtt_msg_device)
	client.subscribe(mqtt_topic_status)
	client.publish(mqtt_topic_status, mqtt_msg_status)
	client.subscribe(mqtt_topic_current)
	client.publish(mqtt_topic_current, mqtt_msg_current)
	client.subscribe(mqtt_topic_current_finished, 0)
	client.publish(mqtt_topic_current_finished, mqtt_msg_current_finished)

def on_message(client, userdata, msg):
	print(msg.topic)
#	print(msg.payload)
	payload = json.loads(msg.payload)
	print(payload)


mqttc = mqtt.Client()
mqttc.loop_start()
while True:
	mqttc.on_publish = on_publish
	mqttc.on_connect = on_connect
	mqttc.on_message = on_message
	mqttc.connect(mqtt_host, mqtt_port, mqtt_keepalive_interval)
	time.sleep(10)

