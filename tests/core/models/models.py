from pydantic import BaseModel

class UserRequest(BaseModel):
    email: str
    password: str

class UserEditRequest(BaseModel):
    name: str
    phonenumber: str