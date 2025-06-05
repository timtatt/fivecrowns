# fivecrowns

This is a repo to play around with building engines to play the rummy-style game - [five crowns](https://en.wikipedia.org/wiki/Five_Crowns_(card_game))


## Arena

The arena is a basic html page to test out bots

## Spec

The interface for a five crowns bot is one of:
- HTTP endpoint 
- gRPC - TBD
- WebSocket - TBD


### Requests

A turn is made up of 2 requests
1. Ask the bot whether it wants to draw from the deck or the discard pile
2. Ask the bot to discard a card from its hand

```js
{
    "action": "draw", // draw, discard
    "playerCount": 4,
    "round": 3, // 3-13, also indicates the number of cards and which one is wild
    "lastTurn": true, // when a player finishes, every other player gets 1 more turn. this indicates if it is the last turn
    "hand": ["10-R"], // list of cards in the players hand
    "newestCard": "", // if action = discard, indicates which card the player has drawn; can come from the deck or discard pile
    "discard": [""], // list of cards in the discard pile. top-most card is at index 0
}
```

### Responses

Below are examples of valid responses to each action

```js
{
    "action": "draw", // draw, discard
    "stack": "deck", // deck, discard
}
```


```js
{
    "action": "discard", // draw, discard
    "card": "10-R", // which card has been discarded
    "sequences": [
        ["9-R", "10-R", "11-R"],
        ["9-Y"]
    ], // list of valid sequences computed by the bot
}
```

