-- https://docs.timescale.com/timescaledb/latest/tutorials/simulate-iot-sensor-data/#step1

-- 创建表 - 传感器
CREATE TABLE sensors
(
    id       SERIAL PRIMARY KEY,
    type     VARCHAR(50),
    location VARCHAR(50)
);

-- 创建表 - 传感器遥感数据
CREATE TABLE sensor_data
(
    id          SERIAL,
    time        BIGINT NOT NULL,
    sensor_id   INTEGER,
    temperature DOUBLE PRECISION,
    cpu         DOUBLE PRECISION,
    FOREIGN KEY (sensor_id) REFERENCES sensors (id)
);
SELECT create_hypertable('sensor_data', 'time', chunk_time_interval => 86400000);

-- 插入测试数据 - 传感器
INSERT INTO sensors (type, location)
VALUES ('a', 'floor'),
       ('a', 'ceiling'),
       ('b', 'floor'),
       ('b', 'ceiling');

-- 插入测试数据 - 传感器遥感数据
INSERT INTO sensor_data (time, sensor_id, cpu, temperature)
SELECT time,
       sensor_id,
       random()       AS cpu,
       random() * 100 AS temperature
FROM generate_series(now() - interval '24 hour', now(), interval '5 minute') AS g1(time),
     generate_series(1, 4, 1) AS g2(sensor_id);

-- 查询所有数据 - 传感器
SELECT *
FROM sensors;

-- 查询所有数据 - 传感器遥感数据
SELECT *
FROM sensor_data
ORDER BY time;

-- 查询所有传感器[30分钟]间隔的[平均值]
SELECT time_bucket(1800000, time) AS period,
       AVG(temperature)                AS avg_temp,
       AVG(cpu)                        AS avg_cpu
FROM sensor_data
GROUP BY period;

-- 查询所有传感器[30分钟]间隔的[平均值]以及[最新温度值]
SELECT time_bucket(1800000, time) AS period,
       AVG(temperature)                AS avg_temp,
       AVG(cpu)                        AS avg_cpu,
       last(temperature, time)         AS last_temp
FROM sensor_data
GROUP BY period;

-- 连表查询所有传感器[30分钟]间隔的[平均值]以及[最新温度值]
SELECT sensors.location,
       time_bucket(1800000, time) AS period,
       AVG(temperature)                AS avg_temp,
       AVG(cpu)                        AS avg_cpu,
       last(temperature, time)         AS last_temp
FROM sensor_data
         JOIN sensors on sensor_data.sensor_id = sensors.id
GROUP BY period, sensors.location;
