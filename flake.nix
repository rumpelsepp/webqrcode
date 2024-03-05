{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs";
  };

  outputs = { self, nixpkgs }:
    with import nixpkgs { system = "x86_64-linux"; };
    let pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in {
      devShell.x86_64-linux = pkgs.mkShell {
        buildInputs = with pkgs; [
          go gopls gotools
        ];
      };
      packages.x86_64-linux.default =
      pkgs.buildGoModule {
        pname = "webqrcode";
        src = self;
        version = self.lastModifiedDate;
        buildInputs = [pkgs.go];
        vendorHash = "sha256-6UyljuhLPfCCkYSLmHp40ltDMAHnD5ZYU7CXSfYH6CQ=";
      };
      formatter.x86_64-linux = pkgs.nixpkgs-fmt;
    };
}
