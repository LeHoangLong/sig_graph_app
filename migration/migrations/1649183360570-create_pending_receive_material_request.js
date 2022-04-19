'use strict'

const Client = require('pg').Client

module.exports.up = async function (next) {
  const client = new Client()
  await client.connect()
  try {
    await client.query('BEGIN')

    await client.query(`
      CREATE TABLE IF NOT EXISTS "pending_receive_material_request" (
        id SERIAL PRIMARY KEY,
        main_node_id TEXT NOT NULL,
        transfer_time TIMESTAMPTZ NOT NULL
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
      CREATE TABLE IF NOT EXISTS "node_from_peer" (
        id SERIAL PRIMARY KEY,
        node_id TEXT NOT NULL,
        is_finalized BOOLEAN NOT NULL,
        created_time TIMESTAMPTZ NOT NULL,
        signature TEXT NOT NULL,
        public_key TEXT NOT NULL,
        pending_receive_material_request_id INTEGER REFERENCES "pending_receive_material_request"(id) NOT NULL,
        type TEXT REFERENCES "node_type"(type) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "material_from_peer" (
        id INTEGER REFERENCES "node_from_peer"(id) PRIMARY KEY,
        name TEXT NOT NULL,
        quantity TEXT NOT NULL,
        unit TEXT NOT NULL
      )
    `)

    await client.query(`
        CREATE TABLE IF NOT EXISTS "parent_hashed_id"(
          id SERIAL PRIMARY KEY,
          owner_node_id INTEGER REFERENCES "node_from_peer"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
          hashed_id TEXT NOT NULL,
          parent_node_id INTEGER REFERENCES "node_from_peer"(id) ON DELETE CASCADE ON UPDATE CASCADE
        )
    `)

    await client.query(`
        CREATE TABLE IF NOT EXISTS "child_hashed_id"(
          id SERIAL PRIMARY KEY,
          owner_node_id INTEGER REFERENCES "node_from_peer"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
          hashed_id TEXT NOT NULL,
          child_node_id INTEGER REFERENCES "node_from_peer"(id) ON DELETE CASCADE ON UPDATE CASCADE
        )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "signature_option" (
        signature TEXT NOT NULL,
        new_node_id TEXT NOT NULL,
        pending_request_id INTEGER REFERENCES "pending_receive_material_request"(id) ON DELETE CASCADE ON UPDATE CASCADE
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
    await client.query(`DROP TABLE IF EXISTS "signature_option"`)
    await client.query(`DROP TABLE IF EXISTS "child_hashed_id"`)
    await client.query(`DROP TABLE IF EXISTS "parent_hashed_id"`)
    await client.query(`DROP TABLE IF EXISTS "material_from_peer"`)
    await client.query(`DROP TABLE IF EXISTS "node_from_peer"`)
    await client.query(`DROP TABLE IF EXISTS "node_type"`)
    await client.query(`DROP TABLE IF EXISTS "pending_receive_material_request"`)
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
