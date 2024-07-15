#!/bin/bash

python3 deamon_feature_extractor_process.py > deamon_feature_extractor_process.log  2>&1 &
python3 deamon_query_process.py > deamon_query_process.log  2>&1 &


# ./venv/bin/python3 deamon_feature_extractor_process.py > deamon_feature_extractor_process.log  2>&1 &
# ./venv/bin/python3 deamon_query_process.py > deamon_query_process.log  2>&1 &
