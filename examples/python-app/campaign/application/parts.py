from __future__ import annotations

from dataclasses import dataclass

from campaign.domain.campaign import Campaign
from campaign.domain.short_link import ShortLink
from serialization import canonical


@dataclass(frozen=True)
class MoneyParts:
    amount: str
    currency: str


@dataclass(frozen=True)
class ShortLinkParts:
    slug: str
    target_url: str
    active: bool


@dataclass(frozen=True)
class CampaignParts:
    id: str
    budget: MoneyParts
    links: tuple[ShortLinkParts, ...]


def campaign_parts(c: Campaign) -> CampaignParts:
    return CampaignParts(
        id=canonical(c.id, str),
        budget=MoneyParts(
            amount=canonical(c.budget.amount, str),
            currency=canonical(c.budget.currency, str),
        ),
        links=tuple(_short_link_parts(link) for link in c.links),
    )


def _short_link_parts(link: ShortLink) -> ShortLinkParts:
    return ShortLinkParts(
        slug=canonical(link.slug, str),
        target_url=canonical(link.target_url, str),
        active=link.active,
    )
