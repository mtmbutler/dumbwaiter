import json

import requests


class Dumbwaiter:
    def __init__(self, url: str):
        self.url = url

    def all_days(self):
        url = self.url + "/days"
        return requests.get(url)

    def new_day(self, day):
        url = self.url + "/days"
        return requests.post(url, data=json.dumps(day))

    def update_day(self, i, day):
        url = self.url + f"/days/{i}"
        return requests.put(url, data=json.dumps(day))

    def delete_day(self, i):
        url = self.url + f"/days/{i}"
        return requests.delete(url)

    def one_day(self, i):
        url = self.url + f"/days/{i}"
        return requests.get(url)
