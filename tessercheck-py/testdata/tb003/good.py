from dataclasses import dataclass


@dataclass(frozen=True)
class Slug:
    _value: str

    def __post_init__(self) -> None:
        object.__setattr__(self, "_value", self._value.lower())


@dataclass(frozen=True)
class PersonNameSpec:
    given: str
    family: str


@dataclass(frozen=True, init=False)
class PersonName:
    _given: str
    _family: str

    def __init__(self, spec: PersonNameSpec) -> None:
        if not spec.given or not spec.family:
            raise ValueError("given and family are required")
        object.__setattr__(self, "_given", spec.given.strip())
        object.__setattr__(self, "_family", spec.family.strip())
