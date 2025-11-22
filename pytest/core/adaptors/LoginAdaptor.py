import requests
from ..data.links import url

class LoginAdaptor:

    def userLogin(self, headers = None, body = None):
        return requests.request("POST", url+"login", headers=headers, data=body)