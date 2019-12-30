import glob
import os
import pickle
from PIL import Image
from feature_extractor import FeatureExtractor

import os,shutil
import time

def movefile(srcfile,dstfile):
    if not os.path.isfile(srcfile):
        print(srcfile + " not exist!")
    else:
        fpath,fname=os.path.split(dstfile)    #分离文件名和路径
        if not os.path.exists(fpath):
            os.makedirs(fpath)                #创建路径
        shutil.move(srcfile,dstfile)          #移动文件
        print ("move "+srcfile+" -> "+ dstfile+"\n")


fe = FeatureExtractor()
while True:

	for img_path in sorted(glob.glob('static/processing-image/*.jpg')):
	    print(img_path)
	    img = Image.open(img_path)  # PIL image
	    feature = fe.extract(img)
	    feature_path = 'static/feature/' + os.path.splitext(os.path.basename(img_path))[0] + '.pkl'
	    pickle.dump(feature, open(feature_path, 'wb'))
	    movefile(img_path, img_path.replace("static/processing-image", "static/img"))

	time.sleep(5)

