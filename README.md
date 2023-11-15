# FlagGuessr

Flag Guessr is a minigame Discord bot about guessing countries based on their flag and hints related to the country.

![](https://i.imgur.com/PIjix2L.png)

## Features

- Pool of 250 countries and territories
- [Hints](#hints), if you don't recognize the flag
- [Options that enhance your experience](#options)
- Country streaks

### Hints

1. Driving side
2. Population
3. Capitals
4. Top Level Domains

![](https://i.imgur.com/mz030o1.png)

### Options

1. Difficulty (Normal & Hard) - on Hard difficulty, guessing a country incorrectly will result in immediately losing the streak
2. Minimum population - [predefined set of populations](https://i.imgur.com/tAwsq6A.png) to filter countries
3. Hide - whether the game embeds should only be visible to the user

![](https://i.imgur.com/9owCmJE.png)

## Ideas for expansion

1. Leaderboard system

## Ready to play?

Invite Flag Guessr using [this link](https://discord.com/oauth2/authorize?client_id=1007647563790417960&scope=bot).

## Have a question, suggestion or a concern?

Join the [support server](https://discord.gg/bgku8tPVqN).

## Technical details

Flag Guessr is completely stateless - information about games is stored in buttons/modals as JSON. This means that there is no database and users don't lose their progress after bot restarts (unless they delete the embeds).

### [Button data](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L41-#L50)

```json
{
   "u":"244563286257827840",
   "d":0,
   "m":10000000,
   "i":22,
   "e":false,
   "s":1,
   "a":0,
   "h":0
}
```

### [Modal data](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L69-#L75)

```json
{
   "d":0,
   "m":10000000,
   "i":22,
   "e":false,
   "s":1
}
```

#### Breakdown


[`u`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L42): **u**ser ID, used for checking who can interact with the game embed

[`d`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L43): [**d**ifficulty (0 = Normal, 1 = Hard)](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L52-#L67)

[`m`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L44): **m**inimum population

[`i`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L45): **i**ndex of the country in the slice, used for getting details

[`e`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L46): whether the game embeds should be **e**phemeral

[`s`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L47): current **s**treak

[`a`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L48): [**a**ction the button should execute](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L9-#L17)

[`h`](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L49): [current **h**int](https://github.com/caneleex/FlagGuessr/blob/main/util/types.go#L19-#L27)
