#!/bin/bash

# Start containers
docker-compose -f environment/docker-compose-dev.yml down
echo "[TipsGO]: vetautet server stop..."