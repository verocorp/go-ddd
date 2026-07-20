from dataclasses import dataclass


@dataclass(frozen=True)
class Money:
    _amount: str
    currency: str

    def __str__(self) -> str:
        return f"{self._amount} {self.currency}"


@dataclass(frozen=True)
class Slot:
    _key: str

    @property
    def key(self) -> str:
        return self._key

    def raw_key(self) -> str:
        return self._key
