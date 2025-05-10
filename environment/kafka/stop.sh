#!/bin/bash

# Stop and remove containers, networks, and volumes defined in docker-compose-cluster.yml
docker-compose -f environment/kafka/docker-compose-kafka-single.yml down
echo "[TipsGO]: Kafka single stop..."