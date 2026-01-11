import requests
from ..data.links import url

class ProfileAdaptor:

    def changeUserData(self, headers = None, body = None):
        return requests.request("PUT", url+"profile", headers=headers, data=body)

    def getUserData(self, headers = None, body = None):
        return requests.request("GET", url+"profile", headers=headers, data=body)