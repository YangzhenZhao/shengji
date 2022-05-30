import { GameScene } from './game_scene';

export class GameDetail {
    private gameScene: GameScene
    public playNumber: string
    constructor(gameScene: GameScene) {
        this.gameScene = gameScene
        this.playNumber = "2"
    }
}