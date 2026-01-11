from faker import Faker
from ..adaptors.RegisterAdaptor import RegisterAdaptor
from ..models.models import UserRequest

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