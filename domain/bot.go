package domain

import "time"

const BotNickname = "bot"
const BotPlayAgainDelay = 3 * time.Second
const BotMoveDelay = 2 * time.Second

var BotNicknames = []Nickname{BotNickname, "b0t"}
