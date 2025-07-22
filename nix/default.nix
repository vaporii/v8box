{ pkgs ? (import <nixpkgs> {}), ... }:

pkgs.buildGoModule {
  pname = "v8box";
  version = "0.1.0";

  src = ../.;

  vendorHash = "sha256-yRgT6vp/EGWaglN7LKgKkaGP/hFzC86EncyZMQ6iRcA=";
}
