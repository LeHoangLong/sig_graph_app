'use strict'

const Client = require('pg').Client

module.exports.up = async function (next) {
  const client = new Client()
  await client.connect()
  try {
    await client.query('BEGIN')
    await client.query(`
      CREATE TABLE IF NOT EXISTS "peer" (
        id SERIAL PRIMARY KEY,
        owner_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE,
        alias TEXT NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "public_key" (
        id SERIAL PRIMARY KEY,
        value TEXT NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "peer_key" (
        id SERIAL PRIMARY KEY,
        owner_id INTEGER REFERENCES "peer"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL UNIQUE
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "user_key" (
        id SERIAL PRIMARY KEY,
        owner_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE,
        public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL UNIQUE,
        private_key TEXT NOT NULL,
        is_default BOOLEAN NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "material" (
        id TEXT PRIMARY KEY,
        public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE,
        name TEXT NOT NULL,
        quantity DECIMAL NOT NULL,
        unit TEXT NOT NULL,
        created_time TIMESTAMPTZ DEFAULT NOW()
      )
    `)
    await client.query('COMMIT')
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
      DROP TABLE IF EXISTS "material"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "user_key"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "peer_key"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "public_key"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "peer"
    `)
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
