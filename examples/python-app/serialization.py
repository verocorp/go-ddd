from __future__ import annotations

from typing import TypeVar

T = TypeVar("T", str, int, float, bytes)

_EXITS = {"__str__": str, "__int__": int, "__float__": float, "__bytes__": bytes}


def canonical(vo: object, expected: type[T]) -> T:
    cls = type(vo)
    defined = [name for name in _EXITS if name in cls.__dict__]
    if len(defined) != 1:
        raise TypeError(f"{cls.__name__} must define exactly one canonical exit, found {defined!r}")
    if _EXITS[defined[0]] is not expected:
        raise TypeError(f"{cls.__name__} defines {defined[0]}; its canonical form is not {expected.__name__}")
    value = getattr(vo, defined[0])()
    if not isinstance(value, expected):
        raise TypeError(f"{cls.__name__}.{defined[0]} returned {type(value).__name__}, not {expected.__name__}")
    return value
