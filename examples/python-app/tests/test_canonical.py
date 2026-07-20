from __future__ import annotations

import pytest

from campaign.domain.campaign import Campaign
from campaign.domain.money import Money, MoneySpec
from campaign.domain.short_link import ShortLink
from campaign.domain.values import Slug
from serialization import canonical, canonical_text


def test_canonical_dispatches_the_single_defined_exit() -> None:
    assert canonical(Slug("promo")) == "promo"
    assert canonical_text(Slug("promo")) == "promo"


def test_canonical_rejects_a_type_with_no_exit() -> None:
    with pytest.raises(TypeError, match="exactly one canonical exit"):
        canonical(Money(MoneySpec(amount="1.00", currency="USD")))


def test_canonical_rejects_a_type_with_two_exits() -> None:
    class _Two:
        def __str__(self) -> str:
            return "x"

        def __int__(self) -> int:
            return 1

    with pytest.raises(TypeError, match="exactly one canonical exit"):
        canonical(_Two())


def test_canonical_text_rejects_a_non_text_exit() -> None:
    class _IntBacked:
        def __int__(self) -> int:
            return 7

    assert canonical(_IntBacked()) == 7
    with pytest.raises(TypeError, match="not str"):
        canonical_text(_IntBacked())


def test_structured_types_define_no_conversion_dunders() -> None:
    for cls in (Money, ShortLink, Campaign):
        for name in ("__str__", "__int__", "__float__", "__bytes__"):
            assert name not in cls.__dict__, f"{cls.__name__} defines {name}"
