"""
Core модули для тестирования
"""
from .services import CreateUserService
from .adaptors import LoginAdaptor, RegisterAdaptor, ProfileAdaptor

__all__ = [
    'CreateUserService',
    'LoginAdaptor',
    'RegisterAdaptor',
    'ProfileAdaptor'
]

print("Инициализация tests.core модулей...")