import requests
from ..data.links import url
class CartAdaptor:

    def getCard(self, headers=None, body=None):
        return requests.request("GET", url + "cart", headers=headers, data=body)

    def addToCard(self, headers=None, body=None):
        return requests.request("POST", url + "cart", headers=headers, data=body)

    def updateCard(self, headers=None, body=None):
        return requests.request("PUT", url + "cart", headers=headers, data=body)

    def removeFromCard(self, headers=None, body=None):
        return requests.request("DELETE", url + "cart", headers=headers, data=body)

    def clearCard(self, headers=None, body=None):
        return requests.request("POST", url + "cart/clear", headers=headers, data=body)