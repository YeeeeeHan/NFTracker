import json
import sys
import requests
from bs4 import BeautifulSoup


URL = sys.argv[1]

headers = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
    # noqa
}

req = requests.Session()

page = req.get(URL, headers=headers)


soup = BeautifulSoup(page.text, 'html.parser')

data = json.loads(soup.find('script', type='application/json').text)
print(json.dumps(data['props']['relayCache'][0][1]["json"]["data"]["collection"], sort_keys=False, indent=4))