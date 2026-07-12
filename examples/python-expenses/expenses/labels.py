"""Labels — a collection value object wrapping freeform key/value pairs."""

from dataclasses import dataclass
from typing import Mapping


@dataclass(frozen=True)
class Labels:
    _values: tuple[tuple[str, str], ...] = ()  # immutable, hashable storage

    def __post_init__(self) -> None:  # canonicalize on EVERY construction path
        # dedupe keys (last wins) AND sort — review-hardened so the raw
        # constructor can't hold a non-canonical value with duplicate keys.
        object.__setattr__(self, "_values", tuple(sorted(dict(self._values).items())))

    @classmethod
    def new(cls, values: Mapping[str, str] | None = None) -> "Labels":
        return cls(tuple((values or {}).items()))

    def as_dict(self) -> dict[str, str]:  # copy out, never a reference
        return dict(self._values)

    def __len__(self) -> int:
        return len(self._values)
