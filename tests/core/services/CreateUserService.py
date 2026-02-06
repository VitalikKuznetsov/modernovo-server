from faker import Faker
from pydantic.v1.utils import to_lower_camel

from ..adaptors.RegisterAdaptor import RegisterAdaptor
from ..models.models import UserRequest
from ..adaptors.LoginAdaptor import LoginAdaptor

class CreateUserService:

    def __init__(self):
        self.faker = Faker()

    def getUser(self):
        email = self.faker.email()
        password = self.faker.password()
        return email, password, RegisterAdaptor().registerUser(body=UserRequest(email=email, password=password).json())

    def getEmail(self):
        email = self.faker.email()
        return email

    def getTokenHeaders(self):
        email, password, request = self.getUser()
        token = LoginAdaptor().userLogin(body=UserRequest(email=email, password=password).json()).json()['token']
        headers = {'cookie': f'session_token={token}'}
        return headers

    def getAdminTokenHeaders(self):
        RegisterAdaptor().registerUser(body=UserRequest(email="admin@mail.ru", password="admin").json())
        token = LoginAdaptor().userLogin(body=UserRequest(email="admin@mail.ru", password="admin").json()).json()['token']
        headers = {'cookie': f'session_token={token}'}
        return headers