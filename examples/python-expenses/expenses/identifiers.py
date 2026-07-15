"""Simple (single-field) value objects used across the expense domain."""

import uuid
from dataclasses import dataclass


@dataclass(frozen=True)
class ReportID:
    _value: str

    def __post_init__(self) -> None:
        if not self._value:
            raise ValueError("report id must not be empty")

    @classmethod
    def generate(cls) -> "ReportID":
        return cls(str(uuid.uuid4()))

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True)
class ReportTitle:
    _value: str

    def __post_init__(self) -> None:
        if not self._value.strip():
            raise ValueError("report title must not be empty")

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True)
class ReceiptNumber:
    _value: str

    def __post_init__(self) -> None:
        if not self._value.strip():
            raise ValueError("receipt number must not be empty")

    def __str__(self) -> str:
        return self._value


@dataclass(frozen=True)
class Category:
    _value: str

    def __post_init__(self) -> None:
        if not self._value.strip():
            raise ValueError("category must not be empty")

    def __str__(self) -> str:
        return self._value
