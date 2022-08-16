import requests


# update go constants from dydx-v3-python

def load_constants():
    url = "https://raw.githubusercontent.com/dydxprotocol/dydx-v3-python/master/dydx3/constants.py"
    res = requests.get(url, proxies={"https": 'http://127.0.0.1:1081'}).text
    with open('constants.py', 'w') as fp:
        fp.write(res)


def update_constants():
    import constants
    # generate go constants code
    text = "package starkex\n\n"
    text += "var ASSET_RESOLUTION = map[string]int64{\n"
    for item in constants.ASSET_RESOLUTION:
        text += '  "{}": {},\n'.format(item, constants.ASSET_RESOLUTION[item])
    text += "}\n\n"
    text += "var SYNTHETIC_ID_MAP = map[string]string{\n"
    for item in constants.SYNTHETIC_ASSET_ID_MAP:
        text += '  "{}": "{}",\n'.format(item, constants.SYNTHETIC_ASSET_ID_MAP[item])
    text += "}\n"
    print(text)
    with open('assets.go', 'w') as fp:
        fp.write(text)


load_constants()
update_constants()
