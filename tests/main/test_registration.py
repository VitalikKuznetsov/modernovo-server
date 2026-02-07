from tests.core.models.models import UserRequest
from tests.core.services.CreateUserService import CreateUserService
from tests.core.adaptors.RegisterAdaptor import RegisterAdaptor
from hamcrest import *

class TestRegistration:

    def testShouldRegisterUser(self):
        email, password, request = CreateUserService().getUser()
        assert_that(request.status_code, is_(200))

    def testShouldErrorAfterDoubleRegister(self):
        email, password, request = CreateUserService().getUser()
        req = RegisterAdaptor().registerUser(body=UserRequest(email=email, password=password).json())
        assert_that(req.status_code, is_(409))
        assert_that(req.json(), "<{'error': 'User with this email already exists'}>")

    def testShoulErrorWithNoEmailRegister(self):
        req = RegisterAdaptor().registerUser(body=UserRequest(email="", password="1234").json())
        assert_that(req.status_code, is_(400))

    def testShoulErrorWithNoPasswordRegister(self):
        req = RegisterAdaptor().registerUser(body=UserRequest(email=CreateUserService().getEmail(), password="").json())
        assert_that(req.status_code, is_(400))

    def testShouldNotRegisterWithNoValidEmail(self):
        req = RegisterAdaptor().registerUser(body=UserRequest(email=CreateUserService().getEmail()[1::3], password="password").json())
        assert_that(req.status_code, is_(409))
        assert_that(req.json(), "<{'error': 'User with this email already exists'}>")