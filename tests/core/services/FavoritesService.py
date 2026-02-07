import json
from ..adaptors.FavoritesAdaptor import FavoritesAdaptor
from .CreateUserService import CreateUserService


class FavoritesService:

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

        print("Не найдено существующих продуктов, используем ID=1")
        return 1

    @staticmethod
    def setup_user_with_favorite():
        print(f"\n Настраиваем пользователя с товаром в избранном")

        user_service = CreateUserService()
        headers = user_service.getTokenHeaders()

        product_id = FavoritesService.get_existing_product_id()

        print(f"Добавляем товар {product_id} в избранное")

        favorites_adaptor = FavoritesAdaptor()
        favorite_request = {"product_id": product_id}

        add_response = favorites_adaptor.addToFavorites(
            headers=headers,
            body=json.dumps(favorite_request)
        )

        print(f"Response Status: {add_response.status_code}")

        if add_response.status_code != 200:
            print(f"Ошибка при добавлении в избранное: {add_response.text}")
            return {
                'headers': headers,
                'product_id': product_id,
                'add_response': add_response,
                'user_service': user_service,
                'success': False
            }

        print(f" Товар успешно добавлен в избранное")

        return {
            'headers': headers,
            'product_id': product_id,
            'add_response': add_response,
            'user_service': user_service,
            'success': True
        }

    @staticmethod
    def clear_user_favorites(headers):
        try:
            favorites_adaptor = FavoritesAdaptor()
            favorites_response = favorites_adaptor.getFavorites(headers=headers)

            if favorites_response.status_code == 200:
                favorites_data = favorites_response.json()
                favorites = favorites_data.get('favorites', [])

                for product_id in favorites:
                    remove_request = {"product_id": product_id}
                    favorites_adaptor.removeFromFavorites(
                        headers=headers,
                        body=json.dumps(remove_request)
                    )

                print(f"Удалено {len(favorites)} товаров из избранного")
        except Exception as e:
            print(f"Ошибка при очистке избранного: {e}")

    @staticmethod
    def get_favorites_info(headers):
        try:
            favorites_adaptor = FavoritesAdaptor()
            response = favorites_adaptor.getFavorites(headers=headers)

            if response.status_code == 200:
                favorites_data = response.json()
                print(f" Избранное: total={favorites_data.get('total')}, items={favorites_data.get('favorites')}")
                return favorites_data
            else:
                print(f" Ошибка получения избранного: {response.status_code}")
                return None
        except Exception as e:
            print(f"Ошибка при получении избранного: {e}")
            return None