import { GameScene } from './game_scene';
import { 
    Poker, SPADE, HEART, CLUB, DIANMOND, SHOW_MASTER, ShowMasterRequest, ShowMasterResponse,
} from './dto'
import { Game } from 'phaser';

const screenHeight =  document.documentElement.clientHeight;
const showMasterYPosition = screenHeight - 230
const showMasterXPositions = [
    300, 370, 440, 510, 580
]

export class GameDetail {
    private gameScene: GameScene
    public playNumber: string
    // joker, spade, heart, club, dianmond
    public masterFlower: string
    public isMasterProtect: boolean
    public isSelfShowMaster: boolean
    public showMasterPosition: number
    public blackJokerCnt: number
    public redJokerCnt: number
    public playNumCntMap: Map<string, number>
    constructor(gameScene: GameScene) {
        this.gameScene = gameScene
        this.playNumber = "2"
        this.masterFlower = ""
        this.showMasterPosition = -1
        this.isMasterProtect = false
        this.isSelfShowMaster = false
        this.blackJokerCnt = 0
        this.redJokerCnt = 0
        this.playNumCntMap = new Map([
            [SPADE, 0],
            [HEART, 0],
            [CLUB, 0],
            [DIANMOND, 0]
        ])
    }

    public onShowMasterRes(res: ShowMasterResponse) {
        this.masterFlower = res.color
        this.isMasterProtect = res.isProtect
        this.isSelfShowMaster = res.isSelfShowMaster
        this.showMasterPosition = res.showMasterPosition
    }

    public onDealPoker(poker: Poker) {
        if (poker.color === "red") {
            this.redJokerCnt += 1
            return
        }
        if (poker.color === "black") {
            this.blackJokerCnt += 1
            return
        }
        if (poker.number !== this.playNumber) {
            return
        }
        this.playNumCntMap.set(poker.color, this.playNumCntMap.get(poker.color)! + 1)
    }

    public isCanShowJoker() {
        if (this.isSelfProtect()) {
            return false
        }
        if (this.masterFlower === "") {
            if (this.blackJokerCnt === 2 || this.redJokerCnt === 2) {
                return true
            }
        }
        return this.blackJokerCnt + this.redJokerCnt >= 3
    }

    public isCanShowColor(color: string) {
        if (this.isSelfProtect()) {
            return false
        }
        let jokerCnt = this.redJokerCnt
        if (color === SPADE || color === CLUB) {
            jokerCnt = this.blackJokerCnt
        }
        if (this.isMasterProtect) {
            return false
        }
        if (this.masterFlower === "") {
            return jokerCnt > 0 && this.playNumCntMap.get(color)! > 0
        }
        return jokerCnt > 0 && this.playNumCntMap.get(color)! === 2
    }

    isSelfProtect(): boolean {
        return this.isSelfShowMaster && this.isMasterProtect
    }

    public sendShowMaster(color: string) {
        let req: ShowMasterRequest = {
            color: color,
            isSelfProtect: false,
            isOppose: false,
        }
        if (this.masterFlower === "") {
            req.isSelfProtect = this.isColorProtect(color)
        } else {
            req.isSelfProtect = true
            req.isOppose = !this.isSelfShowMaster
        }
        this.gameScene.sendMessageToServer(SHOW_MASTER, JSON.stringify(req))
    } 

    isColorProtect(color: string): boolean {
        if (color === "red" && this.redJokerCnt === 2 && this.blackJokerCnt > 0) {
            return true
        } 
        if (color === "black" && this.blackJokerCnt === 2 && this.redJokerCnt > 0) {
            return true
        }
        return this.playNumCntMap.get(color) === 2
    }

    public showMaster() {
        let showCnt = 0
        console.log("hhhhhhhh")
        if (this.gameScene.gameDetail.isCanShowJoker()) {
            let image = this.gameScene.add.image(showMasterXPositions[showCnt], showMasterYPosition, 'jokerImage').setOrigin(0, 0).setDisplaySize(48, 48).setInteractive()
            this.gameScene.showMasterImages.push(image)
            image.on("pointerdown", function (this: GameDetail) {
                this.gameScene.destoryAllShowMaster()
                let master = "black"
                if (this.gameScene.gameDetail.blackJokerCnt < 2) {
                    master = "red"
                }
                this.sendShowMaster(master)
            }.bind(this))
            showCnt += 1
        }
        // 红，黑，梅，方
        if (this.isCanShowColor(HEART)) {
            let image = this.gameScene.add.sprite(showMasterXPositions[showCnt], showMasterYPosition, 'flowerImages', 0).setOrigin(0, 0).setDisplaySize(48, 48).setInteractive()
            this.gameScene.showMasterImages.push(image)
            image.on("pointerdown", function (this: GameDetail) {
                this.gameScene.destoryAllShowMaster()
                this.sendShowMaster(HEART)
            }.bind(this))
            showCnt += 1
        }
        if (this.isCanShowColor(SPADE)) {
            let image = this.gameScene.add.sprite(showMasterXPositions[showCnt], showMasterYPosition, 'flowerImages', 1).setOrigin(0, 0).setDisplaySize(48, 48).setInteractive()
            this.gameScene.showMasterImages.push(image)
            image.on("pointerdown", function (this: GameDetail) {
                this.gameScene.destoryAllShowMaster()
                this.sendShowMaster(SPADE)
            }.bind(this))
            showCnt += 1
        }
        if (this.isCanShowColor(CLUB)) {
            let image = this.gameScene.add.sprite(showMasterXPositions[showCnt], showMasterYPosition, 'flowerImages', 2).setOrigin(0, 0).setDisplaySize(48, 48).setInteractive()
            this.gameScene.showMasterImages.push(image)
            image.on("pointerdown", function (this: GameDetail) {
                this.gameScene.destoryAllShowMaster()
                this.sendShowMaster(CLUB)
            }.bind(this))
            showCnt += 1
        }
        if (this.isCanShowColor(DIANMOND)) {
            let image = this.gameScene.add.sprite(showMasterXPositions[showCnt], showMasterYPosition, 'flowerImages', 3).setOrigin(0, 0).setDisplaySize(48, 48).setInteractive()
            this.gameScene.showMasterImages.push(image)
            image.on("pointerdown", function (this: GameDetail) {
                this.gameScene.destoryAllShowMaster()
                this.sendShowMaster(DIANMOND)
            }.bind(this))
        }
       
    }
}