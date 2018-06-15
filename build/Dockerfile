FROM haskell:latest

# based on implementation by James Gregory <james@jagregory.com>
MAINTAINER Comply <comply@strongdm.com>

# install latex packages
RUN apt-get update -y \
  && apt-get install -y -o Acquire::Retries=10 --no-install-recommends \
    texlive-latex-base \
    texlive-xetex \
    texlive-fonts-recommended \
    latex-xcolor \
    texlive-latex-extra \
    fontconfig \
    unzip \
    lmodern

# will ease up the update process
# updating this env variable will trigger the automatic build of the Docker image
ENV PANDOC_VERSION "2.2.1"

# install pandoc
RUN cabal update && cabal install pandoc-${PANDOC_VERSION}

WORKDIR /source

ENTRYPOINT ["/root/.cabal/bin/pandoc"]

CMD ["--help"]
