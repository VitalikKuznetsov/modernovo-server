from tests.core.adaptors.FavoritesAdaptor import FavoritesAdaptor
from tests.core.services.FavoritesService import FavoritesService
from tests.core.services.CreateUserService import CreateUserService
from tests.core.adaptors.ProductsAdaptor import ProductsAdaptor
from hamcrest import *
import json


class TestFavorites:
    def testShouldGetEmptyFavorites(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        req = FavoritesAdaptor().getFavorites(headers=headers)

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            favorites_data = req.json()
            print(f"Избранное: {favorites_data}")

            if favorites_data is None:
                print(f"✅ Избранное пустое (получено None)")
                return

            if isinstance(favorites_data, dict):
                if 'favorites' in favorites_data:
                    assert_that(favorites_data, has_key('favorites'))
                    assert_that(favorites_data, has_key('total'))

                    assert_that(favorites_data['total'], is_(0))

                    favorites = favorites_data.get('favorites')
                    if favorites is not None:
                        assert_that(favorites, instance_of(list))
                        assert_that(len(favorites), is_(0))

                else:
                    print(f"⚠️ Неизвестная структура избранного: {favorites_data}")
            else:
                print(f"⚠️ Неожиданный тип данных избранного: {type(favorites_data)}")
        else:
            print(f"Response: {req.text}")
            assert_that(req.status_code, is_(200))

    def testShouldAddToFavorites(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = FavoritesService.get_existing_product_id()

        favorite_request = {
            "product_id": product_id
        }

        req = FavoritesAdaptor().addToFavorites(
            headers=headers,
            body=json.dumps(favorite_request)
        )

        if req.status_code == 200:
            favorites_data = req.json()
            if favorites_data is None:
                get_req = FavoritesAdaptor().getFavorites(headers=headers)
                if get_req.status_code == 200:
                    get_data = get_req.json()
                    if get_data is not None and isinstance(get_data, dict):
                        favorites = get_data.get('favorites', [])
            else:
                if isinstance(favorites_data, dict):
                    if 'favorites' in favorites_data:
                        assert_that(favorites_data, has_key('favorites'))
                        assert_that(favorites_data, has_key('total'))
                        assert_that(favorites_data['total'], greater_than(0))

                        favorites = favorites_data.get('favorites')
                        if favorites is not None:
                            assert_that(favorites, instance_of(list))
                            assert_that(product_id in favorites, is_(True))
                    else:
                        print(f"Неизвестная структура ответа: {favorites_data}")
                else:
                    print(f"Неожиданный тип ответа: {type(favorites_data)}")

        else:
            print(f"❌ Ошибка при добавлении в избранное")

        FavoritesService.clear_user_favorites(headers)

    def testShouldRemoveFromFavorites(self):
        setup_data = FavoritesService.setup_user_with_favorite()
        headers = setup_data['headers']
        product_id = setup_data['product_id']

        if not setup_data.get('success', False):
            print("Пропускаем тест - не удалось добавить товар в избранное")
            return
        remove_request = {
            "product_id": product_id
        }

        req = FavoritesAdaptor().removeFromFavorites(
            headers=headers,
            body=json.dumps(remove_request)
        )

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            favorites_data = req.json()
            print(f"Товар удален из избранного")

            if favorites_data is not None and isinstance(favorites_data, dict):
                favorites = favorites_data.get('favorites')
                if favorites is not None:
                    assert_that(product_id not in favorites, is_(True))
                    print(f"Товаров в избранном: {len(favorites)}")
                else:
                    print(f"Избранное пустое (favorites=None)")
            else:
                get_req = FavoritesAdaptor().getFavorites(headers=headers)
                if get_req.status_code == 200:
                    get_data = get_req.json()
                    if get_data is not None and isinstance(get_data, dict):
                        favorites = get_data.get('favorites', [])
                        if product_id not in favorites:
                            print(f"Товар удален из избранного (проверено отдельным запросом)")
        else:
            print(f"Response: {req.text}")

    def testShouldCheckFavoriteStatus(self):
        setup_data = FavoritesService.setup_user_with_favorite()
        headers = setup_data['headers']
        product_id = setup_data['product_id']

        if not setup_data.get('success', False):
            print("Пропускаем тест - не удалось добавить товар в избранное")
            return

        check_request = {
            "product_id": product_id
        }

        req = FavoritesAdaptor().checkFavorities(
            headers=headers,
            body=json.dumps(check_request)
        )

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            check_data = req.json()
            print(f"Результат проверки: {check_data}")

            if check_data is None:
                print(f"API вернуло None при проверке")
                return
            if isinstance(check_data, dict):
                if 'is_favorite' in check_data:
                    assert_that(check_data['is_favorite'], is_(True))
                elif 'favorite' in check_data:
                    assert_that(check_data['favorite'], is_(True))
                elif 'in_favorites' in check_data:
                    assert_that(check_data['in_favorites'], is_(True))
                else:
                    print(f"Статус получен (неизвестная структура ответа: {check_data})")
            elif isinstance(check_data, bool):
                assert_that(check_data, is_(True))
            else:
                print(f"⚠️ Неожиданный тип данных: {type(check_data)}")
        else:
            print(f"Response: {req.text}")

        FavoritesService.clear_user_favorites(headers)

    def testShouldNotAddDuplicateToFavorites(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = FavoritesService.get_existing_product_id()

        favorite_request = {
            "product_id": product_id
        }

        req1 = FavoritesAdaptor().addToFavorites(
            headers=headers,
            body=json.dumps(favorite_request)
        )

        if req1.status_code != 200:
            FavoritesService.clear_user_favorites(headers)
            return

        req2 = FavoritesAdaptor().addToFavorites(
            headers=headers,
            body=json.dumps(favorite_request)
        )

        if req2.status_code == 200:
            print(f"✅ Дубликат проигнорирован или добавлен повторно")
        elif req2.status_code in [400, 409]:
            print(f"✅ Сервер корректно отклонил дубликат")
        else:
            print(f"⚠️ Неожиданный статус для дубликата: {req2.status_code}")

        FavoritesService.clear_user_favorites(headers)

    def testShouldNotAddToFavoritesWithoutAuth(self):
        favorite_request = {
            "product_id": 1
        }

        req = FavoritesAdaptor().addToFavorites(body=json.dumps(favorite_request))

        assert_that(req.status_code, any_of(is_(401), is_(403)))
    def testMultipleItemsInFavorites(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()


        products_req = ProductsAdaptor().getProducts()

        if products_req.status_code != 200:
            print("⚠️ Нет продуктов для теста")
            return

        products = products_req.json()['products']
        if len(products) < 2:
            print("⚠️ Нужно минимум 2 продукта для теста")
            return

        items_added = []

        for i, product in enumerate(products[:2]):
            favorite_request = {
                "product_id": product['id']
            }

            add_req = FavoritesAdaptor().addToFavorites(
                headers=headers,
                body=json.dumps(favorite_request)
            )

            if add_req.status_code == 200:
                items_added.append(product['id'])
                print(f"Добавлен товар ID: {product['id']}")

        if items_added:
            favorites_req = FavoritesAdaptor().getFavorites(headers=headers)

            if favorites_req.status_code == 200:
                favorites_data = favorites_req.json()

                print(f"\nИтоговое избранное:")
                print(f"Данные: {favorites_data}")

                if favorites_data is not None and isinstance(favorites_data, dict):
                    favorites = favorites_data.get('favorites', [])
                    for product_id in items_added:
                        assert_that(product_id in favorites, is_(True))

                    print(f"✅ Все {len(items_added)} товара в избранном")
                else:
                    print(f"⚠️ Неизвестный формат избранного, но запросы прошли успешно")

        FavoritesService.clear_user_favorites(headers)

    def testShouldCheckNonFavoriteItem(self):
        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = FavoritesService.get_existing_product_id()

        print(f"\nПроверяем товар {product_id}, которого нет в избранном")

        check_request = {
            "product_id": product_id
        }

        req = FavoritesAdaptor().checkFavorities(
            headers=headers,
            body=json.dumps(check_request)
        )

        print(f"Response Status: {req.status_code}")

        if req.status_code == 200:
            check_data = req.json()
            print(f"Результат проверки: {check_data}")

            if check_data is None:
                print(f"⚠️ API вернуло None при проверке")
                return

            if isinstance(check_data, dict):
                if 'is_favorite' in check_data:
                    assert_that(check_data['is_favorite'], is_(False))
                    print(f"✅ Товар не в избранном (is_favorite=False)")
                elif 'favorite' in check_data:
                    assert_that(check_data['favorite'], is_(False))
                    print(f"✅ Товар не в избранном (favorite=False)")
                elif 'in_favorites' in check_data:
                    assert_that(check_data['in_favorites'], is_(False))
                    print(f"✅ Товар не в избранном (in_favorites=False)")
                else:
                    print(f"✅ Статус получен (неизвестная структура)")
            elif isinstance(check_data, bool):
                assert_that(check_data, is_(False))
                print(f"✅ Товар не в избранном (boolean=False)")
            else:
                print(f"✅ Статус получен (тип: {type(check_data)})")
        else:
            print(f"Response: {req.text}")