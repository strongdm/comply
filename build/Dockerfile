FROM pandoc/ubuntu

RUN apt-get update -y \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y -o Acquire::Retries=10 --no-install-recommends \
  texlive-latex-base \
  texlive-latex-extra \
  texlive-plain-generic \
  lmodern

WORKDIR /source
