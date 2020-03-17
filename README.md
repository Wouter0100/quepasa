# QuePasa

> A (micro) web-application to make web-based [WhatsApp][0] bots easy to write.

Built on the library [go-whatsapp][1] by [Rhymen][2].

**Implemented features:**

* Register a number with WhatsApp
* Verify a number with a QR code
* Persistence of account data and keys
* Exposes HTTP endpoints for:
  * sending messages
  * receiving messages

**WARNING: This application has not been audited. It should not be regarded as
secure, use at your own risk.**

**This is a third-party effort, and is NOT in any affiliated with [WhatsApp][0].**

### Why?

When you need to communicate over WhatsApp from a different service, for example,
[a help desk](http://zammad.org/) or other web-app, QuePasa provides a simple HTTP
API to do so.

QuePasa stores keys and WhatsApp account data in a postgres database. It does
not come with HTTPS out of the box. Your QuePasa API tokens essentially give
full access to your WhatsApp account (to the extent that QuePasa has
implemented WhatsApp features). Use with caution.

## Usage

### Prerequisites

For local development
* docker
* golang
* postgresql

### Run using Docker

* Add info about database migrations

```bash

make docker_build
# edit docker-compose.yml.sample to your hearts content
docker-compose up
```

## HTTP API

1. Use the `Accept: application/json` header
2. `TOKEN` should be treated like a password.

### Get bot info

A simple method for testing your bot's auth token. Requires no parameters. Returns basic information about the bot.

**request**
```
GET /bot/<TOKEN>/
```

***response***

```json
{
    "id": "129f1757-e706-452e-aa1c-4994a95e1092",
    "number": "+15555555552",
    "user_id": "845ae4d0-f2c3-5342-91a2-5b45cb8db57c",
    "token": "8129c0b4-0b96-4486-84fc-c3dd7b03f846",
    "is_verified": true,
    "created_at": "2018-11-02T11:36:24.273Z",
    "updated_at": "2018-11-02T11:36:24.273Z"
}

```

### Sending

**request**
```
POST /bot/<TOKEN>/send

{
  "recipient": "+15555555552",
  "messsage": "Hello World!"
}
```

**response**
```json
{
  "result": {
    "recipient": "+15555555551",
    "source": "+15555555552",
    "status": "sent",
    "timestamp": "1543420505142"
  }
}
```

### Receive

The "timestamp" query parameter is optional. A maximum of 40 messages per conversation will be returned.

**request**
```
GET /bot/<TOKEN>/receive?timestamp=1541265073783
```

**response**
```json
{
  "messages": [
    {
      "source": "+15555555551",
      "timestamp": "1541265073894",
      "message": {
        "body": "Hello World!",
        "profileKey": "XXTXQ=="
      }
    }
  ],
  "bot": {
    "id": "129f1757-e706-452e-aa1c-4994a95e1092",
    "number": "+15555555552",
    "user_id": "845ae4d0-f2c3-5342-91a2-5b45cb8db57c",
    "token": "8129c0b4-0b96-4486-84fc-c3dd7b03f846",
    "is_verified": true,
    "created_at": "2018-11-02T11:36:24.273Z",
    "updated_at": "2018-11-02T11:36:24.273Z"
  }
}
```

## License

[![License GNU AGPL v3.0](https://img.shields.io/badge/License-AGPL%203.0-lightgrey.svg)](https://gitlab.com/digiresilience/link/quepasa/blob/master/LICENSE.md)

QuePasa is a free software project licensed under the GNU Affero General Public License v3.0 (GNU AGPLv3) by [The Center for Digital Resilience](https://digiresilience.org).

[0]: https://whatsapp.com
[1]: https://github.com/Rhymen/go-whatsapp
[2]: https://github.com/Rhymen
