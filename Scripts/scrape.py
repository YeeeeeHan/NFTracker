import json

import requests
from bs4 import BeautifulSoup

URL = "https://opensea.io/collection/doodles-official"
headers = {
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
    # noqa
}

req = requests.Session()

page = req.get(URL, headers=headers)


soup = BeautifulSoup(page.text, 'html.parser')
# print(soup.prettify())

data = json.loads(soup.find('script', type='application/json').text)
print(json.dumps(data['props']['relayCache'][0][1]["json"]["data"]["collection"], sort_keys=False, indent=4))
# sub = soup.findAll('span', class_='nav-item nav-item-section')[0].a.contents[0]
# module = soup.findAll('span', class_='nav-item nav-item-chapter')[0].a.contents[0]
