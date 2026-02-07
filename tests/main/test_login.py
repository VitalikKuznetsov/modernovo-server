from tests.core.models.models import UserRequest
from tests.core.services.CreateUserService import CreateUserService
from tests.core.adaptors.LoginAdaptor import LoginAdaptor
from hamcrest import *

class TestLogin:

    def testShouldLoginWithExistUser(self):
        email, password, request = CreateUserService().getUser()
        req = LoginAdaptor().userLogin(body=UserRequest(email=email, password=password).json())
        assert_that(req.status_code, is_(200))
        assert_that(req.json()['email'], is_(email))

    def testShouldNotLoginWithNoExistUser(self):
        req = LoginAdaptor().userLogin(body=UserRequest(email=CreateUserService().getEmail(), password="1234").json())
        assert_that(req.status_code, is_(401))
        assert_that(req.json()['error'], is_("Invalid email or password"))

    def testShouldNotLoginWithWrongPassword(self):
        email, password, request = CreateUserService().getUser()
        req = LoginAdaptor().userLogin(body=UserRequest(email=email, password="1234").json())
        assert_that(req.status_code, is_(401))
        assert_that(req.json()['error'], is_("Invalid email or password"))