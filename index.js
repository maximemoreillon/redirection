const express = require('express')
const dotenv = require('dotenv')

dotenv.config()

const app = express()
const port = process.env.EXPRESS_PORT || 80

app.get('/*', (req, res) => {
  const target_url = process.env.TARGET_URL
  if(!target_url) res.status(500).send(`TARGET_URL not set`)
  res.redirect(target_url)
})

app.listen(port, () => {
  console.log(`Redirect service on :${port}`)
})
