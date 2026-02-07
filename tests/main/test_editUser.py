from tests.core.services.CreateUserService import CreateUserService
from tests.core.adaptors.ProfileAdaptor import ProfileAdaptor
from faker import Faker
from tests.core.models.models import UserEditRequest, UserRequest
from tests.core.adaptors.LoginAdaptor import LoginAdaptor
from hamcrest import *

class TestEditUser:

    def testShouldCreateUserData(self):
        email, password, request = CreateUserService().getUser()
        token = LoginAdaptor().userLogin(body=UserRequest(email=email, password=password).json()).json()['token']
        headers = { 'cookie' : f'session_token={token}'}
        req = ProfileAdaptor().getUserData(headers=headers)
        assert_that(req.status_code, is_(200))
        fake = Faker()
        newName, newPhone = fake.name(), fake.phone_number()
        req = ProfileAdaptor().changeUserData(headers=headers, body=UserEditRequest(name=newName, phonenumber=newPhone).json())
        assert_that(req.status_code, is_(200))
        req = ProfileAdaptor().getUserData(headers=headers)
        assert_that(req.status_code, is_(200))
        assert_that(req.json()['name'], newName)
        assert_that(req.json()['phonenumber'], newPhone)

    def testShouldChangeUserData(self):
        email, password, request = CreateUserService().getUser()
        token = LoginAdaptor().userLogin(body=UserRequest(email=email, password=password).json()).json()['token']
        headers = {'cookie': f'session_token={token}'}
        req = ProfileAdaptor().getUserData(headers=headers)
        assert_that(req.status_code, is_(200))
        fake = Faker()
        fName, fPhone = fake.name(), fake.phone_number()
        sName, sPhone = fake.name(), fake.phone_number()
        req = ProfileAdaptor().changeUserData(headers=headers,
                                              body=UserEditRequest(name=fName, phonenumber=fPhone).json())
        assert_that(req.status_code, is_(200))
        req = ProfileAdaptor().getUserData(headers=headers)
        assert_that(req.status_code, is_(200))
        assert_that(req.json()['name'], fName)
        assert_that(req.json()['phonenumber'], fPhone)
        req = ProfileAdaptor().changeUserData(headers=headers,
                                              body=UserEditRequest(name=sName, phonenumber=sPhone).json())
        assert_that(req.status_code, is_(200))
        req = ProfileAdaptor().getUserData(headers=headers)
        assert_that(req.status_code, is_(200))
        assert_that(req.json()['name'], sName)
        assert_that(req.json()['phonenumber'], sPhone)