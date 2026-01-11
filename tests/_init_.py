"""
Pytest модули для тестирования Modernovo Server
"""

from .core.services import CreateUserService
from .core.adaptors import LoginAdaptor, RegisterAdaptor, ProfileAdaptor

__all__ = [
    'CreateUserService',
    'LoginAdaptor',
    'RegisterAdaptor',
    'ProfileAdaptor'
]

import logging

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)