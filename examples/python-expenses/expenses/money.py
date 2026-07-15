"""Money — a compound value object. Exact decimal amount, no float drift.

The amount is a ``DecimalAmount`` value object (not a raw ``Decimal``), so the
multi-representation primitive never leaks; the domain compares and sums via
methods. ``currency`` is a safe, single-representation value, exposed directly.
"""

from dataclasses import dataclass

from expenses.decimal_amount import DecimalAmount


@dataclass(frozen=True)
class MoneySpec:
    """Spec: primitive leaves only (the boundary carrier; public fields)."""

    amount: str
    currency: str


@dataclass(frozen=True)
class Money:
    _amount: DecimalAmount
    _currency: str

    @classmethod
    def from_spec(cls, spec: MoneySpec) -> "Money":
        try:
            amount = DecimalAmount.parse(spec.amount)
        except ValueError as e:
            raise ValueError(f"invalid money amount: {e}") from e
        return cls(amount, spec.currency)

    def __post_init__(self) -> None:  # the rules live here, always run
        if not self._currency:
            raise ValueError("currency is required")
        if len(self._currency) != 3 or not self._currency.isalpha():
            raise ValueError(f"invalid currency code: {self._currency!r}")

    @property
    def amount(self) -> DecimalAmount:
        """The amount as a value object — the domain surface for comparison and
        arithmetic. The raw Decimal never leaves the DecimalAmount."""
        return self._amount

    @property
    def currency(self) -> str:
        """The ISO currency code — a safe, single-representation value."""
        return self._currency

    def is_positive(self) -> bool:
        return self._amount.is_positive()

    def same_currency(self, other: "Money") -> bool:
        return self._currency == other._currency

    def add(self, other: "Money") -> "Money":
        if not self.same_currency(other):
            raise ValueError(
                f"cannot add {other._currency} to a {self._currency} amount"
            )
        return Money(self._amount.add(other._amount), self._currency)

    def __str__(self) -> str:
        return f"{self._amount} {self._currency}"
