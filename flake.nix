{
  description = "Dev environment for simpleauth";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";

    # 1.19.1 release
    go_dep.url = "github:NixOS/nixpkgs/8ba120420fbdd9bd35b3a5366fa0206d8c99ade3";

    # 0.9.0 release
    shellcheck_dep.url = "github:NixOS/nixpkgs/8b5ab8341e33322e5b66fb46ce23d724050f6606";
  };

  outputs = { self, flake-utils, go_dep, shellcheck_dep }@inputs :
    flake-utils.lib.eachDefaultSystem (system:
    let
      go_dep = inputs.go_dep.legacyPackages.${system};
      shellcheck_dep = inputs.shellcheck_dep.legacyPackages.${system};
    in
    {
      devShells.default = go_dep.mkShell.override { stdenv = go_dep.pkgsStatic.stdenv; } {
        packages = [
          go_dep.gotools
          go_dep.gopls
          go_dep.go-outline
          go_dep.gocode
          go_dep.gopkgs
          go_dep.gocode-gomod
          go_dep.godef
          go_dep.golint
          go_dep.go_1_19
          shellcheck_dep.shellcheck
        ];

        shellHook = ''
          GOROOT="$(dirname $(dirname $(which go)))/share/go"
          export GOROOT

          echo "shellcheck" "$(shellcheck --version | grep '^version:')"
          go version
        '';
      };
    });
}
