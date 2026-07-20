from __future__ import annotations

_EXITS = ("__int__", "__float__", "__bytes__", "__str__")


def canonical(vo: object) -> str | int | float | bytes:
    cls = type(vo)
    defined = [name for name in _EXITS if name in cls.__dict__]
    if len(defined) != 1:
        raise TypeError(f"{cls.__name__} must define exactly one canonical exit, found {defined!r}")
    value = getattr(vo, defined[0])()
    if isinstance(value, (str, int, float, bytes)):
        return value
    raise TypeError(f"{cls.__name__}.{defined[0]} returned {type(value).__name__}")


def canonical_text(vo: object) -> str:
    value = canonical(vo)
    if not isinstance(value, str):
        raise TypeError(f"{type(vo).__name__} exits as {type(value).__name__}, not str")
    return value
