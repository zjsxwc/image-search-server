import os
import numpy as np
from PIL import Image
from feature_extractor import FeatureExtractor
import glob
import pickle
from datetime import datetime

import json
import time


def filePutContents(filename, contents):
      fh = open(filename, 'w', encoding='utf-8')
      fh.write(contents)
      fh.close()

# Read image features
fe = FeatureExtractor()
features = []
img_paths = []
def updateMemoryFeatures():
    for feature_path in glob.glob("static/feature/*.pkl"):
        img_path = 'static/img/' + os.path.splitext(os.path.basename(feature_path))[0] + '.jpg'
        if img_path not in img_paths:
            try:
                features.append(pickle.load(open(feature_path, 'rb')))
                img_paths.append(img_path)
            except:
                pass


updateMemoryFeatures()
currentTime = int(time.time())
lastUpdateTime = currentTime
while True:
    # 每隔60秒更新
    currentTime = int(time.time())
    if ((currentTime - lastUpdateTime) > 60):
        updateMemoryFeatures()
    if len(features) == 0:
        continue
    for ask_path in glob.glob("static/askans/*.ask.jpg"):
        ans_path = 'static/askans/' + os.path.splitext(os.path.splitext(os.path.basename(ask_path))[0])[0] + '.ans.json'
        if os.path.exists(ans_path): 
            continue
        try:
            img = Image.open(ask_path)  # PIL image
            query = fe.extract(img)
            del img
            dists = np.linalg.norm(features - query, axis=1)  # Do search
            ids = np.argsort(dists)[:30] # Top 30 results

            scores = [(str(dists[id]), img_paths[id], os.path.splitext(os.path.basename(img_paths[id]))[0]) for id in ids]
            filePutContents(ans_path, json.dumps(scores))
        except BaseException as e:
            print(e)
        else:
            pass