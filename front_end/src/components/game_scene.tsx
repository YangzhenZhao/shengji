// @ts-nocheck
import Phaser from "phaser";
import bgimage from "../assets/bg2.jpg"
import pokerImage from "../assets/poker.png"
import beginGame from "../assets/btn/begin.png"
import prepareOk from "../assets/btn/prepare_ok.png"
import { nanoid } from 'nanoid'

const SET_PLAYER_NAME_REQUEST = "set_player_name"
const JOIN_ROOM_REQUEST = "join_room"
const PREPARE_REQUEST = "prepare"

const ROOM_LIST_RESPONSE = "room_list"
const EXISTS_PLAYERS_RESPONSE = "exists_players"

const screenWidth = document.documentElement.clientWidth;
const screenHeight =  document.documentElement.clientHeight;

// 屏幕左上角为坐标原点 (0, 0), 横轴为 x, 纵轴为 y
const playerTextPositions = [
    {x: 50, y: screenHeight - 36},
    {x: screenWidth * 0.6, y: 30},
    {x: 5, y: screenHeight / 2.4},
    {x: screenWidth - 108, y: screenHeight / 2.4}
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

interface Player {
    name: string
    prepare: boolean
}

export class GameScene extends Phaser.Scene {
    public playerName
    public websocket: WebSocket | null
    public players: Player[]
    public prepareBtn: Phaser.GameObjects.Image
    public prepareOkImg: Phaser.GameObjects.Image[]
    constructor(playerName: string) {
        super("GameScene")
        this.playerName = playerName
        this.existPlayers = [1]
        this.playersTexts = []
        this.players = [{name: playerName, prepare: false}, null, null, null]
        this.websocket = null
        this.getPlayerMsgCnt = 0
        this.prepareOkImg = [null, null, null, null]
    }

    preload () {
        this.load.image("bg2", bgimage)
        this.load.image("beginGame", beginGame)
        this.load.image("prepareOk", prepareOk)
        this.load.spritesheet('poker', pokerImage, {
            frameWidth: 90,
            frameHeight: 120
        });
    }

    create() {
        this.add.image(0, 0, 'bg2').setOrigin(0).setDisplaySize(screenWidth, screenHeight);

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
    
    
        var x = 200;
        var y = screenHeight - 160;
    
        for (let i = 0; i < 25; i++) {
            let image = this.add.sprite(x, y, 'poker', i).setOrigin(0, 0).setInteractive();
            image.on('pointerup', () => {
                if (image.data === null || image.getData("status") === "down") {
                image.setData("status", "up")
                image.y -= 30
    
                this.score += 5
                this.showScore.setText("得分\n"+this.score)
                } else {
                image.setData("status", "down")
                image.y += 30
                }
            })
            x += 26;
        }

        this.websocket = new WebSocket("ws://192.168.1.115:8080/ws")
        this.websocket.onopen = this.onopen.bind(this)
        this.websocket.onmessage = this.onmessage.bind(this)
    }

    onopen() {
        console.log("连接成功")
        

        this.sendMessageToServer(SET_PLAYER_NAME_REQUEST, JSON.stringify({
            UUID: nanoid(),
            playerName: this.playerName,
        }))
    }

    onmessage(message) {
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
                this.players[i] = playersMsgs[i]
                this.playersTexts[i].setText(playersMsgs[i].name).setColor('white')
                console.log(i, "hhhhhhhh", playersMsgs[i].prepare)
                if (playersMsgs[i].prepare) {
                    console.log('hhhhhhhhhhh ...........')
                    this.prepareOkImg[i] = this.add.image(prepareOkPositions[i].x, prepareOkPositions[i].y, 'prepareOk').setOrigin(0, 0).setDisplaySize(150, 125).setInteractive()
                }
                console.log(i, playersMsgs[i])
            }
            this.getPlayerMsgCnt += 1
            if (this.getPlayerMsgCnt === 1) {
                this.prepareBtn.visible = true
                this.prepareBtn.on("pointerdown", function (event) {
                    this.prepareBtn.destroy()
                    this.prepareOkImg[0] = this.add.image(prepareOkPositions[0].x, prepareOkPositions[0].y, 'prepareOk').setOrigin(0, 0).setDisplaySize(150, 125).setInteractive()
                    console.log('hhh hover')
                    this.sendMessageToServer(PREPARE_REQUEST, "")
                }.bind(this))
            }
        }
    }
    
    sendMessageToServer(messageType: string, content: string) {
        this.websocket.send(JSON.stringify({
            "messageType": messageType, "content": content
        }))
    }
}