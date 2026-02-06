import random
import string
from ..adaptors.AdminProductsAdaptor import AdminProductsAdaptor
from ..models.models import ProductRequest


class ProductService:

    @staticmethod
    def generate_product_name():
        random_str = ''.join(random.choices(string.ascii_lowercase, k=6))
        return f"Продукт_{random_str}"

    @staticmethod
    def generate_description():
        adjectives = ["Серый", "Красный", "Синий", "Зеленый", "Черный", "Белый"]
        products = ["конвектор", "обогреватель", "вентилятор", "радиатор", "термостат"]
        return f"Описание {random.choice(adjectives).lower()} {random.choice(products)}а"

    @staticmethod
    def get_sample_product():
        return ProductRequest(
            name=ProductService.generate_product_name(),
            description=ProductService.generate_description(),
            price=round(random.uniform(1000.0, 5000.0), 2),
            ImageUrl=f"RedHeart/img_{random.randint(1, 10)}.png",
            image_urls=[
                f"RedHeart/img_{random.randint(1, 10)}.png",
                f"RedHeart/img_{random.randint(11, 20)}.png"
            ]
        )

    @staticmethod
    def create_test_product(headers):
        product_data = ProductService.get_sample_product()
        adaptor = AdminProductsAdaptor()

        response = adaptor.makeProduct(
            headers=headers,
            body=product_data.json()
        )

        if response.status_code == 201:
            created_product = response.json()
            return {
                'product_id': created_product["id"],
                'product_data': product_data,
                'created_product': created_product,
                'response': response
            }
        else:
            raise Exception(f"Failed to create product: {response.status_code} - {response.text}")

    @staticmethod
    def cleanup_product(product_id, headers):
        try:
            AdminProductsAdaptor().deleteProductById(
                headers=headers,
                id=product_id
            )
            return True
        except Exception as e:
            print(f"Cleanup error: {e}")
            return False