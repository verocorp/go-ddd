from dataclasses import dataclass
from decimal import Decimal


@dataclass(frozen=True)
class Money:
    _amount: Decimal
    _currency: str


@dataclass(frozen=True)
class Currency:
    _value: str

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True)
class Price:
    _amount: Decimal
    _currency: Currency


@dataclass(frozen=True)
class Coordinate:
    _lat: float
    _lon: float


@dataclass(frozen=True)
class Window:
    _label: str
    _spans: tuple[int, ...]
