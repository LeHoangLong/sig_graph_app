'use strict'

const Client = require('pg').Client

module.exports.up = async function (next) {
  const client = new Client()
  await client.connect()
  try {
    await client.query('BEGIN')

    await client.query(`
      CREATE TABLE IF NOT EXISTS "receive_material_request_status" (
        id SERIAL PRIMARY KEY,
        description TEXT UNIQUE NOT NULL
      )
    `)

    await client.query(`
      INSERT INTO "receive_material_request_status" (
        id, description
      ) VALUES (0, 'pending'), (1, 'accepted'), (2, 'rejected')
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "receive_material_request" (
        id SERIAL PRIMARY KEY,
        main_node_id INTEGER UNIQUE REFERENCES "node"(id),
        transfer_time TIMESTAMPTZ NOT NULL,
        status_id INT REFERENCES "receive_material_request_status"(id) ON DELETE NO ACTION ON UPDATE CASCADE NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "outbound_receive_material_request" (
        request_id INTEGER PRIMARY KEY REFERENCES "receive_material_request"(id),
        owner_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        owner_public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE,
        recipient_peer_id INTEGER REFERENCES "peer"(id)
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "inbound_receive_material_request" (
        request_id INTEGER PRIMARY KEY REFERENCES "receive_material_request"(id),
        owner_id INTEGER REFERENCES "user"(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
        sender_public_key_id INTEGER REFERENCES "public_key"(id) ON DELETE CASCADE ON UPDATE CASCADE
      )
    `)
    
    await client.query(`
      CREATE TABLE IF NOT EXISTS "receive_material_request_acknowledgement" (
        request_id INT REFERENCES "receive_material_request" (id) ON UPDATE CASCADE ON DELETE CASCADE NOT NULL UNIQUE,
        response_id TEXT NOT NULL UNIQUE
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "node_from_peer" (
        node_id INTEGER PRIMARY KEY REFERENCES "node"(id),
        receive_material_request_id INTEGER REFERENCES "receive_material_request"(id) NOT NULL
      )
    `)

    await client.query(`
      CREATE TABLE IF NOT EXISTS "signature_option" (
        signature BYTEA NOT NULL,
        new_node_id TEXT NOT NULL,
        request_id INTEGER REFERENCES "receive_material_request"(id) ON DELETE CASCADE ON UPDATE CASCADE
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
    await client.query(`DROP TABLE IF EXISTS "node_from_peer"`)
    await client.query(`DROP TABLE IF EXISTS "receive_material_request_response"`)
    await client.query(`DROP TABLE IF EXISTS "outbound_receive_material_request"`)
    await client.query(`DROP TABLE IF EXISTS "inbound_receive_material_request"`)
    await client.query(`DROP TABLE IF EXISTS "receive_material_request"`)
    await client.query(`DROP TABLE IF EXISTS "receive_material_request_status"`)
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
