{ pkgs ? import <nixpkgs> {} }:
let
  rootDir = builtins.toString ./.;
in pkgs.mkShell {
  buildInputs = with pkgs; [
    terraform_1
    git
    makeWrapper
    go_1_20
    pre-commit
    golangci-lint
  ];

  shellHook = ''
      #-------------------------------------------------------
      # pre-commit for git
      #-------------------------------------------------------
      [ -f "${rootDir}/.git/hooks/pre-commit" ] || pre-commit install
  '';
}
