import paho.mqtt.client as mqtt
import paho.mqtt.publish as publish
import json

mqtt_host = "172.17.48.170"
mqtt_port = 1883
mqtt_host_sub = "172.17.48.170"
mqtt_port_sub = 1883
mqtt_topic_device = "{ver}/Printers/{Id}/Device"
mqtt_topic_device_pub = "{ver}/Printers/{Id}/Device/pub"
#mqtt_topic_status = "{ver}/Printers/{Id}/Status"
#mqtt_topic_status_pub = "{ver}/Printers/{Id}/Status/pub"
#mqtt_topic_current = "{ver}/Printers/{Id}/Task/Current"
#mqtt_topic_current_pub = "{ver}/Printers/{Id}/Task/Current/pub"
mqtt_keepalive_interval = 60
client_name = "client_device"

def on_connect(client, userdata, flags, rc):
        client.subscribe(mqtt_topic_device, 0)
#        client.publish(mqtt_topic_state, mqtt_msg_state)
#        client.subscribe(mqtt_topic_status)
#        client.publish(mqtt_topic_status, mqtt_msg_status)
#        client.subscribe(mqtt_topic_current)
#        client.publish(mqtt_topic_current, mqtt_msg_current)

def on_message(client, userdata, msg):
	global m_in
	global m_out
	m_decode=str(msg.payload.decode("utf-8","ignore"))
	print("ricevuto")
	m_in=json.loads(m_decode)
	m_in='%s'%m_in
	m_out={'name': 'DISPOSITIVO_DEVICE_PRINTER',
		'cmd': 'jsonDEVICE',
		'jsonDEVICE': m_in}
	data_out = json.dumps(m_out)
	publish_mqtt(data_out)
	print('pubblicato')
	print(data_out)

def publish_mqtt(sensor_data):
	mqttc = mqtt.Client('Invio_Dati')
	mqttc.connect(mqtt_host, mqtt_port)
	mqttc.publish(mqtt_topic_device_pub, sensor_data)

def on_publish(client, userdata, mid):
	print("Messaggio pubblicato")

client = mqtt.Client(client_name)
client.on_connect = on_connect
client.on_message = on_message 

client.connect(mqtt_host_sub, mqtt_port_sub, mqtt_keepalive_interval)
client.loop_forever()

