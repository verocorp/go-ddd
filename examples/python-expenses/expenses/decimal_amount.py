"""DecimalAmount — an exact monetary amount as a value object.

The certus/`quanta.Decimal` pattern: a multi-representation primitive (a
``Decimal``, where ``1.5`` and ``1.50`` are the same value) is wrapped in its
own value object so callers compare it *by value*, never by string form. The
wrapped ``Decimal`` never leaves — the domain surfaces are the comparison and
arithmetic methods, and ``__str__`` is the sole serialization/display form.
"""

from dataclasses import dataclass
from decimal import Decimal, InvalidOperation


@dataclass(frozen=True)
class DecimalAmount:
    _value: Decimal

    @classmethod
    def parse(cls, raw: str) -> "DecimalAmount":
        """Construct from a decimal string — the boundary/spec entry point."""
        try:
            return cls(Decimal(raw))
        except InvalidOperation as e:
            raise ValueError(f"invalid amount: {raw!r}") from e

    def __post_init__(self) -> None:
        if not self._value.is_finite():
            raise ValueError(f"amount must be a finite number: {self._value}")

    def is_positive(self) -> bool:
        return self._value > 0

    def add(self, other: "DecimalAmount") -> "DecimalAmount":
        return DecimalAmount(self._value + other._value)

    def exceeds(self, other: "DecimalAmount") -> bool:
        """Value comparison — the domain surface for the report's total cap."""
        return self._value > other._value

    def __str__(self) -> str:
        # Serialization/display only. Equality is by value (Decimal), never by
        # this string — so a non-canonical form (1.5 vs 1.50) round-trips
        # faithfully without ever being compared as text.
        return str(self._value)
