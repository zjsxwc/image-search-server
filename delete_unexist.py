import glob
import os
import requests

for pkl_path in sorted(glob.glob('static/feature/*.pkl')):
    print(pkl_path)
    product_id = os.path.splitext(os.path.basename(pkl_path))[0]
    print(product_id)
    img_path = 'static/img/' + product_id + ".jpg"
    req = requests.get("https://fyh88.com/product/" + product_id)
    if req.status_code != 200:
        print("delete " + product_id)
        if os.path.exists(pkl_path):
            os.remove(pkl_path)
        if os.path.exists(img_path):
            os.remove(img_path)