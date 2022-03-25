'use strict'

const Client = require('pg').Client
const YAML = require('yaml')
const fs = require('fs')

module.exports.up = async function (next) {
  const client = new Client()
  await client.connect()
  const configFile = process.env.CONFIG_FILE
  const file = fs.readFileSync(configFile, 'utf8')
  const config = YAML.parse(file)
  try {
    await client.query('BEGIN')
    const userResponse = await client.query(`
      SELECT id FROM "user" WHERE username = 'test'
    `)

    const publicKeyFile = fs.readFileSync(config.user.public_key, 'utf-8')
    const privateKeyFile = fs.readFileSync(config.user.private_key, 'utf-8')

    const publicKeyResponse = await client.query(`
      INSERT INTO "public_key" (
        value
      ) VALUES (
        $1
      ) RETURNING id
    `, [publicKeyFile])

    await client.query(`
      INSERT INTO "user_key" (
        owner_id,
        public_key_id,
        private_key,
        is_default
      ) VALUES (
        $1,
        $2,
        $3,
        TRUE
      )
    `, [userResponse.rows[0].id, publicKeyResponse.rows[0].id, privateKeyFile])
      
    for (const [index, peer] of config.peers.entries()) {
      const peerResponse = await client.query(`
        INSERT INTO "peer" (
          owner_id,
          alias
        ) VALUES (
          $1,
          $2
        ) RETURNING id
      `, [userResponse.rows[0].id, peer.alias ?? ('peer_' + index)])

      for (const key of peer.public_keys) {
        const publicKeyFile = fs.readFileSync(key, 'utf-8')
        const publicKeyResponse = await client.query(`
          INSERT INTO "public_key" (
            value
          ) VALUES (
            $1
          ) RETURNING id
        `, [publicKeyFile])

        await client.query(`
          INSERT INTO "peer_key" (
            owner_id,
            public_key_id
          ) VALUES (
            $1,
            $2
          )
        `, [peerResponse.rows[0].id, publicKeyResponse.rows[0].id])
      }
    }

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
    await client.query('COMMIT')
  } catch (exception) {
    await client.query('ROLLBACK')
    throw exception
  } finally {
    await client.end()
  }
}
