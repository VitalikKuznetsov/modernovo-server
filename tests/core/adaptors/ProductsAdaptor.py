import requests
from ..data.links import url
class ProductsAdaptor:

    def getProducts(self, headers=None, body=None):
        return requests.request("GET", url + "products", headers=headers, data=body)

    def getProductsById(self, headers=None, body=None, id=None):
        return requests.request("GET", url + f"products/{id}", headers=headers, data=body)

    def getProductsByIdDetail(self, headers=None, body=None, id=None):
        return requests.request("GET", url + f"products/{id}", headers=headers, data=body)