const express = require('express')
const dotenv = require('dotenv')

dotenv.config()

const {
  EXPRESS_PORT = 80,
  TARGET_URL,
  WARNING,
} = process.env


const app = express()

app.set('view engine', 'ejs');

app.all('/*', (req, res) => {
  if (!TARGET_URL) throw `TARGET_URL not set`
  if (WARNING) res.render('index', { target: TARGET_URL });
  else res.redirect(TARGET_URL)
})

app.listen(EXPRESS_PORT, () => {
  console.log(`[Express] Listening on port ${EXPRESS_PORT}`)
})
