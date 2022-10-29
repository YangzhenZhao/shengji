// @ts-nocheck
import Phaser from "phaser";
import bgimage from "../assets/bg2.jpg"
import pokerImage from "../assets/poker.png"
import beginGame from "../assets/btn/begin.png"
import prepareOk from "../assets/btn/prepare_ok.png"
import flowerImages from "../assets/flower_color.jpeg"
import jokerImage from "../assets/joker.webp"
import zhuangImage from "../assets/zhuang.jpg"
import miniStarImg from "../assets/ministar.png"
import goodImg from "../assets/good.jpeg"
import kouGreyImage from "../assets/kou_hei.gif"
import kouLightImage from "../assets/kou_light.gif"
import playCardsGreyImg from "../assets/chu_pai_grey.gif"
import playCardsGreenImg from "../assets/chu_pai_green.gif"
import { nanoid } from 'nanoid'
import { Match } from './match';
import { GameDetail } from './game_detail';
import { 
    Player, Poker, Cards, GameResult, SPADE, HEART, CLUB, DIANMOND, getPokerPosition,
    SET_PLAYER_NAME_REQUEST, JOIN_ROOM_REQUEST, PREPARE_REQUEST, 
    ROOM_LIST_RESPONSE, EXISTS_PLAYERS_RESPONSE, DEAL_POKER, MATCH_BEGIN,
    SHOW_MASTER_DONE, FULL_POKER_NUM, SHOW_MASTER_RESPONSE, showColorIdxMap, REVEIVE_HOLE_CARDS,
    KOU_CARDS, PLAY_TURN, PLAY_CARDS, SHOW_PLAY_CARDS, ShowPlayCardsResponse,
    INCREASE_SCORES, ROUND_END, CardValueMap, BIGGEST_POSITION, GAME_RESULT,
} from './dto'
import dayjs from 'dayjs';

const screenWidth = document.documentElement.clientWidth;
const screenHeight =  document.documentElement.clientHeight;

// 屏幕左上角为坐标原点 (0, 0), 横轴为 x, 纵轴为 y
const playerTextPositions = [
    {x: 50, y: screenHeight - 36},
    {x: screenWidth * 0.6, y: 30},
    {x: 5, y: screenHeight / 2.4},
    {x: screenWidth - 108, y: screenHeight / 2.4}
]

const zhuangPositions = [
    {x: 110, y: screenHeight - 45},
    {x: screenWidth * 0.66, y: 21},
    {x: 65, y: screenHeight / 2.45},
    {x: screenWidth - 48, y: screenHeight / 2.45}
]

const prepareImagePositions = [
    {x: screenWidth * 0.4, y: screenHeight - 330},
    {x: screenWidth * 0.6, y: 30},
    {x: 5, y: screenHeight / 2.4},
    {x: screenWidth - 108, y: screenHeight / 2.4}
]

const prepareOkPositions = [
    {x: screenWidth * 0.4, y: screenHeight - 280},
    {x: screenWidth * 0.57, y: 48},
    {x: 72, y: screenHeight / 2.65},
    {x: screenWidth - 230, y: screenHeight / 2.65}
]
let pokerPositions = []
var x = 200;
for (let i = 0; i < 25; i++) {
    pokerPositions.push({
        x: x,
        y: screenHeight - 160,
    })
    x += 26;
}
let showPokerBeginPositions = [
    {x: 300, y: screenHeight - 320},
    {x: screenWidth * 0.27, y: 48},
    {x: 72, y: screenHeight / 2.65},
    {x: screenWidth - 300, y: screenHeight / 2.65}
]
let showPokerPositions = [
    [],
    [],
    [],
    [],
]
for (let i = 0; i < 25; i++) {
    for (let j = 0; j < 4; j++) {
        showPokerPositions[j].push({
            x: showPokerBeginPositions[j].x + 26 * i,
            y: showPokerBeginPositions[j].y
        })
    }
}
let buckleCardPositions = []
x = 300
for (let i = 0; i < 17; i++) {
    buckleCardPositions.push({
        x: x,
        y: screenHeight - 280,
    })
    x += 30;
}
x = 300
for (let i = 0; i < 16; i++) {
    buckleCardPositions.push({
        x: x,
        y: screenHeight - 160,
    })
    x += 30;
}
const playCardsImgX = screenWidth * 0.5
const playCardsImgY = screenHeight - 220

export class GameScene extends Phaser.Scene {
    public playerName
    public websocket: WebSocket | null
    public players: Player[]
    public prepareBtn: Phaser.GameObjects.Image
    public zhuangColorImg: Phaser.GameObjects.Image
    public kouImg: Phaser.GameObjects.Image
    public playCardsImg: Phaser.GameObjects.Image
    public zhuangImg: Phaser.GameObjects.Image
    public prepareOkImg: Phaser.GameObjects.Image[]
    public playersCards: Cards
    public match: Match
    public gameDetail: GameDetail
    public pokerImages: Phaser.GameObjects.Image[]
    public showPlayCardsImgs: Phaser.GameObjects.Image[][]
    public showMasterImages: Phaser.GameObjects.Image[]
    public buckleCards: Poker[]
    public playCards: Poker[]
    public isFirstGame: boolean
    constructor(playerName: string) {
        super("GameScene")
        this.playerName = playerName
        this.biggestPosition = -1
        this.playersTexts = []
        this.players = [{name: playerName, prepare: false}, null, null, null]
        this.websocket = null
        this.getPlayerMsgCnt = 0
        this.prepareOkImg = [null, null, null, null]
        this.playersCards = {
            spadeCards: [],
            heartCards: [],
            clubCards: [],
            dianmondCards: [],
            jokers: [],
            playNumberCards: [],
            cardNum: 0,
        }
        this.match = new Match(this)
        this.gameDetail = new GameDetail(this)
        this.pokerImages = []
        this.showMasterImages = []
        this.selectBuckleNum = 0
        this.buckleCards = []
        this.playCards = []
        this.showPlayCardsImgs = [[], [], [], []]
        this.heartCheckTimeout = 10000  //10s
        this.heartCheckServerTimeout = 30000 // 30s
        this.heartCheckTimeoutObj = null
        this.heartCheckServerTimeoutObj = null
        this.isFirstGame = true
    }

    preload () {
        console.log(dayjs().format())
        this.load.image("bg2", bgimage)
        this.load.image("beginGame", beginGame)
        this.load.image("prepareOk", prepareOk)
        this.load.image("jokerImage", jokerImage)
        this.load.image("zhuangImage", zhuangImage)
        this.load.image("miniStarImage", miniStarImg)
        this.load.image("goodImage", goodImg)
        this.load.image("kouGreyImage", kouGreyImage)
        this.load.image("kouLightImage", kouLightImage)
        this.load.image("playCardsGreyImg", playCardsGreyImg)
        this.load.image("playCardsGreenImg", playCardsGreenImg)
        this.load.spritesheet('poker', pokerImage, {
            frameWidth: 90,
            frameHeight: 120
        });
        this.load.spritesheet('flowerImages', flowerImages, {
            frameWidth: 250,
            frameHeight: 250
        });
        
    }

    create() {
        this.add.image(0, 0, 'bg2').setOrigin(0).setDisplaySize(screenWidth, screenHeight);

        // this.add.image(500, 500, 'miniStarImage').setOrigin(0, 0).setDisplaySize(50, 50).setInteractive()
        // this.add.image(500, 500, 'goodImage').setOrigin(0, 0).setDisplaySize(50, 50).setInteractive()

        this.playersTexts[0] = this.add.text(playerTextPositions[0].x, playerTextPositions[0].y, this.playerName).setColor('white').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true)
        this.playersTexts[1] = this.playersTexts[1] = this.add.text(playerTextPositions[1].x, playerTextPositions[1].y, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);
        this.playersTexts[2] = this.add.text(playerTextPositions[2].x, playerTextPositions[2].y, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;
        this.playersTexts[3] = this.add.text(playerTextPositions[3].x, playerTextPositions[3].y, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;

        this.prepareBtn = this.add.image(prepareImagePositions[0].x, prepareImagePositions[0].y, 'beginGame').setOrigin(0, 0).setDisplaySize(225, 250).setInteractive()
        // pointerup pointerdown focus
        this.prepareBtn.visible = false
        this.ourPlayNumberText = this.add.text(screenWidth - 200, 10, "我方 2").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.yourPlayNumberText = this.add.text(screenWidth - 200, 30, "对方 2").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
    
        this.showScore = this.add.text(screenWidth - 100, 10, "得分\n0").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.score = 0

        this.websocket = new WebSocket("ws://127.0.0.1:8080/ws")
        this.websocket.onopen = this.onopen.bind(this)
        this.websocket.onmessage = this.onmessage.bind(this)
        this.websocket.onclose = this.onclose.bind(this)
    }

    heartCheckStart() {
        let self = this
        this.heartCheckTimeoutObj = setTimeout(function(){
            self.websocket.send("ping");
            self.heartCheckServerTimeoutObj = setTimeout(function(){
                console.log("服务器无响应，断开连接!!!", dayjs().format())
                // self.websocket.close();//如果onclose会执行reconnect，我们执行ws.close()就行了.如果直接执行reconnect 会触发onclose导致重连两次
            }, self.heartCheckServerTimeout)
        }, this.heartCheckTimeout)
    }

    heartCheckReset() {
        clearTimeout(this.heartCheckTimeoutObj)
        clearTimeout(this.heartCheckServerTimeoutObj)
        this.heartCheckStart()
    }

    onopen() {
        console.log("连接成功")

        this.sendMessageToServer(SET_PLAYER_NAME_REQUEST, JSON.stringify({
            UUID: nanoid(),
            playerName: this.playerName,
        }))

        this.heartCheckStart()
    }

    onclose(e) {
        console.log("oncloase: 与服务器断开连接!!", dayjs().format())
        console.log('websocket 断开: ' + e.code + ' ' + e.reason + ' ' + e.wasClean)
        console.log(e)
    }

    onmessage(message) {
        if (message.data === "pong") {
            this.heartCheckReset()
            return
        }
        const data = JSON.parse(message.data)
        const messageType = data["messageType"]
        const content = data["content"]
        console.log(messageType)
        console.log(content)
        if (messageType === ROOM_LIST_RESPONSE) {
            const roomList = JSON.parse(content)
            this.sendMessageToServer(JOIN_ROOM_REQUEST, roomList[0])
        } else if (messageType === EXISTS_PLAYERS_RESPONSE) {
            const playersMsgs: Player[] = JSON.parse(content)
            this.players[0] = playersMsgs[0]
            for (let i = 1; i < 4; i++) {
                if (playersMsgs[i] === null) {
                    continue
                }
                this.playersTexts[i].setText(playersMsgs[i].name).setColor('white')
                if ((this.players[i] === null || !this.players[i].prepare) && playersMsgs[i].prepare) {
                    this.prepareOkImg[i] = this.add.image(prepareOkPositions[i].x, prepareOkPositions[i].y, 'prepareOk').setOrigin(0, 0).setDisplaySize(150, 125).setInteractive()
                }
                console.log(i, playersMsgs[i])
                this.players[i] = playersMsgs[i]
            }
            this.getPlayerMsgCnt += 1
            if (this.getPlayerMsgCnt === 1) {
                this.prepareBtn.visible = true
                this.prepareBtn.on("pointerdown", function (event) {
                    this.prepareBtn.destroy()
                    this.prepareOkImg[0] = this.add.image(prepareOkPositions[0].x, prepareOkPositions[0].y, 'prepareOk').setOrigin(0, 0).setDisplaySize(150, 125).setInteractive()
                    this.sendMessageToServer(PREPARE_REQUEST, "")
                }.bind(this))
            }
        } else if (messageType === DEAL_POKER) {
            this.dealPoker(JSON.parse(content))
        } else if (messageType === MATCH_BEGIN) {
            for (let i = 0; i < 4; i++) {
                this.prepareOkImg[i].destroy()
            }
        } else if (messageType === SHOW_MASTER_RESPONSE) {
            this.gameDetail.onShowMasterRes(JSON.parse(content))
            this.showZhuangColor()
            this.showPokers()
        } else if (messageType === REVEIVE_HOLE_CARDS) {
            this.showBuckleCards(JSON.parse(content))
        } else if (messageType === PLAY_TURN) {
            this.playCardsImg = this.add.image(playCardsImgX, playCardsImgY, 'playCardsGreenImg').setOrigin(0, 0).setDisplaySize(96, 40).setInteractive()
            // console.log(this.playCardsImg.texture.key, "hhhhhhhhhhhhhhh")
            this.playCardsImg.on('pointerup', function () {
                this.sendMessageToServer(PLAY_CARDS, JSON.stringify(this.playCards))
                for (let i = 0; i < this.playCards.length; i++) {
                    this.removePoker(this.playCards[i])
                }
                this.showPlayCards(0, this.playCards, false)
                this.showPokers()
                this.playCardsImg.destroy()
            }.bind(this))
        } else if (messageType === SHOW_PLAY_CARDS) {
            let showPlayCardsResponse: ShowPlayCardsResponse = JSON.parse(content)
            this.showPlayCards(showPlayCardsResponse.showIdx, showPlayCardsResponse.cards, false)
        } else if (messageType === BIGGEST_POSITION) {
            this.showBiggestPlayCards(Number(content))
        } else if (messageType === INCREASE_SCORES) {
            this.score += Number(content)
            this.showScore.setText("得分\n"+this.score)
        } else if (messageType === ROUND_END) {
            this.destoryShowPlayCardsImgs()
        } else if (messageType === GAME_RESULT) {
            this.initNextGame(JSON.parse(content))
        }
    }

    initNextGame(gameResult: GameResult) {
        this.gameDetail = new GameDetail(this)
        if (this.zhuangColorImg) {
            this.zhuangColorImg.destroy()
        }
        this.showScore.setText("得分\n0")
        this.score = 0
        this.isFirstGame = false
        this.ourPlayNumberText = this.add.text(screenWidth - 200, 10, "我方 " + gameResult.ourRound).setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.yourPlayNumberText = this.add.text(screenWidth - 200, 30, "对方 " + gameResult.otherRound).setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.showZhuang(gameResult.bankerPostion)
    }

    destoryShowPlayCardsImgs() {
        for (let i = 0; i < 4; i++) {
            for (let j = 0; j < this.showPlayCardsImgs[i].length; j++) {
                let attachBiggest = this.showPlayCardsImgs[i][j].getData("biggest")
                if (attachBiggest !== undefined) {
                    attachBiggest.destroy()
                }
                this.showPlayCardsImgs[i][j].destroy()
            }
        }
        this.biggestPosition = -1
    }

    showBiggestPlayCards(idx: number) {
        if (this.biggestPosition !== -1) {
            let cardsLen = this.showPlayCardsImgs[this.biggestPosition].length
            let attachBiggest = this.showPlayCardsImgs[this.biggestPosition][cardsLen - 1].getData("biggest")
            if (attachBiggest !== undefined) {
                attachBiggest.destroy()
            }
        }
        this.biggestPosition = idx
        let cardsLen = this.showPlayCardsImgs[idx].length
        let img = this.showPlayCardsImgs[idx][cardsLen - 1]
        let x = img.x + img.width - 30
        let y = img.y + img.height  - 30
        let attachStarImg = this.add.image(x, y, 'goodImage').setOrigin(0, 0).setDisplaySize(30, 30).setInteractive()
        img.setData("biggest", attachStarImg)
    }

    showPlayCards(idx: number, cards: Poker[], isBiggest: boolean) {
        let image = null
        for (let i = 0; i < cards.length; i++) {
            image = this.add.sprite(
                showPokerPositions[idx][i].x,
                showPokerPositions[idx][i].y,
                'poker',
                getPokerPosition(cards[i])
            ).setOrigin(0, 0).setInteractive()
            this.showPlayCardsImgs[idx].push(image)
        }
        if (isBiggest) {
            let x = image.x + image.width - 30
            let y = image.y + image.height - 30
            let attachStarImg = this.add.image(x, y, 'goodImage').setOrigin(0, 0).setDisplaySize(30, 30).setInteractive()
            image.setData("biggest", attachStarImg)
        }
    }

    showZhuangColor() {
        if (this.zhuangColorImg) {
            this.zhuangColorImg.destroy()
        }
        let isSelfTeamShow = this.gameDetail.isSelfTeamShowMaster()
        let positionX = screenWidth - 122
        let positionY = 10
        if (!isSelfTeamShow) {
            positionY = 30
        }
        if (this.gameDetail.isShowJokerMaster()) {
            this.zhuangColorImg = this.add.image(positionX, positionY, 'jokerImage').setOrigin(0, 0).setDisplaySize(20, 20).setInteractive()
        } else {
            let showColorIdx = showColorIdxMap.get(this.gameDetail.masterFlower)!
            this.zhuangColorImg = this.add.sprite(positionX, positionY, 'flowerImages', showColorIdx).setOrigin(0, 0).setDisplaySize(20, 20).setInteractive()
        }
    }

    appendPoker(poker: Poker) {
        if (poker.number === "joker") {
            this.playersCards.jokers.push(poker)
            this.playersCards.jokers.sort((a, b): Number => {
                if (a.color === b.color) {
                    return 0
                }
                if (a.color === "red") {
                    return 1
                }
                return -1
            })
        } else if (poker.number === this.gameDetail.playNumber) {
            this.playersCards.playNumberCards.push(poker)
        } else if (poker.color === SPADE) {
            this.playersCards.spadeCards.push(poker)
            this.playersCards.spadeCards.sort((a, b): Number => {
                return CardValueMap.get(b.number)! - CardValueMap.get(a.number)!
            })
        } else if (poker.color === HEART) {
            this.playersCards.heartCards.push(poker)
            this.playersCards.heartCards.sort((a, b): Number => {
                return CardValueMap.get(b.number)! - CardValueMap.get(a.number)!
            })
        } else if (poker.color === CLUB) {
            this.playersCards.clubCards.push(poker)
            this.playersCards.clubCards.sort((a, b): Number => {
                return CardValueMap.get(b.number)! - CardValueMap.get(a.number)!
            })
        } else {
            this.playersCards.dianmondCards.push(poker)
            this.playersCards.dianmondCards.sort((a, b): Number => {
                return CardValueMap.get(b.number)! - CardValueMap.get(a.number)!
            })
        }

        this.playersCards.cardNum += 1

        this.gameDetail.onDealPoker(poker)
    }

    removePoker(poker: Poker) {
        if (poker.number === "joker") {
            for (let i = 0; i < this.playersCards.jokers.length; i++) {
                if (poker === this.playersCards.jokers[i]) {
                    this.playersCards.jokers.splice(i, 1)
                    break
                }
            }
        } else if (poker.number === this.gameDetail.playNumber) {
            for (let i = 0; i < this.playersCards.playNumberCards.length; i++) {
                if (poker === this.playersCards.playNumberCards[i]) {
                    this.playersCards.playNumberCards.splice(i, 1)
                    break
                }
            }
        } else if (poker.color === SPADE) {
            for (let i = 0; i < this.playersCards.spadeCards.length; i++) {
                if (poker === this.playersCards.spadeCards[i]) {
                    this.playersCards.spadeCards.splice(i, 1)
                    break
                }
            }
        } else if (poker.color === HEART) {
            for (let i = 0; i < this.playersCards.heartCards.length; i++) {
                if (poker === this.playersCards.heartCards[i]) {
                    this.playersCards.heartCards.splice(i, 1)
                    break
                }
            }
        } else if (poker.color === CLUB) {
            for (let i = 0; i < this.playersCards.clubCards.length; i++) {
                if (poker === this.playersCards.clubCards[i]) {
                    this.playersCards.clubCards.splice(i, 1)
                    break
                }
            }
        } else {
            for (let i = 0; i < this.playersCards.dianmondCards.length; i++) {
                if (poker === this.playersCards.dianmondCards[i]) {
                    this.playersCards.dianmondCards.splice(i, 1)
                    break
                }
            }
        }

        this.playersCards.cardNum -= 1

        this.gameDetail.onRemovePoker(poker)
    }

    dealPoker(poker: Poker) {
        this.appendPoker(poker)
        this.showPokers()
        this.gameDetail.showMaster()

        if (this.playersCards.cardNum === FULL_POKER_NUM) {
            this.waitShowMaster()
        }
    }

    clearPokers() {
        for (let i = 0; i < this.pokerImages.length; i++) {
            let attachStar = this.pokerImages[i].getData("star")
            this.pokerImages[i].destroy()
            if (attachStar !== undefined) {
                attachStar.destroy()
            }
        }
        this.pokerImages = []
    }

    showPokers() {
        this.clearPokers()
        this.playCards = []
        let position = this.showSomePokers(0, this.playersCards.jokers)
        position = this.showSomePokers(position, this.playersCards.playNumberCards, true)
        if (this.gameDetail.masterFlower === SPADE || this.gameDetail.masterFlower === "" || this.gameDetail.masterFlower === "red" || this.gameDetail.masterFlower === "black") {
            if (this.gameDetail.masterFlower === SPADE) {
                position = this.showSomePokers(position, this.playersCards.spadeCards, true)
            } else {
                position = this.showSomePokers(position, this.playersCards.spadeCards, false)
            }
            position = this.showSomePokers(position, this.playersCards.heartCards, false)
            position = this.showSomePokers(position, this.playersCards.clubCards, false)
            this.showSomePokers(position, this.playersCards.dianmondCards, false)
        } else if (this.gameDetail.masterFlower === HEART) {
            position = this.showSomePokers(position, this.playersCards.heartCards, true)
            position = this.showSomePokers(position, this.playersCards.clubCards, false)
            position = this.showSomePokers(position, this.playersCards.dianmondCards, false)
            this.showSomePokers(position, this.playersCards.spadeCards, false)
        } else if (this.gameDetail.masterFlower === CLUB) {
            position = this.showSomePokers(position, this.playersCards.clubCards, true)
            position = this.showSomePokers(position, this.playersCards.dianmondCards, false)
            position = this.showSomePokers(position, this.playersCards.spadeCards, false)
            position = this.showSomePokers(position, this.playersCards.heartCards, false)
        } else {
            position =  this.showSomePokers(position, this.playersCards.dianmondCards, true)
            position = this.showSomePokers(position, this.playersCards.spadeCards, false)
            position = this.showSomePokers(position, this.playersCards.heartCards, false)
            position = this.showSomePokers(position, this.playersCards.clubCards, false)
        }        
    }

    showBuckleCards(holeCards: Poker[]) {
        this.clearPokers()
        this.showKouImg(false)
        for (let i = 0; i < 8; i++) {
            this.appendPoker(holeCards[i])
        }
        let position = this.showSomeBuckleCards(0, this.playersCards.jokers)
        position = this.showSomeBuckleCards(position, this.playersCards.playNumberCards)
        position = this.showSomeBuckleCards(position, this.playersCards.spadeCards)
        position = this.showSomeBuckleCards(position, this.playersCards.heartCards)
        position = this.showSomeBuckleCards(position, this.playersCards.clubCards)
        this.showSomeBuckleCards(position, this.playersCards.dianmondCards)
    }

    showKouImg(isLight: boolean) {
        if (this.kouImg) {
            this.kouImg.destroy()
        }
        let imgName = "kouGreyImage"
        if (isLight) {
            imgName = "kouLightImage"
        }
        this.kouImg = this.add.image(screenWidth * 0.48, screenHeight - 380, imgName).setOrigin(0, 0).setDisplaySize(90, 40).setInteractive()
        if (isLight) {
            this.kouImg.on('pointerup', () => {
                this.kouImg.destroy()
                for (let i = 0; i < 8; i++) {
                    this.removePoker(this.buckleCards[i])
                }
                console.log(this.buckleCards)
                this.showPokers()
                this.sendMessageToServer(KOU_CARDS, JSON.stringify(this.buckleCards))
                this.buckleCards = []
            })
        }
    }

    showSomeBuckleCards(position: number, cards: Poker[]): number {
        for (let i = 0; i < cards.length; i++) {
            let x = buckleCardPositions[position].x
            let y = buckleCardPositions[position].y
            let image = this.add.sprite(x, y, 'poker', getPokerPosition(cards[i])).setOrigin(0, 0).setInteractive()
            let tmpCard = cards[i]
            image.on('pointerup', () => {
                if (image.data === null || image.getData("status") === "down") {
                    image.setData("status", "up")
                    image.y -= 30
                    this.selectBuckleNum += 1
                    this.buckleCards.push(tmpCard)
                    if (this.selectBuckleNum === 8) {
                        this.showKouImg(true)
                    }
                } else {
                    image.setData("status", "down")
                    image.y += 30
                    this.selectBuckleNum -= 1
                    this.removeBuckleCard(tmpCard)
                    if (this.selectBuckleNum < 8) {
                        this.showKouImg(false)
                    }
                }
            })
            this.pokerImages.push(image)
            position += 1
        }
        return position
    }

    removeBuckleCard(poker: Poker) {
        for (let i = 0; i < this.buckleCards.length; i++) {
            if (poker === this.buckleCards[i]) {
                this.buckleCards.splice(i, 1)
                break
            }
        }
    }

    showSomePokers(position: number, cards: Poker[], addStar: boolean): number {
        for (let i = 0; i < cards.length; i++) {
            let x = pokerPositions[position].x
            let y = pokerPositions[position].y
            let image = this.add.sprite(x, y, 'poker', getPokerPosition(cards[i])).setOrigin(0, 0).setInteractive()
            let attachStarImg = null
            if (addStar) {
                x = image.x
                y = image.y + image.height - 20
                attachStarImg = this.add.image(x, y, 'miniStarImage').setOrigin(0, 0).setDisplaySize(20, 20).setInteractive()
                image.setData("star", attachStarImg)
            }
            image.on('pointerup', () => {
                if (image.getData("status") === undefined || image.getData("status") === "down") {
                    image.setData("status", "up")
                    image.y -= 30
                    if (addStar) {
                        attachStarImg.y -= 30
                    }
                    this.playCards.push(cards[i])
                } else {
                    image.setData("status", "down")
                    image.y += 30
                    if (addStar) {
                        attachStarImg.y += 30
                    }
                    this.removePlayCard(cards[i])
                }
                console.log(this.playCards)
            })
            this.pokerImages.push(image)
            position += 1
        }
        return position
    }

    removePlayCard(poker: Poker) {
        for (let i = 0; i < this.playCards.length; i++) {
            if (this.playCards[i] === poker) {
                this.playCards.splice(i, 1)
                return
            }
        }
    }

    destoryAllShowMaster() {
        for (let i = 0; i < this.showMasterImages.length; i++) {
            this.showMasterImages[i].destroy()
        }
    }

    waitShowMaster() {
        console.log(dayjs().format())
        let setIntervalID = setInterval(() => {
            console.log(dayjs().format())
        }, 1000)
        setTimeout(() => {
            console.log(dayjs().format())
            clearInterval(setIntervalID)
            this.destoryAllShowMaster()
            this.sendMessageToServer(SHOW_MASTER_DONE, "")
            if (this.isFirstGame) {
                this.showZhuang(this.gameDetail.showMasterPosition)
            }
        }, 9000)
    }

    showZhuang(pos: number) {
        if (pos === -1) {
            return
        }
        if (this.zhuangImg) {
            this.zhuangImg.destroy()
        }
        this.zhuangImg = this.add.image(zhuangPositions[pos].x, zhuangPositions[pos].y, 'zhuangImage').setOrigin(0, 0).setDisplaySize(30, 20).setInteractive()
    }

    sendMessageToServer(messageType: string, content: string) {
        this.websocket.send(JSON.stringify({
            "messageType": messageType, "content": content
        }))
    }
}