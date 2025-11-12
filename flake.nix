{
  description = "LocalBox - general purpose sandbox for running untrusted code";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      utils,
    }:
    utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
        packages = {
          localbox = pkgs.buildGoModule {
            env.CGO_ENABLED = 0;
            name = "localbox";
            src = self;
            goSum = ./go.sum;
            vendorHash = "sha256-ajutaBjeEeMtmYvyA6uORG+D9NjnDi2sPM/13fQVLWA=";
            checkFlags = [ "-skip" ];
          };
          isolate = pkgs.isolate;
          default = packages.localbox;
        };
      }
    );
}