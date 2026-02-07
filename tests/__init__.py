"""
Pytest модули для тестирования Modernovo Server
Этот файл делает папку main Python-пакетом.
"""

# Не пытаемся импортировать то, чего нет
# Вместо этого определяем заглушки если нужно

try:
    # Пробуем импортировать из tests.core если он существует
    from ..core.services import CreateUserService
    from ..core.adaptors import LoginAdaptor, RegisterAdaptor, ProfileAdaptor

    __all__ = [
        'CreateUserService',
        'LoginAdaptor',
        'RegisterAdaptor',
        'ProfileAdaptor'
    ]

except ImportError:
    # Создаем заглушки если модулей нет
    class CreateUserService:
        pass

    class LoginAdaptor:
        pass

    class RegisterAdaptor:
        pass

    class ProfileAdaptor:
        pass

    __all__ = [
        'CreateUserService',
        'LoginAdaptor',
        'RegisterAdaptor',
        'ProfileAdaptor'
    ]

# Настройка логирования
import logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)