export const SET_PLAYER_NAME_REQUEST = "set_player_name"
export const JOIN_ROOM_REQUEST = "join_room"
export const PREPARE_REQUEST = "prepare"
export const SHOW_MASTER = "show_master"
export const SHOW_MASTER_DONE = "show_master_done"
export const KOU_CARDS = "kou_cards"
export const PLAY_CARDS = "play_cards"

export const ROOM_LIST_RESPONSE = "room_list"
export const EXISTS_PLAYERS_RESPONSE = "exists_players"
export const SHOW_MASTER_RESPONSE = "show_master_result"
export const DEAL_POKER = "deal_poker"
export const MATCH_BEGIN = "match_begin"
export const REVEIVE_HOLE_CARDS = "deal_hole_cards"
export const PLAY_TURN = "play_trun"
export const SHOW_PLAY_CARDS = "show_play_cards"
export const INCREASE_SCORES = "increase_scores"
export const ROUND_END = "round_end"
export const BIGGEST_POSITION = "biggest_position"
export const GAME_RESULT = "game_result"

export const FULL_POKER_NUM = 25

export interface Player {
    name: string
    prepare: boolean
}

export interface GameResult {
    ourRound: string
    otherRound: string
    bankerPostion: number
}

export interface Poker {
    color: string
    number: string
}

export interface Cards {
    spadeCards: Poker[]
    heartCards: Poker[]
    clubCards: Poker[]
    dianmondCards: Poker[]
    jokers: Poker[]
    playNumberCards: Poker[]
    cardNum: number
}

export interface ShowMasterRequest {
    color: string
    isSelfProtect: boolean,
    isOppose: boolean,
}

export interface ShowPlayCardsResponse {
    showIdx: number
    cards: Poker[]
}

export interface ShowMasterResponse {
    color: string,
    isProtect: boolean,
    isSelfShowMaster: boolean
    showMasterPosition: number
}

const positionMap = new Map([
    ['A', 0],
    ["2", 1],
    ['3', 2],
    ["4", 3],
    ["5", 4],
    ["6", 5],
    ['7', 6],
    ["8", 7],
    ['9', 8],
    ["10", 9],
    ["J", 10],
    ["Q", 11],
    ["K", 12],
])

export const CardValueMap = new Map([
    ["2", 1],
    ['3', 2],
    ["4", 3],
    ["5", 4],
    ["6", 5],
    ['7', 6],
    ["8", 7],
    ['9', 8],
    ["10", 9],
    ["J", 10],
    ["Q", 11],
    ["K", 12],
    ['A', 13],
])

export const SPADE = "spade"
export const HEART = "heart"
export const CLUB = "club"
export const DIANMOND = "dianmond"

export const showColorIdxMap = new Map([
    [HEART, 0],
    [SPADE, 1],
    [CLUB, 2],
    [DIANMOND, 3],
])

export const getPokerPosition = (poker: Poker): number => {
    if (poker.color === 'red') {
        return 52
    }
    if (poker.color === 'black') {
        return 53
    }
    let position = positionMap.get(poker.number)
    if (position === undefined) {
        return -1
    }
    if (poker.color === HEART) {
        return position
    }
    if (poker.color === DIANMOND) {
        return 13 + position
    }
    if (poker.color === SPADE) {
        return 26 + position
    }
    return 39 + position
}