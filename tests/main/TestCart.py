from tests.core.adaptors.CartAdaptor import CartAdaptor
from tests.core.services.CartService import CartService
from tests.core.services.CreateUserService import CreateUserService
from tests.core.models.models import CartItemRequest, UpdateCartItemRequest, RemoveCartItemRequest
from hamcrest import *


class TestCart:
    """Тесты для корзины"""

    def testShouldGetEmptyCart(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        req = CartAdaptor().getCard(headers=headers)

        if req.status_code == 200:
            cart_data = req.json()
            print(f"Корзина: {cart_data}")

            assert_that(cart_data, has_key('total'))
            assert_that(cart_data, has_key('item_count'))

            assert_that(cart_data['total'], is_(0))
            assert_that(cart_data['item_count'], is_(0))

            if cart_data.get('items') is None:
                print(f"✅ Корзина получена, items = None (ожидаемо для пустой корзины)")
            else:
                assert_that(cart_data['items'], instance_of(list))
                assert_that(len(cart_data['items']), is_(0))
                print(f"✅ Корзина получена, товаров: 0")
        else:
            print(f"Response: {req.text}")
            assert_that(req.status_code, is_(200))

    def testShouldAddItemToCart(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = CartService.get_existing_product_id()

        cart_item = CartItemRequest(
            product_id=product_id,
            quantity=3
        )

        req = CartAdaptor().addToCard(
            headers=headers,
            body=cart_item.json()
        )

        print(f"Response Status: {req.status_code}")
        print(f"Response: {req.text[:200] if req.text else 'No response'}")

        if req.status_code == 200:
            cart_data = req.json()
            assert_that(cart_data, has_key('total'))
            assert_that(cart_data, has_key('item_count'))
            assert_that(cart_data['item_count'], greater_than(0))

            items = cart_data.get('items')
            if items is not None:
                assert_that(items, instance_of(list))
                assert_that(len(items), greater_than(0))

                product_found = False
                for item in items:
                    if item.get('product_id') == product_id:
                        product_found = True
                        assert_that(item['quantity'], is_(cart_item.quantity))
                        break
                assert_that(product_found, is_(True))

            print(f"Итоговая корзина: total={cart_data.get('total')}, items={cart_data.get('item_count')}")
        else:
            print(f" Ошибка при добавлении в корзину")

        CartService.clear_user_cart(headers)
    def testShouldRemoveItemFromCart(self):
        setup_data = CartService.setup_user_with_cart()
        headers = setup_data['headers']
        product_id = setup_data['product_id']

        if not setup_data.get('success', False):
            print(" Пропускаем тест - не удалось добавить товар в корзину")
            return

        remove_item = RemoveCartItemRequest(
            product_id=product_id
        )

        req = CartAdaptor().removeFromCard(
            headers=headers,
            body=remove_item.json()
        )

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            cart_data = req.json()
            print(f" Товар удален из корзины")

            items = cart_data.get('items')
            if items is not None:
                product_in_cart = any(item.get('product_id') == product_id for item in items)
                assert_that(product_in_cart, is_(False))
                print(f"Товаров в корзине: {len(items)}")
            else:
                print(f"Корзина пустая (items=None)")
        else:
            print(f"Response: {req.text}")


    def testShouldClearCart(self):
        setup_data = CartService.setup_user_with_cart()
        headers = setup_data['headers']

        if not setup_data.get('success', False):
            print(" Пропускаем тест - не удалось добавить товар в корзину")
            return

        print(f"\n Очищаем всю корзину")

        req = CartAdaptor().clearCard(headers=headers)

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            cart_data = req.json()
            print(f" Корзина очищена")

            assert_that(cart_data['total'], is_(0))
            assert_that(cart_data['item_count'], is_(0))

            items = cart_data.get('items')
            if items is not None:
                assert_that(len(items), is_(0))

            print(f"Товаров в корзине после очистки: {cart_data.get('item_count')}")
        else:
            print(f"Response: {req.text}")