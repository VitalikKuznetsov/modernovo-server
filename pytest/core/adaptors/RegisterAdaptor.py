import requests
from ..data.links import url

class RegisterAdaptor:

    def registerUser(self, headers = None, body = None):
        return requests.request("POST", url+"register", headers=headers, data=body)