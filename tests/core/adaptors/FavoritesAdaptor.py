import requests
from ..data.links import url
class FavoritesAdaptor:
    def getFavorites(self, headers=None, body=None):
        return requests.request("GET", url + "favorites", headers=headers, data=body)

    def addToFavorites(self, headers=None, body=None):
        return requests.request("POST", url + "favorites", headers=headers, data=body)

    def checkFavorities(self, headers=None, body=None):
        return requests.request("GET", url + "favorites/check", headers=headers, data=body)

    def removeFromFavorites(self, headers=None, body=None):
        return requests.request("DELETE", url + "favorites", headers=headers, data=body)