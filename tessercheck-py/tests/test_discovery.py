"""The ported discovery classifier's guarantees — the same teeth as the
verified impl's own gate (``examples/python-app/tests/test_discovery.py``),
proven against synthetic trees, plus the generalizations the port adds
(non-Python dirs outside the contract, parameterized app-level set,
namespace/broken ``__init__`` can't hide, no-contexts failure)."""

import pathlib

import pytest

from tessercheck.cli import main
from tessercheck.discovery import (
    APP_LEVEL_PACKAGES,
    classify_root,
    exposes_client,
    totality_errors,
)


def _pkg(root: pathlib.Path, name: str, init: str = "") -> pathlib.Path:
    d = root / name
    d.mkdir()
    (d / "__init__.py").write_text(init, encoding="utf-8")
    return d


def test_totality_guard_teeth_flags_clientless_context(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "billing", '"""a context that forgot its Client"""\n')
    d = classify_root(tmp_path)
    assert d.unclassified == ("billing",)
    assert d.contexts == ()


def test_discovery_teeth_finds_client_bearing_context(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "billing", "from billing.client import Client\n")
    d = classify_root(tmp_path)
    assert d.contexts == ("billing",)
    assert d.unclassified == ()


def test_exposes_client_detects_direct_definition(tmp_path: pathlib.Path) -> None:
    _pkg(
        tmp_path,
        "billing",
        "from typing import Protocol\n\nclass Client(Protocol):\n    def ping(self) -> None: ...\n",
    )
    assert exposes_client(tmp_path / "billing")


def test_exposes_client_detects_asname_reexport(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "billing", "from billing.impl import ApiClient as Client\n")
    assert exposes_client(tmp_path / "billing")


def test_app_level_packages_are_recognized_by_name(tmp_path: pathlib.Path) -> None:
    for name in sorted(APP_LEVEL_PACKAGES):
        _pkg(tmp_path, name, "x = 1\n")
    d = classify_root(tmp_path)
    assert d.contexts == () and d.unclassified == ()


def test_non_python_dirs_are_outside_the_contract(tmp_path: pathlib.Path) -> None:
    (tmp_path / "docs").mkdir()
    (tmp_path / "docs" / "guide.md").write_text("hello\n", encoding="utf-8")
    d = classify_root(tmp_path)
    assert d.contexts == () and d.unclassified == ()


def test_skip_dirs_are_skipped_even_with_python(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "build", "x = 1\n")
    d = classify_root(tmp_path)
    assert d.contexts == () and d.unclassified == ()


def test_namespace_package_cannot_hide(tmp_path: pathlib.Path) -> None:
    (tmp_path / "billing").mkdir()
    (tmp_path / "billing" / "service.py").write_text("x = 1\n", encoding="utf-8")
    d = classify_root(tmp_path)
    assert d.unclassified == ("billing",)


def test_broken_init_cannot_hide(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "billing", "def broken(:\n")
    d = classify_root(tmp_path)
    assert d.unclassified == ("billing",)


def test_extended_app_level_set(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "scripts", "x = 1\n")
    assert classify_root(tmp_path).unclassified == ("scripts",)
    extended = APP_LEVEL_PACKAGES | {"scripts"}
    assert classify_root(tmp_path, extended).unclassified == ()


def test_totality_errors_name_package_and_fix(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "reports", '"""the defect class"""\n')
    errors = totality_errors(tmp_path, classify_root(tmp_path))
    assert len(errors) == 1
    assert "reports" in errors[0]
    assert "Client" in errors[0] and "__init__.py" in errors[0]
    assert "--app-level" in errors[0]


def test_totality_errors_flag_empty_app_root(tmp_path: pathlib.Path) -> None:
    _pkg(tmp_path, "bootstrap", "x = 1\n")
    errors = totality_errors(tmp_path, classify_root(tmp_path))
    assert len(errors) == 1
    assert "no bounded contexts discovered" in errors[0]


def _good_tree(tmp_path: pathlib.Path) -> pathlib.Path:
    _pkg(tmp_path, "billing", "from billing.client import Client\n")
    (tmp_path / "billing" / "client.py").write_text(
        "from typing import Protocol\n\nclass Client(Protocol):\n    def ping(self) -> None: ...\n",
        encoding="utf-8",
    )
    _pkg(tmp_path, "bootstrap", "x = 1\n")
    return tmp_path


def test_cli_app_root_clean_tree_exits_zero(
    tmp_path: pathlib.Path, capsys: pytest.CaptureFixture[str]
) -> None:
    _good_tree(tmp_path)
    assert main(["--app-root", str(tmp_path)]) == 0
    assert capsys.readouterr().err == ""


def test_cli_app_root_unclassified_is_loud_error(
    tmp_path: pathlib.Path, capsys: pytest.CaptureFixture[str]
) -> None:
    _good_tree(tmp_path)
    _pkg(tmp_path, "reports", '"""forgot its Client"""\n')
    rc = main(["--app-root", str(tmp_path)])
    err = capsys.readouterr().err
    assert rc == 2
    assert "reports" in err and "Client" in err


def test_cli_app_level_extension_declares_plumbing(
    tmp_path: pathlib.Path, capsys: pytest.CaptureFixture[str]
) -> None:
    _good_tree(tmp_path)
    _pkg(tmp_path, "scripts", "x = 1\n")
    assert main(["--app-root", str(tmp_path)]) == 2
    capsys.readouterr()
    assert main(["--app-root", str(tmp_path), "--app-level", "scripts"]) == 0


def test_cli_app_level_requires_app_root(tmp_path: pathlib.Path) -> None:
    with pytest.raises(SystemExit) as exc:
        main(["--app-level", "scripts", str(tmp_path)])
    assert exc.value.code == 2


def test_cli_app_root_is_default_check_target(
    tmp_path: pathlib.Path, capsys: pytest.CaptureFixture[str]
) -> None:
    _good_tree(tmp_path)
    (tmp_path / "billing" / "money.py").write_text(
        "from dataclasses import dataclass\n\n\n@dataclass\nclass Money:\n"
        "    amount: int\n    currency: str\n",
        encoding="utf-8",
    )
    rc = main(["--app-root", str(tmp_path)])
    out = capsys.readouterr().out
    assert rc == 1
    assert "TB001" in out


def test_cli_explicit_paths_scope_checks_but_discovery_still_total(
    tmp_path: pathlib.Path, capsys: pytest.CaptureFixture[str]
) -> None:
    _good_tree(tmp_path)
    _pkg(tmp_path, "reports", '"""forgot its Client"""\n')
    rc = main(["--app-root", str(tmp_path), str(tmp_path / "billing")])
    err = capsys.readouterr().err
    assert rc == 2
    assert "reports" in err
