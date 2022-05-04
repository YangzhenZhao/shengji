// @ts-nocheck
import Phaser from "phaser";
import bgimage from "../bg2.jpg"
import pokerImage from "../poker.png"
import { nanoid } from 'nanoid'

const SET_PLAYER_NAME_REQUEST = "set_player_name"
const JOIN_ROOM_REQUEST = "join_room"

const ROOM_LIST_RESPONSE = "room_list"
const EXISTS_PLAYERS_RESPONSE = "exists_players"

let graphics: any

interface Player {
    name: string
    isPrepare: boolean
}

export class GameScene extends Phaser.Scene {
    public playerName
    public websocket: WebSocket | null
    public players: Player[]
    constructor(playerName: string) {
        super("GameScene")
        this.playerName = playerName
        this.existPlayers = [1]
        this.players = [{name: playerName, isPrepare: false}, null, null, null]
        this.websocket = null
    }

    preload () {
        this.load.image("bg2", bgimage)
        this.load.spritesheet('poker', pokerImage, {
            frameWidth: 90,
            frameHeight: 120
        });
        this.websocket = new WebSocket("ws://192.168.1.115:8080/ws")
        this.websocket.onopen = this.onopen.bind(this)
        this.websocket.onmessage = this.onmessage.bind(this)
    }

    create() {
        const screenWidth = document.documentElement.clientWidth;
        const screenHeight =  document.documentElement.clientHeight;
        const myPosition = this.existPlayers.length
        let teammatePosition = 7 - myPosition
        let otherSidePosition1 = 1
        let otherSidePosition2 = 2
        if (myPosition <= 2) {
            teammatePosition = 3 - myPosition
            otherSidePosition1 = 3
            otherSidePosition2 = 4
        }

        this.add.image(0, 0, 'bg2').setOrigin(0).setDisplaySize(screenWidth, screenHeight);

        this.add.text(50, screenHeight - 36, this.playerName).setColor('white').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true)
        
        if (this.existPlayers.length >= teammatePosition) {
            this.add.text(screenWidth * 0.6, 30, existPlayers[teammatePosition - 1].name).setColor('white').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true)
        } else {
            this.add.text(screenWidth * 0.6, 30, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);
        }
        
        if (this.existPlayers.length >= otherSidePosition1) {
            this.add.text(5, screenHeight / 2.4, this.existPlayers[otherSidePosition1 - 1].name).setColor('white').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;
        } else {
            this.add.text(5, screenHeight / 2.4, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;
        }
        
        if (this.existPlayers.length >= otherSidePosition2) {
            this.add.text(screenWidth - 108, screenHeight / 2.4, existPlayers[otherSidePosition2 - 1].name).setColor('white').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;
        } else {
            this.add.text(screenWidth - 108, screenHeight / 2.4, '空闲位置').setColor('red').setFontSize(28).setShadow(2, 2, "#333333", 2, true, true);;
        }

        this.ourPlayNumberText = this.add.text(screenWidth - 200, 10, "我方 2").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.yourPlayNumberText = this.add.text(screenWidth - 200, 30, "对方 2").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
    
        this.showScore = this.add.text(screenWidth - 100, 10, "得分\n0").setFontSize(20).setShadow(2, 2, "#333333", 2, true, true);
        this.score = 0
    
        // graphics = this.add.graphics({ x: 0, y: 0 });
    
        // graphics.lineStyle(4, 0xffd700, 1);
        // graphics.strokeCircle(screenWidth / 2, screenHeight - 210, 40);
    
        // this.waitTime = this.add.text(screenWidth / 2 - 16, screenHeight - 216, "0秒").setColor('red').setFontSize(18)
        // this.second = 0
        // this.time.addEvent({ delay: 1000, callback: this.changeWaitTime, callbackScope: this, loop: true })
    
    
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
    }

    onopen() {
        console.log("连接成功")
        this.sendMessageToServer(SET_PLAYER_NAME_REQUEST, JSON.stringify({
            UUID: nanoid(),
            playerName: this.playerName,
        }))
    }

    onmessage(message) {
        console.log(message.data)
    }

    // changeWaitTime() {
    //     this.second += 2
    //     this.waitTime.setText(this.second+"秒")
    // }
    
    sendMessageToServer(messageType: string, content: string) {
        this.websocket.send(JSON.stringify({
            "messageType": messageType, "content": content
        }))
    }
}