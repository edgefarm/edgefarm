all: ../sensor_pb2.py 

../sensor_pb2.py : sensor.proto 
	~/bin/protoc sensor.proto --python_out=..

.PHONY: all 
