"""Shared AST helpers used by both the classifier and the checks.

Kept in their own module so ``classify`` and ``checks`` can each depend on them
without depending on each other.
"""

import ast


def _name_of(node: ast.expr) -> str | None:
    """The bare name of a ``Name``/``Attribute`` decorator target."""
    if isinstance(node, ast.Name):
        return node.id
    if isinstance(node, ast.Attribute):
        return node.attr
    return None


def _is_true(node: ast.expr) -> bool:
    return isinstance(node, ast.Constant) and node.value is True


def _dataclass_frozen(decorators: list[ast.expr]) -> tuple[bool, bool, ast.expr | None]:
    """(is_dataclass, is_frozen, decorator_node) for a class's decorator list.

    The decorator node is where TB001 points and where a ``# tessercheck:ignore`` is
    expected — it's the ``@dataclass`` line the user would change to add
    ``frozen=True``, not the ``class`` line the ClassDef reports.
    """
    for dec in decorators:
        target = dec.func if isinstance(dec, ast.Call) else dec
        if _name_of(target) != "dataclass":
            continue
        frozen = False
        if isinstance(dec, ast.Call):
            for kw in dec.keywords:
                if kw.arg == "frozen" and _is_true(kw.value):
                    frozen = True
        return True, frozen, dec
    return False, False, None


def _dataclass_init_false(decorators: list[ast.expr]) -> bool:
    """True when the ``@dataclass`` decorator declares ``init=False`` — the
    explicit signal that the class hand-writes its own construction path (the
    spec-taking ``__init__`` of a compound value object / entity)."""
    for dec in decorators:
        target = dec.func if isinstance(dec, ast.Call) else dec
        if _name_of(target) != "dataclass":
            continue
        if isinstance(dec, ast.Call):
            for kw in dec.keywords:
                if (
                    kw.arg == "init"
                    and isinstance(kw.value, ast.Constant)
                    and kw.value.value is False
                ):
                    return True
        return False
    return False


def _annotation_base(ann: ast.expr) -> str | None:
    """Base name of an annotation: ``list[X]``/``List[X]``/``list`` -> ``list``."""
    if isinstance(ann, ast.Name):
        return ann.id
    if isinstance(ann, ast.Subscript):
        return _annotation_base(ann.value)
    if isinstance(ann, ast.Attribute):
        return ann.attr
    return None


def _is_str_call(node: ast.expr) -> bool:
    """``str(x)`` or ``x.__str__()`` — a stringification, not a value."""
    if not isinstance(node, ast.Call):
        return False
    func = node.func
    if isinstance(func, ast.Name) and func.id == "str" and len(node.args) == 1:
        return True
    if isinstance(func, ast.Attribute) and func.attr == "__str__":
        return True
    return False
