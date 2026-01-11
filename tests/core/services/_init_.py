"""
Сервисы для тестирования
"""

from .CreateUserService import CreateUserService

__all__ = ['CreateUserService']

def get_all_services():
    """Получить все доступные сервисы"""
    return {
        'create_user': CreateUserService
    }