package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tx7do/kratos-transport/broker"
	svcV1 "kratos-cqrs/api/logger/service/v1"
)

func (s *LoggerJobService) InsertSensorData(event broker.Event) error {
	fmt.Println("InsertSensorData() Topic: ", event.Topic(), " Payload: ", string(event.Message().Body))

	var sensorData []*svcV1.SensorData
	err := json.Unmarshal(event.Message().Body, &sensorData)
	if err != nil {
		s.log.Debug("InsertSensorData Unmarshal", err.Error())
		return err
	}

	err = s.sensorData.BatchInsertSensorData(context.Background(), sensorData)
	if err != nil {
		s.log.Debug("InsertSensorData Insert", err.Error())
		return err
	}

	err = event.Ack()
	if err != nil {
		s.log.Debug("InsertSensorData Ack", err.Error())
		return err
	}

	return nil
}

func (s *LoggerJobService) InsertSensor(event broker.Event) error {
	fmt.Println("InsertSensor() Topic: ", event.Topic(), " Payload: ", string(event.Message().Body))

	var sensor svcV1.Sensor
	err := json.Unmarshal(event.Message().Body, &sensor)
	if err != nil {
		s.log.Debug("InsertSensor Unmarshal", err.Error())
		return err
	}

	err = s.sensor.Create(context.Background(), &sensor)
	if err != nil {
		s.log.Debug("InsertSensor Insert", err.Error())
		return err
	}

	return nil
}
