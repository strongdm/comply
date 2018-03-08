package site

const cssACE = `
`

const indexACE = `= doctype html
html lang=en
  head
    meta charset=utf-8
    title Comply
    link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/bulma/0.6.2/css/bulma.min.css"
    meta name="viewport" content="width=device-width, initial-scale=1"
  body
    section.hero.is-primary
      .hero-body
        .container
          h1.title {{.Msg}}
          p.subtitle My first
`
