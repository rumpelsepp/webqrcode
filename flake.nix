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
        vendorHash = "sha256-rUQ/iXG/lnvi01NXBOsu3CoG3YrYUFwPZDadTFThH3g=";
      };
      formatter.x86_64-linux = pkgs.nixpkgs-fmt;
    };
}
