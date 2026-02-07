import requests
from ..data.links import url
class AdminProductsAdaptor:

    def getProducts(self, headers=None, body=None):
        return requests.request("GET", url + "admin/products", headers=headers, data=body)

    def makeProduct(self, headers=None, body=None):
        return requests.request("POST", url + "admin/products", headers=headers, data=body)

    def updateProductById(self, headers=None, body=None, id=None):
        return requests.request("PUT", url + f"admin/products/{id}", headers=headers, data=body)

    def deleteProductById(self, headers=None, body=None, id=None):
        return requests.request("DELETE", url + f"admin/products/{id}", headers=headers, data=body)