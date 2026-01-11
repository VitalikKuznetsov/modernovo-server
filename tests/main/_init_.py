"""
Основные тесты
"""

from .TestRegistration import TestRegistration
from .TestLogin import TestLogin

__all__ = [
    'TestRegistration',
    'TestLogin'
]

import os
import sys

test_modules = []
for file in os.listdir(os.path.dirname(__file__)):
    if file.startswith('test_') and file.endswith('.py'):
        module_name = file[:-3]
        test_modules.append(module_name)

__all__.extend(test_modules)