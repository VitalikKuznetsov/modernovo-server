from tests.core.adaptors.ProductsAdaptor import ProductsAdaptor
from tests.core.services.ProductService import ProductService
from tests.core.services.CreateUserService import CreateUserService
from hamcrest import *


class TestProducts:

    def testShouldGetAllProducts(self):
        req = ProductsAdaptor().getProducts()

        assert_that(req.status_code, is_(200))
        response_data = req.json()

        assert_that(response_data, has_key('products'))
        assert_that(response_data, has_key('total'))
        assert_that(response_data['products'], instance_of(list))
        assert_that(response_data['total'], instance_of(int))

        assert_that(len(response_data['products']), is_(response_data['total']))

    def testShouldGetProductById(self):
        headers = CreateUserService().getAdminTokenHeaders()

        product_info = ProductService.create_test_product(headers)
        product_id = product_info['product_id']
        product_data = product_info['product_data']

        try:
            req = ProductsAdaptor().getProductsById(id=product_id)

            assert_that(req.status_code, is_(200))
            product = req.json()

            assert_that(product['id'], is_(product_id))
            assert_that(product['name'], is_(product_data.name))
            assert_that(product['description'], is_(product_data.description))
            assert_that(product['ImageUrl'], is_(product_data.ImageUrl))
            assert_that(product['image_urls'], is_(product_data.image_urls))

            assert_that(product['price'], instance_of((int, float)))

        finally:
            ProductService.cleanup_product(product_id, headers)

    def testShouldNotGetProductWithInvalidId(self):
        req = ProductsAdaptor().getProductsById(id="invalid_id_999999")

        assert_that(req.status_code, is_(400))

    def testShouldGetProductDetailById(self):
        headers = CreateUserService().getAdminTokenHeaders()

        product_info = ProductService.create_test_product(headers)
        product_id = product_info['product_id']
        product_data = product_info['product_data']

        try:
            req = ProductsAdaptor().getProductsByIdDetail(id=product_id)

            assert_that(req.status_code, is_(200))
            product = req.json()

            assert_that(product['id'], is_(product_id))
            assert_that(product['name'], is_(product_data.name))
            assert_that(product['description'], is_(product_data.description))
        finally:
            ProductService.cleanup_product(product_id, headers)

    def testShouldReturnProductsWithCorrectStructure(self):
        req = ProductsAdaptor().getProducts()

        assert_that(req.status_code, is_(200))
        response_data = req.json()
        products = response_data['products']

        if products:
            product = products[0]
            assert_that(product, has_key('id'))
            assert_that(product['id'], instance_of(int))
            assert_that(product, has_key('name'))
            assert_that(product['name'], instance_of(str))
            assert_that(product, has_key('description'))
            assert_that(product['description'], instance_of(str))
            assert_that(product, has_key('price'))

            price = product['price']
            assert_that(isinstance(price, (int, float)), is_(True),
                        f"Price should be int or float, got {type(price)}: {price}")

            assert_that(product, has_key('ImageUrl'))
            assert_that(product['ImageUrl'], instance_of(str))
            assert_that(product, has_key('image_urls'))
            assert_that(product['image_urls'], instance_of(list))

    def testShouldReturnEmptyProductsList(self):
        req = ProductsAdaptor().getProducts()

        assert_that(req.status_code, is_(200))
        response_data = req.json()

        assert_that(response_data['products'], instance_of(list))
        assert_that(response_data['total'], instance_of(int))

    def testShouldGetProductByExistingId(self):
        req = ProductsAdaptor().getProducts()
        assert_that(req.status_code, is_(200))

        products = req.json()['products']

        if products:
            product_id = products[0]['id']

            req_single = ProductsAdaptor().getProductsById(id=product_id)

            assert_that(req_single.status_code, is_(200))
            product = req_single.json()

            assert_that(product['id'], is_(product_id))
            assert_that(product['name'], is_(products[0]['name']))
            assert_that(product['price'], is_(products[0]['price']))

    def testShouldVerifyProductStructure(self):
        req = ProductsAdaptor().getProducts()
        assert_that(req.status_code, is_(200))

        products = req.json()['products']

        if not products:
            print("Нет продуктов для проверки структуры")
            return

        product = products[0]
        required_fields = {
            'id': int,
            'name': str,
            'description': str,
            'price': (int, float),
            'ImageUrl': str,
            'image_urls': list
        }

        for field, expected_type in required_fields.items():
            assert_that(product, has_key(field), f"Продукт должен содержать поле '{field}'")

            value = product[field]
            if isinstance(expected_type, tuple):
                assert_that(isinstance(value, expected_type), is_(True),
                            f"Поле '{field}' должно быть одним из типов {expected_type}, получен {type(value)}")
            else:
                assert_that(value, instance_of(expected_type),
                            f"Поле '{field}' должно быть типа {expected_type}, получен {type(value)}")

        for url in product['image_urls']:
            assert_that(url, instance_of(str))