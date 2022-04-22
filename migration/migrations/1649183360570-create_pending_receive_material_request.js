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
      CREATE TABLE IF NOT EXISTS "node_from_peer" (
        node_id INTEGER PRIMARY KEY REFERENCES "node"(id),
        pending_receive_material_request_id INTEGER REFERENCES "pending_receive_material_request"(id) NOT NULL
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
    await client.query(`DROP TABLE IF EXISTS "material_from_peer"`)
    await client.query(`DROP TABLE IF EXISTS "node_from_peer"`)
    await client.query(`DROP TABLE IF EXISTS "pending_receive_material_request"`)
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
