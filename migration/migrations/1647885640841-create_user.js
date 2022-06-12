'use strict'

const Client = require('pg').Client

module.exports.up = async function (next) {
  const client = new Client()
  await client.connect()
  try {
    await client.query('BEGIN')
    await client.query('COMMIT')
    await client.query(`
      CREATE TABLE IF NOT EXISTS "user" (
        id SERIAL PRIMARY KEY,
        username TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        salt TEXT NOT NULL
      )
    `)

    // create test user
    await client.query(`
        INSERT INTO "user" (username, password_hash, salt) VALUES ('test', '', '') ON CONFLICT DO NOTHING
    `)
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}

module.exports.down = async function (next) {
  const client = new Client()
  await client.connect()
  try {
    await client.query('BEGIN')
    await client.query(`
      DROP TABLE IF EXISTS "user"
    `)
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
