from __future__ import annotations

import pytest

from campaign.domain.campaign import Campaign
from campaign.domain.money import Money, MoneySpec
from campaign.domain.short_link import ShortLink
from campaign.domain.values import Slug
from serialization import canonical


def test_canonical_dispatches_the_single_defined_exit() -> None:
    assert canonical(Slug("promo"), str) == "promo"


def test_canonical_rejects_a_type_with_no_exit() -> None:
    with pytest.raises(TypeError, match="exactly one canonical exit"):
        canonical(Money(MoneySpec(amount="1.00", currency="USD")), str)


def test_canonical_rejects_a_type_with_two_exits() -> None:
    class _Two:
        def __str__(self) -> str:
            return "x"

        def __int__(self) -> int:
            return 1

    with pytest.raises(TypeError, match="exactly one canonical exit"):
        canonical(_Two(), str)


def test_canonical_rejects_a_mismatched_expectation() -> None:
    with pytest.raises(TypeError, match="not int"):
        canonical(Slug("promo"), int)

    class _IntBacked:
        def __int__(self) -> int:
            return 7

    assert canonical(_IntBacked(), int) == 7
    with pytest.raises(TypeError, match="not str"):
        canonical(_IntBacked(), str)


def test_canonical_rejects_an_exit_returning_the_wrong_type() -> None:
    class _Lying:
        def __int__(self) -> int:
            return "seven"  # type: ignore[return-value]

    with pytest.raises(TypeError, match="returned str, not int"):
        canonical(_Lying(), int)


def test_structured_types_define_no_conversion_dunders() -> None:
    for cls in (Money, ShortLink, Campaign):
        for name in ("__str__", "__int__", "__float__", "__bytes__"):
            assert name not in cls.__dict__, f"{cls.__name__} defines {name}"
