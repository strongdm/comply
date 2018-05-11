class Comply < Formula
    desc ""
    homepage ""
    url "file:///Users/jmccarthy/Downloads/comply-1.1.3.tgz"
    sha256 "01f9920e5e9fbd29386e4a4131fac78c002730e49c3f870a0ee2b958c3ce75f6"

    depends_on "go" => :build

    def install
        ENV["GOPATH"] = buildpath
        ENV.prepend_create_path "PATH", buildpath/"bin"
        (buildpath/"src/github.com/strongdm/comply").install buildpath.children
        cd "src/github.com/strongdm/comply" do
            system "make", "brew"
            bin.install "bin/comply"
        end
    end

    test do
        system "#{bin}/sdm", "logout"
    end
end
