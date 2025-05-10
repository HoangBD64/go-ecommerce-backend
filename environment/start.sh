#!/bin/bash

# Start containers
docker-compose -f environment/docker-compose-dev.yml up
echo "[TipsGO]: vetautet server start..."