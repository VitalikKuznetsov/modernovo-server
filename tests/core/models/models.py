from typing import List, Optional

from pydantic import BaseModel

class UserRequest(BaseModel):
    email: str
    password: str

class UserEditRequest(BaseModel):
    name: str
    phonenumber: str

class ProductRequest(BaseModel):
    name: str
    description: str
    price: float
    ImageUrl: str
    image_urls: List[str]

    def json(self, **kwargs):
        return super().json(**kwargs)


class ProductUpdateRequest(BaseModel):
    name: Optional[str] = None
    description: Optional[str] = None
    price: Optional[float] = None
    ImageUrl: Optional[str] = None
    image_urls: Optional[List[str]] = None

    def json(self, **kwargs):
        return super().json(**kwargs)


class ProductResponse(BaseModel):
    id: int
    name: str
    description: str
    price: float
    ImageUrl: str
    image_urls: List[str]


class CartItemRequest(BaseModel):
    product_id: int
    quantity: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class UpdateCartItemRequest(BaseModel):
    product_id: int
    quantity: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class RemoveCartItemRequest(BaseModel):
    product_id: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class CartResponse(BaseModel):
    id: int
    user_id: int
    items: List[dict]
    total_price: float
    total_items: int


class FavoriteItemRequest(BaseModel):
    product_id: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class CheckFavoriteRequest(BaseModel):
    product_id: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class RemoveFavoriteRequest(BaseModel):
    product_id: int

    def json(self, **kwargs):
        return super().json(**kwargs)


class FavoriteResponse(BaseModel):
    favorites: List[int]
    total: int