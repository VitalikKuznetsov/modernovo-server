"""
Адаптеры для API
"""

from .LoginAdaptor import LoginAdaptor
from .RegisterAdaptor import RegisterAdaptor
from .ProfileAdaptor import ProfileAdaptor
from .ProductsAdaptor import ProductsAdaptor
from .CartAdaptor import CartAdaptor
from .FavoritesAdaptor import FavoritesAdaptor
from .AdminProductsAdaptor import AdminProductsAdaptor

__all__ = ['LoginAdaptor', 'RegisterAdaptor', 'ProfileAdaptor', 'ProductsAdaptor', 'CartAdaptor', 'FavoritesAdaptor', 'AdminProductsAdaptor']