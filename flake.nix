{
  description = "naga.srin.cc — Zola site + now-playing Go service";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { self, nixpkgs }:
    let
      # Keep in sync with the zola tag in ./Dockerfile. The site's config.toml
      # tracks whichever schema that version expects, so a mismatch between the
      # devshell and the build image means `zola serve` and CI disagree.
      systems = [ "aarch64-darwin" "x86_64-darwin" "aarch64-linux" "x86_64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = [ pkgs.zola pkgs.go ];

          shellHook = ''
            echo "zola $(zola --version | cut -d' ' -f2) · $(go version | cut -d' ' -f3)"
            echo "  zola serve            dev server"
            echo "  zola build            render to public/"
            echo "  cd now-playing && go test ./...   API tests"
          '';
        };
      });
    };
}
