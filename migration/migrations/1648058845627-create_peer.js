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
        CREATE TABLE IF NOT EXISTS "node_type"(
          type TEXT PRIMARY KEY
        )
    `)

    await client.query(`
          INSERT INTO "node_type" (
            type
          ) VALUES (
            'material'
          )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "node" (
        id SERIAL PRIMARY KEY,
        node_id TEXT NOT NULL,
        public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE,
        is_finalized BOOLEAN NOT NULL,
        previous_node_hashed_ids TEXT[] NOT NULL,
        next_node_hashed_ids TEXT[] NOT NULL,
        created_time TIMESTAMPTZ NOT NULL,
        signature BYTEA NOT NULL,
        type TEXT REFERENCES "node_type"(type) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL
      )
    `)

    await client.query(`
        CREATE TABLE IF NOT EXISTS "node_edge"(
          id SERIAL PRIMARY KEY,
          owner_node_id INTEGER REFERENCES "node"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
          referenced_node_id INTEGER REFERENCES "node"(id) ON DELETE CASCADE ON UPDATE CASCADE
        )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "material" (
        node_id INTEGER PRIMARY KEY REFERENCES "node"(id),
        name TEXT NOT NULL,
        quantity DECIMAL NOT NULL,
        unit TEXT NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "supported_peer_protocol" (
        id SERIAL PRIMARY KEY,
        protocol TEXT NOT NULL UNIQUE,
        major_version INTEGER NOT NULL CHECK (major_version >= 0),
        minor_version INTEGER NOT NULL CHECK (minor_version >= 0)
      )
    `)
    
    await client.query(`
      CREATE UNIQUE INDEX protocol_index ON "supported_peer_protocol" (
        protocol,
        major_version,
        minor_version
      )
    `)

    await client.query(`
      INSERT INTO "supported_peer_protocol" (
        protocol,
        major_version,
        minor_version
      ) VALUES (
        'graphql',
        0,
        0
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "peer_endpoint" (
          id SERIAL PRIMARY KEY,
          peer_id INTEGER REFERENCES "peer"(id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL,
          url TEXT NOT NULL CHECK(LENGTH(url) > 0),
          protocol_id INTEGER REFERENCES "supported_peer_protocol"(id) ON UPDATE CASCADE ON DELETE RESTRICT NOT NULL
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
      DROP TABLE IF EXISTS "peer_endpoint"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "supported_peer_protocol"
    `)
    await client.query(`
      DROP TABLE IF EXISTS "material"
    `)
    await client.query(`DROP TABLE IF EXISTS "node_edge"`)
    await client.query(`
      DROP TABLE IF EXISTS "node"
    `)
    await client.query(`DROP TABLE IF EXISTS "node_type"`)
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
