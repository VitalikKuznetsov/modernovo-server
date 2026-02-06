from tests.core.models.models import UserRequest
from tests.core.services.CreateUserService import CreateUserService
from tests.core.adaptors.ProductsAdaptor import ProductsAdaptor
from hamcrest import *

class TestProducts:

    def testProducts(self):
        req = ProductsAdaptor().getProducts()
        print(req.json())