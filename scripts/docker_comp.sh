#!/bin/bash
./scripts/utils/func.sh
echo $get_curr_path
docker-compose -f docker-compose.yml up -d
