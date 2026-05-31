# Websocket

This is the official websocket used to transmit and listen to audio being broadcasted.

## Websocket Endpoint

The websocket endpoint is ws://localhost:5000/v1/ws. But when you connect, you have to retreive a session key.

## How To Get a Session Key

You will automatically receive a session key when connecting to websocket.

To fully connect you can use wscat:

```bash
wscat -c ws://localhost:5000/v1/ws
```