{ pkgs ? import <nixpkgs> {} }:
let
  rootDir = builtins.toString ./.;
in pkgs.mkShell {
  buildInputs = with pkgs; [
    terraform_1
    git
    makeWrapper
    go_1_21
  ];
}
