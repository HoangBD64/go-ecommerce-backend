#!/bin/bash

# Start containers
docker-compose -f environment/kafka/docker-compose-kafka-single.yml up
echo "[TipsGO]: Kafka single start..."