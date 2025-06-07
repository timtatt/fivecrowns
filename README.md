# fivecrowns

This is a repo to play around with building engines to play the rummy-style game - [five crowns](https://en.wikipedia.org/wiki/Five_Crowns_(card_game))


## Arena

The arena is a basic html page to test out bots

To run the arena
```
go run .
open http://localhost:3000/arena
```

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
    "action": "draw", // draw, discard, score
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

Draw is the first step of a turn where the bot choose whether to take the top card of the deck OR take the top card from discard pile
```js
{
    "action": "draw", // draw, discard, score
    "stack": "deck", // deck, discard
}
```


Discard is the second part of a turn where the bot will arrange the cards into sequences and discard a card to finish their turn
```js
{
    "action": "discard", // draw, discard, score
    "card": "10-R", // which card has been discarded
    "flop": false, // indicates that the sequences are sufficient to 'flop' (i.e. reveal hand and start final turns for other players)
    "sequences": [
        ["9-R", "10-R", "11-R"],
        ["9-Y"]
    ], // list of valid sequences computed by the bot
}
```


Score is an action which can be used to test the engine. It is the same as 'discard' but doesn't try to discard a card.
```js
{
    "action": "score", // draw, discard, score
    "flop": false, // indicates that the sequences are sufficient to 'flop' (i.e. reveal hand and start final turns for other players)
    "sequences": [
        ["9-R", "10-R", "11-R"],
        ["9-Y"]
    ], // list of valid sequences computed by the bot
}
```


## Bots

### Smooth Brain Bot
- Cannot build sequences
- Randomly picked a stack to draw from
- Randomly picks a card to discard

### Grug Brain Bot
- Can build sequences
- Priotises highest scoring sequences
- Uses wilds to rid of most # of cards


### Big Brain Bot
- Builds sequences
- Can split sequences to maximise card usage
- Prefers sets to runs, especially in rounds 4, 5, 7, 8
- Uses wilds to maximise used cards

### Galaxy Brain Bot
- Reads the discard pile and uses probability to decide on best actions
- Can go into damage control mode when the turn count is high depending on number of players and curent round i.e. when there have been many turns, higher probability that someone will flop soon
- Observes which cards opponent to left is picking up; will avoid handing them cards
- Uses AI to trash talk opponents based on the cards they have discarded
