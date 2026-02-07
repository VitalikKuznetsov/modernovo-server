import random
from ..adaptors.CartAdaptor import CartAdaptor
from ..models.models import CartItemRequest
from .CreateUserService import CreateUserService


class CartService:

    @staticmethod
    def get_existing_product_id():
        from ..adaptors.ProductsAdaptor import ProductsAdaptor

        try:
            products_response = ProductsAdaptor().getProducts()
            if products_response.status_code == 200:
                products = products_response.json()['products']
                if products:
                    print(f"Найден существующий продукт ID: {products[0]['id']}")
                    return products[0]['id']
        except Exception as e:
            print(f" Ошибка при получении продуктов: {e}")

        print(" Не найдено существующих продуктов, используем ID=1")
        return 1

    @staticmethod
    def create_test_cart_item(product_id=None):
        if product_id is None:
            product_id = CartService.get_existing_product_id()

        return CartItemRequest(
            product_id=product_id,
            quantity=random.randint(1, 5)
        )

    @staticmethod
    def setup_user_with_cart():

        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = CartService.get_existing_product_id()

        cart_item = CartItemRequest(
            product_id=product_id,
            quantity=2
        )

        cart_adaptor = CartAdaptor()
        add_response = cart_adaptor.addToCard(
            headers=headers,
            body=cart_item.json()
        )

        print(f"Response Status: {add_response.status_code}")

        if add_response.status_code != 200:
            print(f" Ошибка при добавлении в корзину: {add_response.text}")
            return {
                'headers': headers,
                'product_id': product_id,
                'cart_item': cart_item,
                'add_response': add_response,
                'user_service': user_service,
                'success': False
            }


        return {
            'headers': headers,
            'product_id': product_id,
            'cart_item': cart_item,
            'add_response': add_response,
            'user_service': user_service,
            'success': True
        }

    @staticmethod
    def clear_user_cart(headers):
        try:
            cart_adaptor = CartAdaptor()
            response = cart_adaptor.clearCard(headers=headers)
            print(f"Clear cart response: {response.status_code}")
            return response
        except Exception as e:
            print(f" Ошибка при очистке корзины: {e}")
            return None

    @staticmethod
    def get_cart_info(headers):
        try:
            cart_adaptor = CartAdaptor()
            response = cart_adaptor.getCard(headers=headers)

            if response.status_code == 200:
                cart_data = response.json()
                print(f" Корзина: total={cart_data.get('total')}, items={cart_data.get('item_count')}")
                return cart_data
            else:
                print(f" Ошибка получения корзины: {response.status_code}")
                return None
        except Exception as e:
            print(f" Ошибка при получении корзины: {e}")
            return None